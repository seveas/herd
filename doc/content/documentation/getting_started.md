---
Title: Getting started
Weight: 1
---

Welcome! The documentation on these pages will tell you all about how herd works, from installing it
to customizing it by writing your own plugins. But let's start off easy and first just get it
installed and run some basic commands.

# Installing

Installing herd is pretty easy: it's a single binary that you can download and put on your `$PATH`
or `%PATH%`. Binaries for Linux, Windows BSD and mac can be found on the [download page](/download).
Once installed, let's confirm that it works:

```console
$ herd --version
herd version 0.14.0
```

Great, basic functionality confirmed! Now let's check the help output to see what we can do.

```console
$ herd --help

Replace your ssh for loops with a tool that
- can find hosts for you
- can handle thousands of hosts in parallel
- does not fork a command for every host
- stores all its history, including output, for you to reuse
- can run interactively!

Usage:
  herd [command]

Examples:
  herd run '*' os=Debian -- dpkg -l bash
  herd interactive *vpn-gateway*

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  interactive Interactive shell for running commands on a set of hosts
  keyscan     Scan ssh keys and output them in known_hosts format, similar to ssh-keyscan
  list        Query your datasources and list hosts matching globs and filters
  ping        Connect to hosts and check if they are alive
  run         Run a single command on a set of hosts
  run-script  Run a script on a set of hosts

Flags:
      --connect-timeout duration     Per-host ssh connect timeout (default 3s)
  -h, --help                         help for herd
      --host-timeout duration        Per-host timeout for commands (default 10s)
      --load-timeout duration        Timeout for loading host data from providers (default 30s)
  -l, --loglevel string              Log level (default "INFO")
      --no-color                     Disable the use of the colors in the output
      --no-magic-providers           Do not use magic autodiscovery, only explicitly configured providers
      --no-pager                     Disable the use of the pager
      --no-refresh                   Do not try to refresh cached data
  -o, --output string                When to print command output (all at once, per host or per line) (default "all")
  -p, --parallel int                 Maximum number of hosts to run on in parallel
      --profile string               Write profiling and tracing data to files starting with this name
      --refresh                      Force caches to be refreshed
  -s, --sort strings                 Sort hosts by these fields before running commands (default [name])
      --splay duration               Wait a random duration up to this argument before and between each host
      --ssh-agent-timeout duration   SSH agent timeout when checking functionality (default 50ms)
      --strict-loading               Fail if any provider fails to load data
  -t, --timeout duration             Global timeout for commands (default 1m0s)
      --timestamp                    In tail mode, prefix each line with the current time
  -v, --version                      version for herd

Use "herd [command] --help" for more information about a command.


Configuration: /Users/seveas/Library/Preferences/herd/config.yaml, /etc/herd/config.yaml
Datadir: /Users/seveas/Library/Application Support/herd
History: /Users/seveas/Library/Application Support/herd/history
Cache: /Users/seveas/Library/Caches/herd
Providers: aws,azure,cache,consul,google,http,json,known_hosts,plain,plugin,prometheus,puppet,transip
```

That is a _lot_ of information. Let's take it step by step and start with the basics: finding hosts.

# Finding hosts

Contrary to other SSH clients, herd really doesn't want you to tell it a hostname to connect to. It
wants to find hosts, filter them by some criteria and connect to a set of hosts in parallel. So how
does it find hosts?

The simplest way is by looking in your `~/.ssh/known_hosts` or in PuTTY's host key cache. For this
to work, you do need to disable known_hosts hashing in `~/.ssh/config`

```
HashKnownHosts no
```

Once you disable that, any host you ssh to can be found by herd as well when you enable known_hosts
discovery in the herd configuration:

```yaml
Providers:
    known_hosts:
```

This is not enabled by default, because in larger environments `known_hosts` tends to accumulate
decommissioned hosts, and herd will try to connect to them.

Let's see an example:

```console
$ herd list | head -n 5
server-1.example.com
server-2.example.com
server-3.example.com
server-4.example.com
server-5.example.com
```

Great! We found servers. Lets see if we can connect to them. And we'll show our first trick: we can
select servers using a regular expression on their name.

```console
$ herd ping name=~server-..example.com
{{<ansi green>}}server-1.example.com{{</ansi>}} connection successful in 1.596s
{{<ansi red>}}server-2.example.com Timed out while connecting to server after 3s{{</ansi>}}
{{<ansi green>}}server-3.example.com{{</ansi>}} connection successful in 737ms
{{<ansi green>}}server-4.example.com{{</ansi>}} connection successful in 746ms
{{<ansi green>}}server-5.example.com{{</ansi>}} connection successful in 669ms
{{<ansi green>}}server-6.example.com{{</ansi>}} connection successful in 847ms
{{<ansi green>}}server-7.example.com{{</ansi>}} connection successful in 737ms
```

Looks like that worked! And we see another few features: herd uses color to indicate success and
failure, and despite errors on a host, it still does what it needs to for other hosts.

# SSH agent setup

If you didn't get any output, it's more than likely that your SSH agent is not running. Herd can
only authenticate using an SSH agent, so you must have one running. Modern linux distributions, and
even MacOS automatically start one when you log in. On windows, PuTTY's Pageant can be used as well.
This is what it looks like when your agent is running and has keys:

```console
$ ssh-add -l
2048 SHA256:rQeIwveJHKV2y1Oqs8AR732TynAne5kFTtg6tybbyT4 /Users/seveas/.ssh/id_rsa (RSA)
```

If you get this instead:

```console
$ ssh-add -l
Could not open a connection to your authentication agent.
```

Then your SSH agent is not running. You can start one in the current shell with:

```console
$ eval $(ssh-agent)
Agent pid 55198
```

If alternatively, you get this:

```console
$ ssh-add -l
The agent has no identities.
```

Then your agent is running but has no keys. Use `ssh-add` (or `ssh-add -k` on a mac) to load the
keys.

# Your first command

Now that we can connect, it's time to run our first command:

```console
$ herd run name=~server-..example.com -- date
7 done, 7 ok, 0 fail, 0 error in 3s
{{<ansi green>}}server-1.example.com  completed successfully after 1s{{</ansi>}}
    Sun Feb 12 14:01:06 EST 2023
{{<ansi green>}}server-2.example.com  completed successfully after 1s{{</ansi>}}
    Sun Feb 12 14:01:06 EST 2023
{{<ansi green>}}server-3.example.com  completed successfully after 1s{{</ansi>}}
    Sun Feb 12 14:01:06 EST 2023
{{<ansi green>}}server-4.example.com  completed successfully after 2s{{</ansi>}}
    Sun Feb 12 14:01:07 EST 2023
{{<ansi green>}}server-5.example.com  completed successfully after 2s{{</ansi>}}
    Sun Feb 12 14:01:07 EST 2023
{{<ansi green>}}server-6.example.com  completed successfully after 2s{{</ansi>}}
    Sun Feb 12 14:01:07 EST 2023
{{<ansi green>}}server-7.example.com  completed successfully after 1s{{</ansi>}}
    Sun Feb 12 14:01:06 EST 2023
History saved to /Users/seveas/Library/Application Support/herd/history/2023-02-12_230105.json
```

And here we see a few more features: the string `--` is used to separate the host query and the command, herd
shows how long a command (including connection set up time) took, and it stored the history of the
command. More importantly, it stored the output as well, so you can easily look at a command's
output again without running it.

That's it! You're now ready to dive into the real meat of herd and set up host discovery and discover
all the ways you can use herd to find information about your hosts and run commands on them.
