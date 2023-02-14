---
Title: Download herd
---
Herd is a single binary that you can download here for a variety of operating systems. All you need to do is put it in a directory on your $PATH or %PATH%

{{< download >}}

## Older versions

Older versions of herd can be found on the [GitHub releases page](https://github.com/seveas/herd/releases)

## Dependencies

Herd needs an SSH agent. It can use OpenSSH's agent as well as PuTTY. The minimum supported version
of OpenSSH is 8.4, released in September 2020. For macOS users, this means you will need at least
macOS 12 (Monterey). Debian users will need to use Debian 11 (Bullseye) or newer, though there is a
backport available in buster-backports too. The minimum supported version of Fedora is 33.

## Homebrew on mac

On MacOS, herd can be installed with the Homebrew package manager from a separate tap.

```console
$ brew tap seveas/herd
$ brew install herd
```

## Install from source

If you already have go installed, you can very easily install `herd`:

```console
$ go install github.com/seveas/herd/cmd/herd
```

## Dependencies for development

If you intend to work on the source code, you can also clone the [git repository](https://github.com/seveas/herd)
and run `make` to build the source. To build the code, you will need a working go install. To rebuild the
generated parser, you will need to install antlr. To regenerate the protobuf files, you will need to
install the protobuf tooling and the go code generator. On a mac this would look like

```console
$ brew install antlr protobuf go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ git clone https://github.com/seveas/herd.git
$ cd herd
$ make
```

And on a Debian system:

```console
$ apt install antlr protobuf-compiler golang-go
$ go install google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ git clone https://github.com/seveas/herd.git
$ cd herd
$ make
```
