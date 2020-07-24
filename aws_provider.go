// +build !no_aws

package herd

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

func init() {
	availableProviders["aws"] = NewAwsProvider
	magicProviders["aws"] = func(r *Registry) {
		p := NewAwsProvider("aws").(*AwsProvider)
		if p.AccessKeyId != "" && p.SecretAccessKey != "" {
			r.AddMagicProvider(NewCacheFromProvider(p))
		}
	}
}

type AwsProvider struct {
	BaseProvider    `mapstructure:",squash"`
	AccessKeyId     string
	SecretAccessKey string
	Partition       string
	Regions         []string
}

func NewAwsProvider(name string) HostProvider {
	p := &AwsProvider{BaseProvider: BaseProvider{Name: name}, Partition: "aws"}

	if v, ok := os.LookupEnv("AWS_ACCESS_KEY_ID"); ok {
		p.AccessKeyId = v
	}
	if v, ok := os.LookupEnv("AWS_ACCESS_KEY"); ok {
		p.AccessKeyId = v
	}
	if v, ok := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); ok {
		p.SecretAccessKey = v
	}
	if v, ok := os.LookupEnv("AWS_SECRET_KEY"); ok {
		p.SecretAccessKey = v
	}
	return p
}

func (p *AwsProvider) Equivalent(o HostProvider) bool {
	if c, ok := o.(*Cache); ok {
		o = c.Source
	}
	op, ok := o.(*AwsProvider)
	return ok &&
		p.AccessKeyId == op.AccessKeyId &&
		p.SecretAccessKey == op.SecretAccessKey &&
		p.Partition == op.Partition &&
		reflect.DeepEqual(p.Regions, op.Regions)
}

func (p *AwsProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(p)
}

func (p *AwsProvider) setRegions() error {
	resolver := endpoints.DefaultResolver().(endpoints.EnumPartitions)
	var partition endpoints.Partition
	for _, partition = range resolver.Partitions() {
		if partition.ID() == p.Partition {
			break
		}
	}
	if partition.ID() != p.Partition {
		return fmt.Errorf("No such partition: %s", p.Partition)
	}
	svc := partition.Services()[endpoints.Ec2ServiceID]
	p.Regions = make([]string, 0)
	for region := range svc.Regions() {
		p.Regions = append(p.Regions, region)
	}
	return nil
}

func (p *AwsProvider) Load(ctx context.Context, mc chan CacheMessage) (Hosts, error) {
	if len(p.Regions) == 0 {
		if err := p.setRegions(); err != nil {
			return Hosts{}, err
		}
	}
	hosts := make(Hosts, 0)
	rc := make(chan loadresult)
	for _, region := range p.Regions {
		name := fmt.Sprintf("%s@%s", p.Name, region)
		mc <- CacheMessage{Name: name, Finished: false, Err: nil}
		go func(region, name string) {
			hosts, err := p.loadRegion(region)
			mc <- CacheMessage{Name: name, Finished: true, Err: err}
			rc <- loadresult{hosts: hosts, err: err}
		}(region, name)
	}
	todo := len(p.Regions)
	errs := &MultiError{}
	for todo > 0 {
		r := <-rc
		if r.err != nil {
			errs.Add(r.err)
		}
		hosts = append(hosts, r.hosts...)
		todo -= 1
	}
	if !errs.HasErrors() {
		return hosts, nil
	}
	return hosts, errs
}

func (p *AwsProvider) loadRegion(region string) (Hosts, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(p.AccessKeyId, p.SecretAccessKey, ""),
		Region:      aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := ec2.New(sess)
	ret := Hosts{}
	var token *string = nil
	sv := aws.StringValue
	iv := aws.Int64Value
	for {
		out, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{NextToken: token, MaxResults: aws.Int64(1000)})
		if err != nil {
			return nil, err
		}
		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				name := *instance.PublicDnsName
				if name == "" {
					name = *instance.InstanceId
				}
				attrs := HostAttributes{
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
				ret = append(ret, NewHost(name, attrs))
			}
		}
		if out.NextToken == nil {
			break
		}
		token = out.NextToken
	}

	return ret, nil
}
