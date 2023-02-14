---
Title: Running commands
Weight: 4
---

Once you've hand crafted your query to select the hosts you want to run a command on, it is time to
actually run a command. This starts simple, but Herd is quite flexible in how it runs commands and
how it shows what it is doing. Let's start with the basics:

```console
$ herd run app=web -- uptime
14 done, 14 ok, 0 fail, 0 error in 3s
{{<ansi green>}}web-1.example.com                                 completed successfully after 1s{{</ansi>}}
     02:13:44 up 85 days, 10:28,  1 user,  load average: 3.70, 3.75, 3.49
{{<ansi green>}}web-2.example.com                                 completed successfully after 1s{{</ansi>}}
     02:13:44 up 231 days, 19:29,  1 user,  load average: 3.89, 4.41, 4.17
{{<ansi green>}}web-3.example.com                                 completed successfully after 1s{{</ansi>}}
     02:13:46 up 94 days, 20:05,  1 user,  load average: 0.00, 0.04, 0.07
```

If the commands you want to run are short, and the amount of hosts to run them on is not that high,
this is all you need. But if either of those is not true, you'll quickly run into timeouts.

## Timeouts

Herd has very aggressive default timeouts, mostly to force you to think about reasonable timeouts
for your commands if they take more than a few seconds. There are three timeouts:

- `--connect-timeout` is how long TCP connections and SSH sessions may take to establish (3s by default)
- `--host-timeout` is how long a command may take on a host, including connection setup (10s by default)
- `--timeout` is a global timeout, 1 minute by default

These parameters take go-style arguments, so `1` means one nanosecond, `1s` means one second, `1m` one
minute and `1h` one hour.

## Thundering herds

If you try to run things on too many host at once, and they all access a single resource, you will
cause a problem known as a thundering herd (guess what inspired the name of this tool!). The first
thundering herd is the resources of the computer you run herd on. On a reasonable new macbook,
trying to do more than a few hundred simultaneous ssh connections slows things down tremendously. If
you encounter this, you can limit the parallelism of herd with the `--parallel` or `-p` parameter.

To play even nicer with shared resources, the `--splay` parameter introduces a random delay before
connecting to each host. This nicely separates the load. For example, to run puppet on a set of
hosts without overloading the puppet infrastructure, you can do something like:

```console
$ herd run role=web --parallel 50 --splay 30s --host-timeout 5m --timeout 30m -o tail -- sudo puppet agent -t
```

If you want to change parallelism while a command is running, you can send herd a `SIGURS1` signal
to increase parallelism by 50% or a `SIGUSR2` signal to decrease parallelism by 50%. When increasing
parallelism, new connections will be made immediately. When decreasing parallelism, running commands
will not be interrupted but new ones will only be started when enough tasks have finished.

## Output formatting

By default, herd shows a summary line, then per host a line indicating success/failure and then the
output of the command. There are a few more output modes, and the puppet example above uses one of
them: tail mode. In this mode, herd does not wait for output to arrive, but shows the output as it
comes in, prefixed with the name of the host that sent it.

```console
$ herd run *.example.com -o tail -- 'echo hello; sleep 3; echo world'
server-1.example.com  hello
server-2.example.com  hello
server-3.example.com  hello
{{<ansi green>}}server-1.example.com  completed successfully after 5s{{</ansi>}}
server-1.example.com  world
server-2.example.com  world
{{<ansi green>}}server-2.example.com  completed successfully after 5s{{</ansi>}}
server-3.example.com  world
{{<ansi green>}}server-3.example.com  completed successfully after 5s{{</ansi>}}
3 done, 3 ok, 0 fail, 0 error in 5s
```

Another useful output mode is the inline mode. Like the default mode, it waits for all output to
arrive, and like tail mode it prefixes the output with the hostname and does not show a summary
line. This is very useful if you want to compare the output of different hosts. Combining it with
the ability to sort the list of hosts by the content of the output makes it really easy to do things
like comparing package versions between hosts or finding the hosts that have been up the longest

```console
$ herd run app=web -o inline -s stdout -- uptime --since
14 done, 14 ok, 0 fail, 0 error in 3s
{{<ansi green >}}server-04.example.com{{</ansi>}}  2022-06-27 07:43:56
{{<ansi green >}}server-09.example.com{{</ansi>}}  2022-09-09 03:00:02
{{<ansi green >}}server-05.example.com{{</ansi>}}  2022-10-27 23:33:46
{{<ansi green >}}server-11.example.com{{</ansi>}}  2022-11-03 19:35:23
{{<ansi green >}}server-12.example.com{{</ansi>}}  2022-11-11 06:08:03
{{<ansi green >}}server-07.example.com{{</ansi>}}  2022-11-17 10:25:57
{{<ansi green >}}server-14.example.com{{</ansi>}}  2022-11-20 15:45:23
{{<ansi green >}}server-02.example.com{{</ansi>}}  2022-11-23 21:58:14
{{<ansi green >}}server-06.example.com{{</ansi>}}  2022-11-24 12:08:46
{{<ansi green >}}server-10.example.com{{</ansi>}}  2022-12-03 12:53:43
{{<ansi green >}}server-01.example.com{{</ansi>}}  2022-12-12 05:16:01
{{<ansi green >}}server-13.example.com{{</ansi>}}  2022-12-12 05:18:00
{{<ansi green >}}server-03.example.com{{</ansi>}}  2022-12-12 05:18:00
{{<ansi green >}}server-08.example.com{{</ansi>}}  2023-02-01 03:58:33
```

## History

The complete history of what you run with herd, including the output of all commands, is saved for
your convenience. This makes it possible to revisit the output of commands without having to re-run
them. The history is saved as a set of json files, and at the end of each `herd` invocation, it will
show you where the history of that invocation is stored.

## Interactive mode and scripting

Herd also has an interactive mode and a scripting interpreter. The syntax of the language used by
these scripts is fairly simple and resembles the command line parameter syntax. The main difference
is that typing is stricter and strings need to be quoted.

Here is an example script.

```sh
#!/usr/bin/herd run-script
#
# We find hosts where openssl to too old and upgrade it. Then we run puppet to make sure everything
# is happy on the host.

# We can run fast but not too fast
set Parallel 200

# We only want to run this on staging hosts and in the test site
add hosts environment == "staging" + site == "test-site"

# We check the version of openssl
run dpkg -l openssl | grep ^ii.*1.1.0l-1~deb9u6

# If grep succeeds, openssl has been updated, and we don't need to upgrade
remove hosts err == nil

# Now slowly upgrade
set Parallel 10
set Timeout 30m
set HostTimeout 5m
set Output "tail"
run sudo apt-get install openssl=1.1.0l-1~deb9u6 && sudo puppet agent -t
```

The rules of the syntax are as follows:

- A script consists of one or more lines, separated by newlines
- Each line is interpreted separately, both in scripted and in interactive mode. There is no way to
  split a command over multiple lines
- Lines starting with `#` are comments. The `#` character has no special meaning in other places on
  a line
- Each line may contain only one command
- String values must be quoted in double quotes like `"this"`
- Regular expressions in host filters must be enclosed in forward slashes, like `/this/`
- Duration values are written as numbers followed with `s`, `m`, or `h` and are not quoted. For
  example: `20s`

Each script operates on a set of hosts, and command either manipulate the set of hosts, run commands
on those hosts or set operating parameters.

| Command        | Parameters                                                                                                       |
|----------------|------------------------------------------------------------------------------------------------------------------|
| `set`          | Parameter and value to set                                                                                       |
| `add hosts`    | Sets of (glob, filters, sampling) pairs, similar to how you search on the command line                           |
| `remove hosts` | Sets of (glob, filters, sampling) pairs, similar to how you search on the command line                           |
| `list hosts`   | None. This command does not yet support `--attributes`, `--count` or `--group`                                   |
| `run`          | Command, unquoted. The rest of the line is passed verbatim to `sh -c` on the remove end, so no quoting is needed |

The parameters you can set correspond to the command line flags of the same name

| Parameter        | Type     | Example  |
|------------------|----------|----------|
| `Output`         | String   | `"tail"` |
| `Parallel`       | Integer  | `100`    |
| `ConnectTimeout` | Duration | `3s`     |
| `HostTimeout`    | Duration | `10s`    |
| `Timeout`        | Duration | `1m`     |
| `NoPager`        | Boolean  | `false`  |
