package azure

import (
	"context"
	"os"

	"github.com/seveas/katyusha"
	"github.com/seveas/katyusha/provider/cache"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/viper"
)

func init() {
	katyusha.RegisterProvider("azure", newProvider, magicProvider)
}

type azureProvider struct {
	name       string
	authorizer autorest.Authorizer
	config     struct {
		Prefix       string
		Environment  string
		AdResource   string
		Subscription string
		TenantId     string
		ClientId     string
		// Client credentials
		ClientSecret string
		// Certificate auth
		CertificatePath     string
		CertificatePassword string
		// User/password auth
		Username string
		Password string
	}
}

func newProvider(name string) katyusha.HostProvider {
	p := &azureProvider{name: name}
	return p
}

func magicProvider() katyusha.HostProvider {
	sub, ok := os.LookupEnv("AZURE_SUBSCRIPTION")
	if !ok {
		return nil
	}
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil
	}
	p := newProvider("azure").(*azureProvider)
	p.authorizer = authorizer
	p.config.Subscription = sub
	return cache.NewFromProvider(p)
}

func (p *azureProvider) Name() string {
	return p.name
}

func (p *azureProvider) Prefix() string {
	return p.config.Prefix
}

func (p *azureProvider) Equivalent(o katyusha.HostProvider) bool {
	op := o.(*azureProvider)
	return p.config.Subscription == op.config.Subscription
}

func (p *azureProvider) ParseViper(v *viper.Viper) error {
	if err := v.Unmarshal(&p.config); err != nil {
		return err
	}
	var ac auth.AuthorizerConfig
	// Let's see if we have authentication data
	if p.config.TenantId != "" && p.config.ClientId != "" {
		if p.config.ClientSecret != "" {
			ac = auth.NewClientCredentialsConfig(p.config.ClientId, p.config.ClientSecret, p.config.TenantId)
		} else if p.config.CertificatePath != "" && p.config.CertificatePassword != "" {
			ac = auth.NewClientCertificateConfig(p.config.CertificatePath, p.config.CertificatePassword, p.config.ClientId, p.config.TenantId)
		} else if p.config.Username != "" && p.config.Password != "" {
			ac = auth.NewUsernamePasswordConfig(p.config.Username, p.config.Password, p.config.ClientId, p.config.TenantId)
		}
	}
	if ac != nil {
		authorizer, err := ac.Authorizer()
		if err != nil {
			return err
		}
		p.authorizer = authorizer
	} else {
		// Last resort: let's try to authorize from environment data
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err != nil {
			return nil
		}
		p.authorizer = authorizer
	}
	return nil
}

func (p *azureProvider) Load(ctx context.Context, lm katyusha.LoadingMessage) (hosts katyusha.Hosts, err error) {
	lm(p.name, false, nil)
	c := compute.NewVirtualMachinesClient(p.config.Subscription)
	c.Authorizer = p.authorizer
	res, err := c.ListAll(ctx, "")
	vms := res.Values()
	if err != nil {
		return nil, err
	}
	for res.NotDone() {
		err := res.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		vms = append(vms, res.Values()...)
	}
	hosts = make(katyusha.Hosts, len(vms))
	for i, vm := range vms {
		attrs := make(katyusha.HostAttributes)
		for k, v := range vm.Tags {
			attrs[k] = *v
		}
		attrs["VMSize"] = string(vm.HardwareProfile.VMSize)
		attrs["ProvisioningState"] = *vm.ProvisioningState
		attrs["VMID"] = *vm.VMID
		attrs["Location"] = *vm.Location
		hosts[i] = katyusha.NewHost(*vm.Name, attrs)
	}

	return hosts, nil
}
