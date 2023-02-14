---
Title: Custom host providers
Weight: 6
---
If your host data lives in a data source not covered by the default providers, you can create a
custom provider that would integrate seamlessly with herd.  This takes only a few steps and a bit of
boilerplate code on top of the actual code to fetch your host data.

In this document we go over what it would take to create an example provider named dibbler. The text
assumes you know how to work with git, github and the go programming language.

# Code organization

If you plan to contribute your custom provider to herd, all you need to do is create the actual
provider. Create a new branch in the herd repository, and in the directory `provider/dibbler/` you
create the file `provider.go` which will hold your code.

If your provider will not be included in herd, it needs to be built as an external plugin. For this,
you create a new repository on GitHub named herd-provider-dibbler with the following files:

- `provider/dibbler/provider.go`
- `cmd/herd-provider-dibbler/main.go`

The first file is for your provider code, the second for the executable wrapper. After creating
them, run `go mod init github.com/your-username/herd-provider-dibbler` to initialize the go module.

# The provider

Every provider must implement the
[`HostProvider`](https://pkg.go.dev/github.com/seveas/herd#HostProvider) interface. Unfortunately,
the API documentation is rather barren at this point, so an example will have to suffice for now.

```go
package dibbler

import (
	"context"

	"github.com/seveas/herd"
	"github.com/spf13/viper"
)

func init() {
	// This is how we tell Herd about your provider. We pass it a name and two initializers: one for
	// explicit initialization and one for a magic provider. Our simple provider does not do magic.
	herd.RegisterProvider("dibbler", newProvider, nil)
}

type dibblerProvider struct {
	name   string
	config struct {
		// A Prefix is required here
		Prefix string
	}
}

func newProvider(name string) herd.HostProvider {
	return &dibblerProvider{name: name}
}

func (p *dibblerProvider) Name() string {
	return p.name
}

func (p *dibblerProvider) Prefix() string {
	return p.config.Prefix
}

// If your config has more entries, e.g. credentials or an API url, this function must return
// whether the providers can be used interchangably.
func (p *dibblerProvider) Equivalent(o herd.HostProvider) bool {
	return true
}

func (p *dibblerProvider) ParseViper(v *viper.Viper) error {
	return v.Unmarshal(&p.config)
}

// This is the main function, it needs to return a set of hosts and/or an error
func (p *dibblerProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	// Tell herd that we've started loading
	lm(p.name, false, nil)
    // This is where your discovery code will go, for now we return something bogus
	hosts := herd.NewHostSet()
	host := herd.NewHost("server-01.example.com", "10.0.0.1", herd.HostAttributes{"app": "web", "env", "staging"})
	hosts.AddHost(host)
    return hosts, nil
}
```

# Leveraging the HTTP provider

If your provider makes API calls to fetch its data, you can use the HTTP provider to do the actual
HTTP fetching and your code can focus on turning the returned data into a hostset. Here's an example
to illustrate this option:

```go
package dibbler

import (
	"context"

	"github.com/seveas/herd"
	"github.com/seveas/herd/provider/http"
	"github.com/spf13/viper"
)

func init() {
	herd.RegisterProvider("dibbler", newProvider, nil)
}

type dibblerProvider struct {
	name   string
	hp     *http.HttpProvider
	config struct {
		// We don't need any config for the HTTP settings, that's handled by the HTTP provider
		Prefix string
	}
}

func newProvider(name string) herd.HostProvider {
	return &dibblerProvider{name: name, hp: http.NewProvider(name).(*http.HttpProvider)}
}

func (p *dibblerProvider) Name() string {
	return p.name
}

func (p *dibblerProvider) Prefix() string {
	return p.config.Prefix
}

func (p *dibblerProvider) Equivalent(o herd.HostProvider) bool {
	op := o.(*dibblerProvider)
	// We're equivalent if the embedded http providers are
	return p.hp.Equivalent(op.hp)
}

func (p *dibblerProvider) ParseViper(v *viper.Viper) error {
	// First we let the HTTP provider parse things
	if err := p.hp.ParseViper(v); err != nil {
		return err
	}
	return v.Unmarshal(&p.config)
}

func (p *dibblerProvider) Load(ctx context.Context, lm herd.LoadingMessage) (*herd.HostSet, error) {
	lm(p.name, false, nil)
	// We use the embedded HTTP provider to do the fetching
	data, err := p.hp.Fetch(ctx)
	if err != nil {
		return nil, err
	}
    // This is where you parse the returned data and create a HostSet, for now we return something bogus
	hosts := herd.NewHostSet()
	host := herd.NewHost("server-01.example.com", "10.0.0.1", herd.HostAttributes{"app": "web", "env", "staging"})
	hosts.AddHost(host)
    return hosts, nil
}
```

As you can see, config parsing and equivalence testing are delegated to the http provider, and
fetching the data is as simple as calling `p.hp.Fetch(ctx)`.

# The plugin wrapper

If you work in a separate repository, your main.go must contain the following.  This creates a
wrapper around your provider that will allow it to be used as a go plugin. End users don't notice
the difference, except that your provider is shipped as a separate binary.

```go
package main

import (
	// Import your provider
	_ "github.com/your-github-account/herd-provider-dibbler/provider/dibbler"

	// And the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("dibbler"); err != nil {
		panic(err)
	}
}
```

# Adding the provider to herd

How to make your provider available depends on whether you implemented it as part of the herd
repository and on how you want to distribute it. If you've added it as part of the herd repository,
just run `make` to build a new herd binary and you're done.

If you're provider will not be part of herd itself, but a standalone repository, you have two
options:

- Install the plugin as a separate binary in your `$PATH`. This can be done with `go install
  github.com/your-username/your-plugin-name/cmd/your-plugin-name`. Herd will be able to find this
  plugin and use it as if it were built-in.
- Build a custom version of herd. This is a good option if you use your own distribution mechanism
  to distribute herd in your environment anyway. To do this, create a file named
  `cmd/herd/custom_providers.go` in your copy of the herd repository and use it to import your
  provider.

# Configuring herd for your plugin provider

You configure herd the same as how you would for any other provider. Here is an example
configuration with some custom paramenters and caching.

```yaml
Providers:
  dibbler:
    provider: cache
    lifetime: 8h
    source:
      provider: dibbler
	  checksum: d1a253884a92a9bbec381bf4c36bebad46299a59b79ac601b726fa8fea6e2877
      location: ankh-morpork
      companion: gaspode
```

Note the `checksum` parameter. If provided, the plugin machinery will verify that the plugin binary
matches this sha256 checksum and will not execute the plugin if it does not match.
