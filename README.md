Herd - Lightning-fast parallel ssh client
=============================================

Herd is a massively parallel ssh client that will replace all your hacky for loops, pssh, xargs
and gnu parallel oneliners with something that is faster, more flexible and simply more awesome!

Herd can find your hosts in you `known_hosts` file, in consul and many other places. Using a
powerful query syntax, you select the hosts you want and run any command in parallel on all of them.
Output can be streamed as it happens, or shown at the end for easy inspection. On top of that,
Herd is chock full of helpful features:

- It can query your inventory systems and display information and statistics about your hosts
- Timeouts, parallelism and delays are all configurable for very flexible ways of running commands
- A variety of output modes show exactly what you need. Be it timestamps at the start of all lines,
  output separated by host or one big river of output. It can do all that and more.
- Output is also stored in a machine-readable way so you can post process it as much as you want
- An interactive shell is available that even lets you create host lists based on the result of
  running commands on hosts

Examples
--------

Here are some examples to get you going. The full documentation can be found on
https://herd.seveas.net

- List all hosts in a specific domain

  `herd list *.example.com`

- Run a command on all of them

  `herd run *.example.com -- uptime`

- Run a long-running command, with output appearing immediately and reduced parallelism

  `herd run *.example.com --parallel 10 --host-timeout 5m --timeout 1h -- sudo puppet agent -t`

- Show some information about hosts in an inventory system

  `herd list site=site-5 --attributes ip_address,os,memory`

- Rolling restart of a service registered in consul

  `herd run --splay 1m --parallel 1 consul_service=smtpd -- sudo systemctl restart postfix`

Installing
----------
There are pre-built binaries on the [download page](https://herd.seveas.net/download/). If
you prefer to install from source, install [go](https://golang.org), clone this repository and
just run `make`.
