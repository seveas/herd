---
Title: Contributing to herd
Weight: 8
---
# GitHub

Herd is developed on GitHub. Here you can find

- [Source code](https://github.com/seveas/herd/)
- [Issue tracker](https://github.com/seveas/herd/issues/)
- [Pull requests](https://github.com/seveas/herd/pulls/)

If you have questions about Herd, or problems using it, please use the issue tracker. If you do not
have a GitHub account and would prefer not to create one, you can contact the author at
[dennis@kaarsemaker.net](mailto:dennis@kaarsemaker.net).

If you'd like to make a code or documentation contribution, a pull request would be the best way.
Before implementing big new features, or new providers, it would be nice to first open an issue to
discuss them to make sure they fit in the design and goals of Herd.

# Development guidelines

The main guideline is to run `make test` and `make test-integration` before filing a pull request
and to clean up your commits so each individual commit makes sense on its own and passes tests.
Rebasing and fixing up PR branches is perfectly fine, and we prefer clean branches over not
rewriting history.

Part of `make test` is a linter to make sure that new code follows Go standards about code
formatting and code use. These tests must pass before we can merge any contribution.

# FAQ

Some questions pop up more often than others. Here are the most common ones

#### Why not use ansible? Or puppet? Or ....?

Herd tries to fill a different niche. It's meant for troubleshooting, quickly gathering info and
running one-off commands in environments of all sizes, from dozens to thousands of servers. It is
not meant for configuration management, orchestration or deployments. You could build such things on
top, but there are many excellent products already for these tasks.

#### Why is my ssh configuration not respected?

Because Herd is not using your ssh client. Herd will read your OpenSSH configuration and [use some
settings](/configuration_reference/#openssh-configuration), but not all parameters are supported.

#### Can't Herd read my ssh key from a file?

No. Herd does not want access to your private key and enforces the use of an ssh agent.

#### Why is Herd having trouble authenticating?

Most likely because your SSH agent is not running, not forwarded or has not keys loaded. See the
[SSH agent documentation](/documentation/getting_started/#ssh-agent-setup)

# What do the release names mean?

They are all characters from the Discworld books by [Terry
Pratchett](https://www.terrypratchettbooks.com/). More specifically, they are characters from the
Tiffany Aching series. A highly recommended series, and a very inspirational witch, cheese maker and
granddaughter to the best shepherd to ever roam the chalk.
