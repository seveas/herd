---
Title: Querying for hosts
Weight: 3
---

Herd has a powerful search function that can search hosts based on any attribute and give you
immediate insight into the state of your server farm. In the [host discovery](../host_discovery/)
documentation you can read all about how you can set up Herd to find your hosts, in this document we
will use that information to select hosts and show information about them.

All the queries show here work for all Herd subcommands, not just for `herd list`. The host output
formatting and statistical functions only work with `herd list`.

## Selecting by name

The simplest way to select hosts is to select them by name. This supports globbing, so we can find
hosts with similar names with a single command:

```console
$ herd list server-*.example.com
server-1.example.com
server-2.example.com
server-3.example.com
server-4.example.com
server-5.example.com
server-6.example.com
server-7.example.com
```

If you don't use a glob, this even works for hosts that are not in your inventory at all. Herd
will use names from the command line, as long as they resolve

```console
$ herd list www.google.com
www.google.com
```

A special type of glob-like match is when the name starts with `file:`. This loads a list of hosts
from a file. Because no matter how powerful the queries are that you learn later on in this file,
sometimes you just have a list of hosts to work on.

## Multiple sets of hosts

You can also specify multiple sets of hosts this way:

```console
$ herd list web-*.example.com + db-*.example.com
db-1.example.com
db-2.example.com
db-3.example.com
web-1.example.com
web-2.example.com
web-3.example.com
```

A `+` is required between these globs, as the full syntax is `<glob> <attribute filters> [<+|->
<glob> <attribute filters>]...`. And that `-` is not a typo, you can do set arithmetic with this.
For example, this invocation finds all example.com and example.org servers, except web servers:

```console
$ herd list *.example.com + *.example.org - web*
db-1.example.com
db-1.example.org
db-2.example.com
db-2.example.org
mail.example.com
mail.example.org
```

## Attribute matching

Hostname matching is fun and all, but the real search power is in working with host attributes. A
lot of the built-in providers add attributes to hosts, and if you use your own providers or
inventory, you can find a lot of information about hosts. For example, if you use the puppet
provider, you can easily find all the hosts that run a certain os:

```console
$ herd list os:distro:codename=stretch
db-1.example.com
web-3.example.com
```

You can filter on as many attributes as you want, all attributes must match. You can also combine
globs and filters

```console
$ herd list web* os:distro:codename=bullseye VMSize=Standard_D4s_v4
web-1.example.org
db-3.example.org
```

You can also filter for inequality or do regular expression matches on attributes

| Operator     | Meaning                           | Example                       |
|--------------|-----------------------------------|-------------------------------|
| `=` or `==`  | Equality test                     | `os:distro:codename=bullseye` |
| `!=` or `~=` | Inequality test                   | `os:distro:id!=debian`        |
| `=~`         | Regular expression match          | `availability_zone=~us`       |
| `!~`         | Regular expression does not match | `availability_zone!~us`       |

Combined with set arithmetic, this can lead to queries that really give you only the hosts you are
looking for.

### Attribute types

Not all host attributes are strings, but command line parameters are strings. To be able to match on
no-string attributes, herd will convert command line arguments to numbers, booleans or nil if
needed. If you use the interactive mode or scripts, this does not apply as proper types are
supported in scripts. So on the command line filters like `core_count=1` or `isblocked=true` do what
you want.

### Multi-valued attributes

Some attributes may have more than one value, such as the `service` attribute provided by the consul
provider which is a list of strings. You can use these in filters, but the operators take on a
different meaning

| Operator     | Meaning                           | Example        | Matches         | Does not match  |
|--------------|-----------------------------------|----------------|-----------------|-----------------|
| `=` or `==`  | At least one value is exactly (∃) | `service=web`  | `[db mail web]` | `[db mail]`     |
| `!=` or `~=` | No value is exactly           (∄) | `service!=web` | `[db mail]`     | `[db mail web]` |
| `=~`         | At least one value matches        | `service=~web` | `[darkweb]`     | `[weeble]`      |
| `!~`         | No value matches                  | `service!~web` | `[weeble]`      | `[darkweb]`     |

### Built-in attributes

Host attributes come from the host providers you use, but there are also some built-in attributes
that are always available:

| Attribute       | Type            | Meaning                                                                                                                  |
|-----------------|-----------------|--------------------------------------------------------------------------------------------------------------------------|
| `name`          | String          | The name of the host                                                                                                     |
| `domainname`    | String          | The domainname of the host                                                                                               |
| `random`        | Integer         | A not-really-random number for stable not-really-random sorting                                                          |
| `stdout`        | String          | The output of the last command in interactive/scripted mode                                                              |
| `stderr`        | String          | The output of the last command in interactive/scripted mode                                                              |
| `exitstatus`    | Integer         | The exit status of the last command in interactive/scripted mode, `-1` when there wan error establishing a connection    |
| `err`           | Error           | The error that occurred during the last command in interactive/scripted mode. Note that a non-zero exit is also an error |
| `herd_provider` | List of strings | The name(s) of the provider(s) that found information about this host                                                    |

## Sampling

With globs and attribute filters, you can filter for hosts matching certain criteria. Once you've
got your filters down, you can sample your data if you want. For example, you may want to run a
command on one host per site do do connectivity checks to the outside:

```console
$ herd run -o inline site=~iad site:1 -- ping -q -c1 -n 1.1.1.1 \| grep rtt
5 done, 5 ok, 0 fail, 0 error in 1s
host-12.site-4.example.com  rtt min/avg/max/mdev = 0.860/0.860/0.860/0.000 ms
host-34.site-2.example.com  rtt min/avg/max/mdev = 1.910/1.910/1.910/0.000 ms
host-56.site-3.example.com  rtt min/avg/max/mdev = 2.037/2.037/2.037/0.000 ms
host-87.site-1.example.com  rtt min/avg/max/mdev = 0.633/0.633/0.633/0.000 ms
host-90.site-5.example.com  rtt min/avg/max/mdev = 1.037/1.037/1.037/0.000 ms
```

You can sample on multiple attributes by separating them with `:`. If your attribute names have a
`:` in them, you need to double them when sampling. For example, to get 2 hosts for each (az, os)
tuple, you can use `availability:zone:os::distro::codename:2` as sampling parameter.

## Full syntax for filters

Combining all these features, the full syntax for queries is:

- Sets of queries, separated by `+` or `-` charaters which perform set arithmetic on the individual
  queries.
- Individual queries have three parts, in this order. All parts are optional
  - A glob to match hostnames. Omitting this implies all hosts
  - Attribute filters. Omitting this leaves the result of the glob as-is
  - Sampling. Omitting this leaves the result of the filters as-is

## Sorting

As you can see in the previous examples, Herd by default sorts it output by hostname in ascending
order. This also applies to the output of `herd run`. You can however change this with the `--sort`
parameter. It accepts a comma-separated list of strings that are the fields to use for sorting. The
`name` field is always implies as last in that list. If your query selected hosts from a file, no
sorting is performed, and the hosts are shown in the same order as in that file.

```console
$ herd list --sort site,app
db-1.site-a.example.com
db-2.site-a.example.com
web-1.site-a.example.com
web-2.site-a.example.com
db-1.site-b.example.com
db-2.site-b.example.com
web-1.site-b.example.com
web-2.site-b.example.com
```

A special attribute to sort on is the `random` attribute. It's not really random, but a checksum of
the hosts name, so it will always be the same for a host but can be used to shuffle hosts into a
somewhat random yet repeatable order.

# More host information

So far all we've done is selecting and sorting hosts. That's all good if you just want to know the
names of hosts, or run commands on them. But herd has all thin information about hosts that you also
may be interested in, and it is more than happy to show you.

```console
$ herd list app=web --attributes os:distro:codename,instance_type
name                os:distro:codename   instance_type
web-1.example.com   stretch              m4.large
web-2.example.com   buster               m5.large
web-3.example.com   buster               m5.large
```

If yoo want to see all attributes, use `--all-attributes`. And if you want machine readable output,
you can use `--csv` and `--separator`.

For fully customized output, you can provide a go [text/template](https://pkg.go.dev/text/template)
template. For example: `herd list app=web --template 'Server {{.Name}} uses {{index .Attributes
"os:distro:codename"}}'`. This can also be used to show data in yaml format: `--template
'{{.|yaml}}'`.

## Counting hosts

Sometimes you're more interested in how many hosts match certain attributes. Herd has you covered
there too! You can easily show counts by one or more attributes:

```console
$ herd list app=web --count site,distribution
site     distribution   count
site-2   buster         5
site-1   buster         4
site-1   stretch        2
site-2   stretch        1
```

As you can see, the default ordering here is by count, but this can still be overwritten with the
`--sort` parameter. You can also group hosts by an attribute to get some more statistical
information. For example, to check whether hosts are balanced properly across availability zones

```console
$ herd list --count app --group availability_zone
app    us-east-1a   us-east-1b   us-east-1c   total   average   stddev
web    18           21           18           57      19        1.41
db     6            6            6            18      6         0
mail   1            0            0            2       0.33      0.47
api    1            1            3            5       1.67      0.94
```
