package azure

import (
	"context"
	"fmt"
	"os"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/cache"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("azure", newProvider, magicProvider)
}

type azureProvider struct {
	name       string
	credential azcore.TokenCredential
	config     struct {
		Prefix       string
		Subscription string
		TenantId     string
		ClientId     string
		// Client credentials
		ClientSecret string // #nosec G117 -- Credential field required for Azure API auth
		// Certificate auth
		CertificatePath     string
		CertificatePassword string
	}
}

func newProvider(name string) herd.HostProvider {
	p := &azureProvider{name: name}
	return p
}

func magicProvider() herd.HostProvider {
	sub, ok := os.LookupEnv("AZURE_SUBSCRIPTION")
	if !ok {
		return nil
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil
	}
	p := newProvider("azure").(*azureProvider)
	p.config.Subscription = sub
	p.credential = cred
	return cache.NewFromProvider(p)
}

func (p *azureProvider) Name() string {
	return p.name
}

func (p *azureProvider) Prefix() string {
	return p.config.Prefix
}

func (p *azureProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*azureProvider)
	return p.config.Subscription == op.config.Subscription
}

func (p *azureProvider) ParseViper(v *viper.Viper) error {
	viper.SetDefault("TenantId", "")
	if err := v.Unmarshal(&p.config); err != nil {
		return err
	}
	if p.credential == nil {
		// Let's see if we have authentication data
		if p.config.TenantId != "" && p.config.ClientId != "" {
			if p.config.ClientSecret != "" {
				cred, err := azidentity.NewClientSecretCredential(p.config.TenantId, p.config.ClientId, p.config.ClientSecret, nil)
				if err != nil {
					return err
				}
				p.credential = cred
			} else if p.config.CertificatePath != "" && p.config.CertificatePassword != "" {
				certData, err := os.ReadFile(p.config.CertificatePath)
				if err != nil {
					return err
				}
				certs, key, err := azidentity.ParseCertificates(certData, []byte(p.config.CertificatePassword))
				if err != nil {
					return err
				}
				cred, err := azidentity.NewClientCertificateCredential(p.config.TenantId, p.config.ClientId, certs, key, nil)
				if err != nil {
					return err
				}
				p.credential = cred
			}
		}
	}
	// If credential is still nil, try a default credential
	if p.credential == nil {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return err
		}
		p.credential = cred
	}
	return nil
}

func (p *azureProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	lm(p.name, false, nil)

	hosts := herd.NewHostSet()
	factory, err := armcompute.NewClientFactory(p.config.Subscription, p.credential, nil)
	if err != nil {
		return nil, err
	}
	client := factory.NewVirtualMachinesClient()
	pager := client.NewListAllPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed to get next page of virtual machines: %v", err)
		}
		logrus.Debugf("Got %d virtual machines", len(page.Value))
		for _, vm := range page.Value {
			attrs := make(herd.HostAttributes)
			for k, v := range vm.Tags {
				attrs[k] = *v
			}
			attrs["VMSize"] = string(*vm.Properties.HardwareProfile.VMSize)
			attrs["ProvisioningState"] = *vm.Properties.ProvisioningState
			attrs["ID"] = *vm.ID
			attrs["Location"] = *vm.Location
			hosts.AddHost(herd.NewHost(*vm.Name, "", attrs))
		}
	}

	return hosts, nil
}
