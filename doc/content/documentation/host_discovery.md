---
Title: Host discovery
Weight: 2
---

Herd has a pluggable infrastructure for finding hosts in many places, from flat files to cloud
provider API's. To make herd find hosts, you can configure which providers to use and how to use
them. This section of the documentation explains which providers there are and how to configure
them.

# Magic providers

As you could see in the [getting started](../getting_started/) documentation, Herd can even find
hosts if you do not configure it. It does this with so called magic providers, that detect hosts
based on information on your filesystem or in your environment. For example, the known_hosts
provider knows where known_hosts files usually live and will load them automatically if the exist.
If you do not want to any magic providers to find hosts, you can use the `--no-magic-provider`
command line flag to skip this.

# Provider configuration

The provider section of the configuration looks as like the following in your config.yaml:

```yaml
Providers:
  consul:
    provider: cache
    lifetime: 1h
    source:
      provider: consul
      timeout: 10s
      address: http://consul.service.consul:8500
  aws-main:
    provider: cache
    lifetime: 8h
    source:
      provider: aws
      accesskeyid: AK1............
      secretaccesskey: SE1..............
      excluderegions: [eu-south-1]
  aws-backup:
    provider: cache
    lifetime: 8h
    source:
      provider: aws
      accesskeyid: AK2............
      secretaccesskey: SE2..............
      excluderegions: [eu-south-1]
  tailscale:
    provider: tailscale
    domain: tail.example.com
  seveas:
    provider: seveas
    prefix: 'svs:'
```

As you can see, the name of a provider section foes not imply which provider is used, you have to
specify the. As a result of this, you can easily use multiple instances of the same provider type,
and can wrap a provider in a cache provider without affecting other parts of the configuration.

If you want to temporarily disable a provider in your configuration, you can add `enabled: false` to
its parameters.

# Host attributes

Providers do not just return a list of host names. The power of Herd is that it can filter by any
attribute that providers return. For example, the AWS provider returns instance types or vpc id's.
When using multiple providers, sometimes the names of these attributes can collide. To prevent this,
each provider supports a `prefix` parameter. This prefix will be prepended to all attribute names
coming from this provider. In our example above, the `seveas` provider uses this parameter to add a
`svs:` prefix to all attribute names.

# Caching

Some providers can take quite a while to return data. Herd has built-in caching, implemented as a
separate provider. The cache provider accepts the following parameters:

| Parameter       | Type      | Meaning                                                   | Example | Default            |
|-----------------|-----------|-----------------------------------------------------------|---------|--------------------|
| `lifetime`      | Duration  | How long data should be cached                            | `8h`    | `1h`               |
| `file`          | File path | Where data should be cached, relative to Herd's cache dir |         | `${name}.cache`    |
| `prefix`        | String    | Attribute prefix                                          | `aws:`  |`''` (empty string) |
| `strictloading` | Boolean   | Don't use stale cached data if loading fresh data fails   | `true`  | `false`            |
| `provider`      |           | The config of the source proivder                         |         |                    |

The cache provider is not only useful for caching, it also helps with fault tolerance. If a source
provider fails to load, the cache provider will use cached data. If a providerd partially fails to
load, such as when the consul provider can load data from certain datacenters but not all, the cache
will complement the fresh data with cached data from the failed parts.

# Custom providers

If you have your own inventory database or API, you can plug this into herd in two ways:

- Regularly export the data to an inventory file
- Write a custom provider to query your infrastructure database or API

Infrastructure files are explained a little lower, custom providers [have their own
section](../custom_providers/) in the documentation. It takes only a little bit of Go code to write
a custom provider.

# Built-in provider reference

## SSH known\_hosts

This provider is loaded by default and does not require configuration. It will load
`~/.ssh/known_hosts` and `/etc/ssh/ssh_known_hosts`. If you want it to load other files, you can use
the `files` parameter in the configuration.

| Parameter | Type            | Meaning                       | Example                       | Default                                         |
|-----------|-----------------|-------------------------------|-------------------------------|-------------------------------------------------|
| `files`   | List of strings | The known_hosts files to load | `[/etc/ssh/more_known_hosts]` | `[~/.ssh/known_hosts /etc/ssh/ssh_known_hosts]` |

This provider provides no attributes.

## PuTTY

This provider only works on Windows. It looks in the registry to find host keys and host
configurations. It returns all hosts it finds.

This provider takes no parameters and provides no host attributes.

## Inventory files

As a simple way of integrating not-yet-supported data sources with herd, you can use the plain or
json providers. The plain provider looks for a file named `${datadir}/inventory` and parses that as
a newline-separated list of hostnames. The json provider looks for `${datadir}/inventory` and parses
that as a list of host objects. The format of that file is quite simple, as can be seen in the
following example:

```json
[
  {
    "Name": "host-1.example.com",
    "Address": "10.0.0.1",
    "Attributes": {
      "application": "frontend",
      "owner": "seveas",
    }
  },
  {
    "Name": "host-1.example.com",
    "Address": "10.0.0.1",
    "Attributes": {
      "application": "database",
      "owner": "not-seveas",
    }
  }
]
```

In these host objects, addresses are optional and there can be as many attributes as you want.
Attribute values can be of any type.

These providers take the following parameters:

| Parameter | Type      | Meaning                                                  | Example | Default                         |
|-----------|-----------|----------------------------------------------------------|---------|---------------------------------|
| `prefix`  | String    | Attribute prefix                                         | `aws:`  | `''` (empty string)             |
| `file`    | File path | Where data should be cached, relative to Herd's data dir |         | `inventory` or `inventory.json` |

The plain provider provides no host attributes, the json provider provides the attributes from the
json file.

## HTTP API

The HTTP API provider is not the most useful one on its own, unless your http API happens to
provider inventory data in the exact format the json provider expects. Its use is mostly in
embedding in custom providers that use HTTP APIs.

This provider, and providers that embed it, accept the following parameters:

| Parameter  | Type      | Meaning                                                      | Example                     | Default             |
|------------|-----------|--------------------------------------------------------------|-----------------------------|---------------------|
| `prefix`   | String    | Attribute prefix                                             | `aws:`                      | `''` (empty string) |
| `url`      | String    | The URL for the API                                          | `https://hosts.exanple.com` | (not set)           |
| `username` | String    | HTTP Basic authentication                                    | seveas                      | (not set)           |
| `password` | String    | HTTP Basic authentication                                    | hunter2                     | (not set)           |
| `headers`  | See below | HTTP headers to send in request (e.g. authorization tokens)  |                             |                     |

An example of headers configuration:

```yaml
Providers:
  api:
    provider: http
    Url: https://hosts.example.com/all
    Headers:
      Accept: application/json
      Authorization: Bearer secret-token-123456
```

## Consul

The consul provider finds hosts in all datacenters in consul. If the name consul.service.consul
resolves, or the environment variable `CONSUL_HTTP_ADDR` is set, the consul provider is
automatically configured.

This provider accepts the following parameters:

| Parameter            | Type            | Meaning                                    | Example                         | Default                                                                 |
|----------------------|-----------------|--------------------------------------------|---------------------------------|-------------------------------------------------------------------------|
| `prefix`             | String          | Attribute prefix                           | `consul:`                       | `''` (empty string)                                                     |
| `address`            | String          | Address of the consul server               | http://consul.example/com:8500` | `http://consul.service.consul:8500` or the value of `$CONSUL_HTTP_ADDR` |
| `datadcenters`       | List of strings | Which datacenter to query                  | [dc1 dc2]                       | `[]` (empty list, meaning all datacenters)                              |
| `excludedatacenters` | List of strings | Which datacenters not to query             | [dc3]                           | `[]`                                                                    |
| `maxconcurrency`     | Boolean         | How many datancenters to query in parallel | 10                              | 30                                                                      |

This provider provides the following host attributes

| Attribute                 | Type            | Meaning                        | Example            |
|---------------------------|-----------------|--------------------------------|--------------------|
| `datacenter`              | String          | The consul datacenter          | `dc1`              |
| `node_address`            | String          | Address according to consul    | `10.0.0.1`         |
| `service`                 | List of strings | Services on the host           | `[web mysql smtp]` |
| `service_healthy`         | List of strings | Healthy services on the host   | `[web mysql]`      |
| `service_unhealthy`       | List of strings | Unhealthy services on the host | `[smtp]`           |
| `service:${service_name}` | List of strings | Tags of the service            | `[staging]`        |

## Prometheus

If you monitor your hosts with prometheus, you can tell herd to find the hosts there. You can point
it to the job(s) that do the node_exporter checks for easiest results. This provider uses the HTTP
API provider and accepts the same parameters. The url should point to the prometheus targets api.

This provider accepts the following parameters:

| Parameter  | Type            | Meaning                                   | Example                                         | Default             |
|------------|-----------------|-------------------------------------------|-------------------------------------------------|---------------------|
| `prefix`   | String          | Attribute prefix                          | `prom:`                                         | `''` (empty string) |
| `url`      | String          | The URL for the API                       | `https://prometheus.exanple.com/api/v1/targets` | (not set)           |
| `jobs`     | List of strings | The jobs that have all the hosts to query | `[node]`                                        | (not set)           |
| `username` | String          | HTTP Basic authentication                 | seveas                                          | (not set)           |
| `password` | String          | HTTP Basic authentication                 | hunter2                                         | (not set)           |
| `headers`  | See above       | HTTP headers to send in request           |                                                 |                     |

This provider provides the following host attributes

| Attribute              | Type      | Meaning                                             | Example                               |
|------------------------|-----------|-----------------------------------------------------|---------------------------------------|
| `scrape_pool`          | Strig     | The prometheus scrape pool                          |                                       |
| `scrape_url`           | String    | The URL of the prometheus endpoiint for this target | `http://10.0.0.1:9100/metrics`        |
| `last_scrape`          | Timestamp | The time this host was last scraped                 | `2023-02-12T15:04:45.642518461+01:00` |
| `last_scrape_duration` | Float     | How long the last scrape took                       | `0.050588943`                         |
| `health`               | String    | Host health from prometheus' perspective            | `up`                                  |

## Puppet

The puppet provider queries puppetdb for host facts and adds those as host attributes. It also
populates Herd's internal host key database, complementing keys forund in your known_hosts and/or
PuTTY. It embeds the HTTP providerd and accepts the same parameters, plus its own.

This provider accepts the following parameters:

| Parameter  | Type            | Meaning                                   | Example                                         | Default             |
|------------|-----------------|-------------------------------------------|-------------------------------------------------|---------------------|
| `prefix`   | String          | Attribute prefix                          | `puppet:`                                       | `''` (empty string) |
| `url`      | String          | The URL for the API                       | `http://puppetdb.example.com:8080/`             | (not set)           |
| `facts`    | List of strings | The jobs that have all the hosts to query | `[ssh os dmi drives ruby]`                      | `[ssh os dmi]`      |
| `username` | String          | HTTP Basic authentication                 | seveas                                          | (not set)           |
| `password` | String          | HTTP Basic authentication                 | hunter2                                         | (not set)           |
| `headers`  | See above       | HTTP headers to send in request           |                                                 |                     |

Puppet facts are usually nested, but Herd does not support this nesting, so the lists get flattened,
with keys separated with `:`. For example, the `os` fact usually looks like this:

```console
$ facter os
{
  architecture => "amd64",
  distro => {
    codename => "stretch",
    description => "Debian GNU/Linux 9 (stretch)",
    id => "Debian",
    release => {
      full => "9.13",
      major => "9",
      minor => "13"
    }
  },
  family => "Debian",
  hardware => "x86_64",
  name => "Debian",
  release => {
    full => "9.13",
    major => "9",
    minor => "13"
  },
  selinux => {
    enabled => false
  }
}
```

The attributes that herd then populates are:

```console
$ herd list --no-refresh --template '{{ .Attributes|yaml }}' server-1.example.com | grep '^os:'
os:architecture: amd64
os:distro:codename: stretch
os:distro:description: Debian GNU/Linux 9 (stretch)
os:distro:id: Debian
os:distro:release:full: "9.13"
os:distro:release:major: "9"
os:distro:release:minor: "13"
os:family: Debian
os:hardware: x86_64
os:name: Debian
os:release:full: "9.13"
os:release:major: "9"
os:release:minor: "13"
os:selinux:enabled: false
```

## AWS

If you use AWS EC2, herd can query its API to get your hosts' information.  You will need an access
key id and its associated secret, and the provider exposes a lot of EC2 attributes, such as the
owner id, state or private IP.

This provider accepts the following parameters:

| Parameter          | Type            | Meaning                                         | Example                  | Default                                |
|--------------------|-----------------|-------------------------------------------------|--------------------------|----------------------------------------|
| `prefix`           | String          | Attribute prefix                                | `aws:`                   | `''` (empty string)                    |
| `accesskeyid`      | String          | IAM access key ID                               | `AKAGSDIUYGASD6AD`       | (not set)                              |
| `secretaccesskey`  | String          | IAM secret access key                           | `87wc76g2d6378r2675`     | (not set)                              |
| `partition`        | String          | AWS partition                                   | `aws-us-gov`             | (not set)                              |
| `regions`          | List of strings | Which regions to query                          | `[us-east-1 af-south-1]` | `[]` (empty list, meaning all regions) |
| `excluderegions`   | List of strings | Which regions to exclude                        | `[ap-east-1]`            | `[]`                                   |
| `usepublicaddress` | Boolean         | Whether to use hosts' public address to connect | `true`                   | `false`                                |

This provider also has a magic provider. If `AWS_ACCESS_KEY_ID` or `AWS_ACCESS_KEY` is set in the
environment, it is used as IAM access key ID, and if `AWS_SECRET_ACCESS_KEY` or `AWS_SECRET_KEY` is
set in the environment, it is uses as IAM secret access key.

This provider provides the following host attributes. All host tags are also added as attributes.

| Attribute                 | Type    | Meaning                                             | Example                               |
|---------------------------|---------|-----------------------------------------------------|---------------------------------------|
| `architecture`            | String  | Hardware architecture of the VM                     | `x86_64`                              |
| `availability_zone`       | String  | Theavailability zone a host is in                   | `us-east-1a`                          |
| `core_count`              | Integer | How many cores a VM has                             | `1`                                   |
| `hypervisor`              | String  | The type of underlying hypervisor                   | `xen`                                 |
| `image_id`                | String  | The AMI used to create the host                     | `ami-07b038657b94fb025`               |
| `instance_id`             | String  | AWS' internal instance identifier                   | `i-0c878318779a1cae1`                 |
| `instance_type`           | String  | EC2 instance type                                   | `m5.large`                            |
| `launch_time`             | String  | When the server was created                         | `2020-12-01T13:04:55Z`                |
| `owner_id`                | String  | Numerical identifier of the owner                   | `493178567742`                        |
| `placement_group`         |         |                                                     |                                       |
| `private_dns_name`        | String  | AWS automatic dns name                              | `ip-10.0.0.1.ec2.internal`            |
| `private_ip`              | String  | Private IP address                                  | `10.0.0.1`                            |
| `public_dns_name`         |         |                                                     |                                       |
| `public_ip`               | String  | Public IP address                                   | `54.162.196.34`                       |
| `reservation_id`          | String  | The reservation the host belongs to                 | `r-87f4a7ee5189fc7b4`                 |
| `root_device_name`        | String  | The disk device for `/`                             | `/dev/sda1`                           |
| `root_device_type`        | String  | The storage type used for `/`                       | `ebs`                                 |
| `state_code`              | Integer | A numerical code corresponding to the machine state | `16`                                  |
| `state_transition_reason` |         |                                                     |                                       |
| `state`                   | String  | The state of the host                               | `running`                             |
| `subnet_id`               | String  | Which subnet the host is in                         | `subnet-7b3e09d1`                     |
| `tenancy`                 |         |                                                     |                                       |
| `threads_per_core`        | Integer | The number of threads per core (hyperthreading)     | `2`                                   |
| `virtualization_type`     | String  | What virtualization type is used                    | `hvm`                                 |
| `vpc_id`                  | String  | Which virtual private cloud the host is in          | `vpc-ff1265491a`                      |

## Azure

This provider retrieves host data from azure. You will need to provide a subscription, tenant and
client identifiers and a set of credentials. If you do not provide credentials, herd will try to
read them from the environment or from MSI.

This provider accepts the following parameters:

| Parameter             | Type   | Meaning                                         | Example                                | Default             |
|-----------------------|--------|-------------------------------------------------|----------------------------------------|---------------------|
| `prefix`              | String | Attribute prefix                                | `azure:`                               | `''` (empty string) |
| `subscription`        | String | Azure subscription id                           | `93692f54-4e0a-484c-94e0-9dd370d3aa73` | (not set)           |
| `tenantid`            | String | Azure tenant id                                 | `0fd787b9-b56c-4c9d-b97c-47c9a25e6108` | (not set)           |
| `clientid`            | String | Azure client id                                 | `b5f906b6-67e8-41b3-9234-4fc14302538c` | (not set)           |
| `clientsecret`        | String | Secret to use for authentication                | `G967TIGgiG8G%fffTfyutasd78sad`        | (not set)           |
| `certificatepath`     | String | Authentication certificate                      | `/etc/azure/cert.p12`                  | (not set)           |
| `certificatepassword` | String | Password for the encrypted certificate          | `hunter2`                              | (not set)           |
| `username`            | String | AAD username                                    | `seveas`                               | (not set)           |
| `password`            | String | AAD password                                    | `hunter2`                              | (not set)           |

Only one set of credentials (client secret, certificate or username/password) is needed.

This provider provides the following host attributes. All host tags are also added as attributes.

| Attribute           | Type   | Meaning                           | Example                                                                                                                                       |
|---------------------|--------|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `VMSize`            | String | The size of the virtual machine   | `Standard_D4s_v4`                                                                                                                             |
| `ProvisioningState` | String | The state of server provisioning  | `Succeeded`                                                                                                                                   |
| `VMID`              | String | Unique identifier for the VM      | /subscriptions/e3113d22-2af2-4f34-ae21-4bc378baa41f/resourceGroups/MYGROUP/providers/Microsoft.Compute/virtualMachines/server-01.example.com` |
| `Location`          | String | Geographic location of the server | `eastus`                                                                                                                                      |

## Google cloud

This provider accepts the following parameters:

| Parameter          | Type            | Meaning                                         | Example                    | Default                              |
|--------------------|-----------------|-------------------------------------------------|----------------------------|--------------------------------------|
| `prefix`           | String          | Attribute prefix                                | `google:`                  | `''` (empty string)                  |
| `key`              | String          | Path to the API key                             | `/etc/google/api_key.json` | (not set)                            |
| `project`          | String          | Which project to query                          | `example-project`          | (not set)                            |
| `zones`            | List of strings | Which zones to query                            | `[us-east1]`               | `[]` (empty list, meaning all zones) |
| `usepublicaddress` | Boolean         | Whether to use hosts' public address to connect | `true`                     | `false`                              |

This provider provides the following host attributes. All host labels are also added as attributes.

| Attribute          | Type    | Meaning                              | Example                                                                                                      |
|--------------------|---------|--------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `can_ip_forward`   | Boolean |                                      | `true`                                                                                                       |
| `cpuplatform`      | String  | The CPU variant used for this host   | `Intel Haswell`                                                                                              |
| `description`      | String  | A description of this host           |                                                                                                              |
| `id`               | Integer | Numerical identifier for this host   | `479572295727027251`                                                                                         |
| `fingerprint`      | String  |                                      |                                                                                                              |
| `kind`             | String  |                                      | `compute#instance`                                                                                           |
| `labelfingerprint` | String  |                                      |                                                                                                              |
| `last_start`       | String  | When the host was last started       | `2022-11-12T19:45:12.943-07:00`                                                                              |
| `last_stop`        | String  | When the host was last stopped       |                                                                                                              |
| `last_suspend`     | String  | When the host was last suspended     |                                                                                                              |
| `machinetype`      | String  | The type of virtual machine          | `https://www.googleapis.com/compute/v1/projects/example-project/zones/us-east1-b/machineTypes/n1-standard-4` |
| `mincpuplatform`   | String  |                                      |                                                                                                              |
| `instancename`     | String  | The name of the instance             | `server-1.example.com`                                                                                       |
| `startrestricted`  | Boolean |                                      | `false`                                                                                                      |
| `statusmessage`    | String  |                                      |                                                                                                              |
| `region`           | String  |                                      | `https://www.googleapis.com/compute/v1/projects/example-project/regions/us-east1`                            |
| `zone`             | String  | The availability zone the host is in | `us-east1-b`                                                                                                 |


## TransIP

For VPS'es hosted at TransIP, herd can find information with this provider.  You will need to enable
API access for your account in the API controlpanel and generate an API key for herd to use.

This provider accepts the following parameters:

| Parameter | Type   | Meaning                     | Example                | Default             |
|-----------|--------|-----------------------------|------------------------|---------------------|
| `prefix`  | String | Attribute prefix            | `transip:`             | `''` (empty string) |
| `user`    | String | The user to authenticate as | `seveas`               | (not set)           |
| `key`     | String | Path to the API key         | `/etc/transip/api_key` | (not set)           |

This provider provides the following host attributes:

| Attribute          | Type            | Meaning                                         | Example                                |
|--------------------|-----------------|-------------------------------------------------|----------------------------------------|
| `availabilityzone` | String          | Which availability zone the host is in          | `ams0`                                 |
| `cpus`             | Integer         | How many CPU's does the host have               | `2`                                    |
| `currentsnapshots` | Integer         | How many snapshots exist for this host          | `0`                                    |
| `description`      | String          | A description of the host                       | `example-server`                       |
| `disksize`         | Integer         | The disk size in KB                             | `157286400`                            |
| `ipaddress`        | String          | IP address of the host                          | `141.138.139.206`                      |
| `isblocked`        | Boolean         |                                                 | `false`                                |
| `iscustomerlocked` | Boolean         |                                                 | `false`                                |
| `islocked`         | Boolean         |                                                 | `false`                                |
| `macaddress`       | String          | Mac address of the host                         | `52:12:46:a3:7c:11`                    |
| `maxsnapshots`     | Integer         | How many snapshots may exist for this host      | `1`                                    |
| `memorysize`       | Integer         | Memory size in KB                               | `4194304`                              |
| `operatingsystem`  | String          | Which operating system was originally installed | `Ubuntu 22.04`                         |
| `productname`      | String          | Which product is this host based on             | `unknown`                              |
| `status`           | String          | Status of the host                              | `running`                              |
| `tags`             | List of strings | All the tags you added to the host              | `[example]`                            |
| `uuid`             | String          | A unique identifier                             | `d6480e70-c895-42b2-8713-03a9ebe87e6b` |

## Tailscale

Herd can autodetect all the hosts in your tailnet. This provider is not yet included by default, as
it does not build on Windows at this time. It can be installed as a provider plugin.

```console
$ go install github.com/seveas/herd/cmd/herd-provider-tailscale@latest
```

This provider accepts the following parameters

| Parameter | Type   | Meaning                                           | Example          | Default             |
|-----------|--------|---------------------------------------------------|------------------|---------------------|
| `prefix`  | String | Attribute prefix                                  | `ts:`            | `''` (empty string) |
| `domain`  | String | Override the domainname for hosts on your tailnet | `ts.example.com` | (not set)           |


This provider provides the following host attributes:

| Attribute       | Type            | Meaning                                         | Example                                                                              |
|-----------------|-----------------|-------------------------------------------------|--------------------------------------------------------------------------------------|
| `active`        | Boolean         | Whether the host is online                      | `false`                                                                              |
| `addrs`         |                 |                                                 |                                                                                      |
| `apiurl`        | List of strings | The API urls for the host                       | `[http://100.66.63.21:39443 http://[fd7a:115c:a1e0:ab12:4843:c596:6243:4a15]:39443]` |
| `capabilities`  |                 |                                                 |                                                                                      |
| `created`       | String          | When tailscale was first installed on the host  | `2022-12-01T04:55:12.4522578Z`                                                       |
| `curaddr`       | String          |                                                 |                                                                                      |
| `dnsname`       | String          | Tailscal DNS name for the host                  | `server-1.tail8b340.ts.net`                                                          |
| `exitnode`      | Boolean         | Whether the host functions as an exit node      | `false`                                                                              |
| `hostname`      | String          | The hostname of the host                        | `server-1`                                                                           |
| `id`            | String          | An opaque unique identifier                     | `bGNogYu2ryb9nA`                                                                     |
| `inengine`      | Boolean         |                                                 | `true`                                                                               |
| `inmagicsock`   | Boolean         |                                                 | `true`                                                                               |
| `innetworkmap`  | Boolean         |                                                 | `true`                                                                               |
| `keepalive`     | Boolean         |                                                 | `false`                                                                              |
| `lasthandhake`  | String          |                                                 |                                                                                      |
| `lastseen`      | String          |                                                 |                                                                                      |
| `lastwrite`     | String          |                                                 |                                                                                      |
| `os`            | String          | The os of the host                              | `linux`                                                                              |
| `publickey`     | String          | The publick key of the host                     | `nodekey:b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c`           |
| `relay`         | String          | Which tailscale relay is used                   | `ams`                                                                                |
| `rxbytes`       | Integer         |
| `shareenode`    | Boolean         |
| `tailscale_ips` | List of strings | The Tailscale VPN IP addresses of the host      | [100.66.63.21 fd7a:115c:a1e0:ab12:4843:c596:6243:4a15]                               |
| `txbytes`       | Integer         |
| `userid`        | Integer         | The ID of the user this host belongs to         | 24365                                                                                |
