---
Title: Configuration reference
Weight: 5
---
# Configuration sources

Herd uses the [viper](https://github.com/spf13/viper) library for configuration and command line
parsing. This means that each parameter can be set in 3 ways: in a configuration file, in the
environment or as a command-line argument.

The only exception to this are the providers, which can only be configured in the configuration
file. A per-provider reference can be found [elsewhere in the documentation](/host_discovery/#built-in-provider-reference)

The location of the configuration file depends on the operating system you are on. `herd -h` will
show you where the configuration files live on your machine.

# Configuration file syntax

The configuration file is a yaml file, all variables are at the top level. Here's an example that
sets some timeout defaults and adds one provider:

```yaml
Timeout:     10m
HostTimeout: 1m

Providers:
  tailscale:
    provider: tailscale
    prefix: "ts:"
```

# Configuration variables

Variables are named roughly the same in all three places, but capitalization differs. The yaml
configuration is case-insensitive, environment variables are all uppercase and separate words with
underscores and command-line parameters use dashes instead of underscores and separate words with
dashes. So the `HostTimeout` variable in yaml is spelled `HERD_HOST_TIMEOUT` in the environment and
corresponds to the `--host-timeout` environment variable.

### General behaviour

| Variable   | Type    | Meaning                                                                                                     |
|------------|---------|-------------------------------------------------------------------------------------------------------------|
| `LogLevel` | String  | The log level to use. Default is `info`, supports the standard log levels                                   |
| `NoPager`  | Boolean | Herd auto-starts a pager when the output it would show is bigger than your screen, this flag inhibits that. |
| `NoColor`  | Boolean | Herd uses color in its output to indicate error/success and for other purposes. This flag inhibits that.    |

### Data loading

| Variable           | Type     | Meaning                                                                                                |
|--------------------|----------|--------------------------------------------------------------------------------------------------------|
| `NoMagicProviders` | Boolean  | Do not try to autodetect hosts based on data in your environment or known locations on your filesystem |
| `Refresh`          | Boolean  | Force a refresh of cached data                                                                         |
| `NoRefresh`        | Boolean  | Do not try to refresh cached data                                                                      |
| `StrictLoading`    | Boolean  | If one or more providers fail to load data, abort before running commands or showing hosts             |
| `LoadTimeout`      | Duration | The time providers can take to load data                                                               |

### Host list output

| Variable     | Type            | Meaning                                                      |
|--------------|-----------------|--------------------------------------------------------------|
| `Sort`       | List of strings | How to sort hosts before showing them                        |
| `Attributes` | List of strings | Which attributes to include in the output                    |
| `Header`     | Boolean         | Whether to show a header with the attribute names            |
| `CSV`        | Boolean         | Output csv data instead of nicely formatted data             |
| `Separator`  | String          | Use a non-standard separator for CSV data                    |
| `Template`   | String          | A text/template template to completely customize list output |
| `Count`      | List of strings | Show counts of hosts by attribute values                     |
| `Group`      | String          | Group hosts by attribute values                              |

### Command running and output

| Variable          | Type            | Meaning                                                                                                                               |
|-------------------|-----------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `Parallel`        | Integer         | Limit the amount of hosts that commands run in parallel on                                                                            |
| `Splay`           | Duration        | Wait a random duration up to the specified argument before connecting to each host to spread command starts                           |
| `ConnectTimeout`  | Duration        | Maximum time allowed for connection set up                                                                                            |
| `SshAgentTimeout` | Duration        | Maximum time allowed for the SSH agent to respond when detecting SSH agent pipelining                                                 |
| `HostTimeout`     | Duration        | Maximum time, including connection set up time, a command may take per host                                                           |
| `Timeout`         | Duration        | Total timeout for a parallel invocation. Any command not finished will be terminated, any command not started yet will not be started |
| `Output`          | String          | The output format to use, one of `all`, `per-host`, `inline` and `tail`                                                               |
| `Sort`            | List of strings | How to sort hosts before showing their results, not used for `tail` and `per-host` output                                             |
| `Timestamp`       | Boolean         | Show a timestamp in front of command output in tail mode                                                                              |

# Theming

All colors used by herd can be customized with a `Colors` section in the configuration. The colors
it understands are the standard ansi colors with some attributes. For a full specification, see
https://github.com/mgutz/ansi

```yaml
Colors:
  LogDebug:   black+h
  LogInfo:    default
  LogWarn:    yellow
  LogError:   red+b
  Command:    cyan
  Summary:    black+h
  Provider:   green
  HostStdout: default
  HostStderr: yellow
  HostOK:     green
  HostFail:   yellow
  HostError:  red
  HostCancel: black+h
```

# Other environment variables

- Herd will automatically start a pager when its output spans more than one screen. This defaults to `less`, but can be overridden with the `PAGER` environment variable.
- `SSH_AUTH_SOCK` and `SSH_CONNECTION` are used to connect to an ssh agent
- `XDG_CONFIG_HOME`, `XDG_DATA_HOME` and `XDG_CONFIG_HOME` are used if available to find configuration and data locations
- `HOME` and `PATH` are used for finding configuration files and provider plugins

Some providers also use environment variables, this is documented in the provider documentation.

# OpenSSH configuration

To avoid having to duplicate your SSH configuration, Herd will respect some of the configuration
parameters set in `~/.ssh/config`. At the moment only 5 things are respected:

- `User`, which defaults to your local username
- `Port`, which defaults to 22
- `IdentityFile`, to limit which keys from your agent will be used for a host
- `StrictHostKeyChecking`, which defaults to `accept-new` for Herd
- `VerifyHostKeyDns` to enable checking host keys in DNS. Herd does _not_ do DNSSEC verification.
