package aws

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/cache"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/seveas/scattergather"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("aws", newProvider, magicProvider)
}

type awsProvider struct {
	name   string
	config struct {
		Prefix           string
		AccessKeyId      string
		SecretAccessKey  string
		Partition        string
		Regions          []string
		ExcludeRegions   []string
		UsePublicAddress bool
	}
}

func newProvider(name string) herd.HostProvider {
	p := &awsProvider{name: name}
	p.config.Partition = "aws"
	return p
}

func magicProvider() herd.HostProvider {
	p := newProvider("aws").(*awsProvider)
	if v, ok := os.LookupEnv("AWS_ACCESS_KEY_ID"); ok {
		p.config.AccessKeyId = v
	}
	if v, ok := os.LookupEnv("AWS_ACCESS_KEY"); ok {
		p.config.AccessKeyId = v
	}
	if v, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); ok {
		p.config.SecretAccessKey = v
	}
	if v, ok := os.LookupEnv("AWS_SECRET_KEY"); ok {
		p.config.SecretAccessKey = v
	}
	if p.config.AccessKeyId != "" && p.config.SecretAccessKey != "" {
		return cache.NewFromProvider(p)
	}
	return nil
}

func (p *awsProvider) Name() string {
	return p.name
}

func (p *awsProvider) Prefix() string {
	return p.config.Prefix
}

func (p *awsProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*awsProvider)
	return p.config.AccessKeyId == op.config.AccessKeyId &&
		p.config.SecretAccessKey == op.config.SecretAccessKey &&
		p.config.Partition == op.config.Partition &&
		reflect.DeepEqual(p.config.Regions, op.config.Regions)
}

func (p *awsProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

func (p *awsProvider) setRegions() error {
	resolver := endpoints.DefaultResolver().(endpoints.EnumPartitions)
	var partition endpoints.Partition
	for _, partition = range resolver.Partitions() {
		if partition.ID() == p.config.Partition {
			break
		}
	}
	if partition.ID() != p.config.Partition {
		return fmt.Errorf("No such partition: %s", p.config.Partition)
	}
	svc := partition.Services()["ec2"]
	p.config.Regions = make([]string, 0)
	for region := range svc.Regions() {
		p.config.Regions = append(p.config.Regions, region)
	}
	return nil
}

func (p *awsProvider) Load(ctx context.Context, lm herd.LoadingMessage) (hosts *herd.HostSet, err error) {
	lm(p.name, false, nil)
	defer func() { lm(p.name, true, err) }()
	if len(p.config.Regions) == 0 {
		if err := p.setRegions(); err != nil {
			return nil, err
		}
	}
	logrus.Debugf("AWS regions: %v", p.config.Regions)
	sg := scattergather.New[*herd.HostSet](int64(len(p.config.Regions)))
	for _, region := range p.config.Regions {
		if len(p.config.ExcludeRegions) != 0 && stringInList(p.config.ExcludeRegions, region) {
			continue
		}
		sg.Run(ctx, func() (*herd.HostSet, error) {
			name := fmt.Sprintf("%s@%s", p.name, region)
			lm(name, false, nil)
			hosts, err := p.loadRegion(ctx, region)
			lm(name, true, err)
			return hosts, err
		})
	}

	allHosts, err := sg.Wait()
	return herd.MergeHostSets(allHosts), err
}

func (p *awsProvider) loadRegion(ctx context.Context, region string) (*herd.HostSet, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(p.config.AccessKeyId, p.config.SecretAccessKey, ""),
		Region:      aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	ret := new(herd.HostSet)
	var token *string = nil
	sv := aws.StringValue
	iv := aws.Int64Value
	for {
		out, err := svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{NextToken: token, MaxResults: aws.Int64(1000)})
		if err != nil {
			return nil, err
		}
		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				name := *instance.PublicDnsName
				if name == "" {
					name = *instance.InstanceId
				}
				attrs := herd.HostAttributes{
					"architecture":            sv(instance.Architecture),
					"hypervisor":              sv(instance.Hypervisor),
					"image_id":                sv(instance.ImageId),
					"instance_id":             sv(instance.InstanceId),
					"instance_type":           sv(instance.InstanceType),
					"launch_time":             *instance.LaunchTime,
					"availability_zone":       sv(instance.Placement.AvailabilityZone),
					"placement_group":         sv(instance.Placement.GroupName),
					"tenancy":                 sv(instance.Placement.Tenancy),
					"private_dns_name":        sv(instance.PrivateDnsName),
					"private_ip":              sv(instance.PrivateIpAddress),
					"public_dns_name":         sv(instance.PublicDnsName),
					"public_ip":               sv(instance.PublicIpAddress),
					"root_device_name":        sv(instance.RootDeviceName),
					"root_device_type":        sv(instance.RootDeviceType),
					"state":                   sv(instance.State.Name),
					"state_code":              iv(instance.State.Code),
					"state_transition_reason": sv(instance.StateTransitionReason),
					"subnet_id":               sv(instance.SubnetId),
					"virtualization_type":     sv(instance.VirtualizationType),
					"vpc_id":                  sv(instance.VpcId),
					"owner_id":                sv(reservation.OwnerId),
					"reservation_id":          sv(reservation.ReservationId),
				}
				if instance.CpuOptions != nil {
					attrs["core_count"] = iv(instance.CpuOptions.CoreCount)
					attrs["threads_per_core"] = iv(instance.CpuOptions.ThreadsPerCore)
				}
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						name = *tag.Value
					} else {
						attrs[*tag.Key] = *tag.Value
					}
				}
				addr := sv(instance.PrivateIpAddress)
				if p.config.UsePublicAddress {
					addr = sv(instance.PublicIpAddress)
				}
				ret.AddHost(herd.NewHost(name, addr, attrs))
			}
		}
		if out.NextToken == nil {
			break
		}
		token = out.NextToken
	}

	return ret, nil
}

func stringInList(haystack []string, needle string) bool {
	for _, twig := range haystack {
		if twig == needle {
			return true
		}
	}
	return false
}
