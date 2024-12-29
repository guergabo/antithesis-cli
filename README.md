<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/user-attachments/assets/bc44d6d0-2fcf-4aa6-88da-6a444e3cbcc2">
    <img alt="antithesis logo" src="https://github.com/user-attachments/assets/4db6f803-6dcc-4118-8eb8-8ef188194a8f" height="100">
  </picture>
</p>

<hr />

[![Build](https://github.com/guergabo/antithesis-cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/guergabo/antithesis-cli/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/guergabo/antithesis-cli)](https://goreportcard.com/report/github.com/guergabo/antithesis-cli)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Antithesis** is a platform for building the impossible.

## Getting Started

### Install Antithesis

Install with Homebrew on macOS:

```console
brew tap guergabo/antithesis
brew install antithesis
```

Install with Go:

```console
go install github.com/guergabo/antithesis-cli@latest
```

Alternatively, you can download the latest `antithesis` binary from the
[Releases](https://github.com/guergabo/antithesis-cli/releases) page.

### Create a Project

To initialize an **Antithesis** project:

```console
antithesis init quickstart ./output
```

### Create a Test Run

To create your first **Antithesis** test run, see our
[Getting Started](https://www.antithesis.com/docs/getting_started/) guide.

```console 
antithesis run \
  --name='quickstart' \
  --description='Running a quick antithesis test.' \
  --tenant='tenant' \
  --username='username' \
  --password='password' \
  --config='<registry>/<namespace>/config:latest' \
  --image='<registry>/<namespace>/sut:latest' \
  --image='<registry>/<namespace>/test-template:latest' \
  --duration=15 \
  --email='xxx@gmail.com'
```

## Getting Help

See `antithesis --help` or our [documentation](https://www.antithesis.com/docs/) for
further information, or reach out on [Discord](https://discord.gg/XqRGqpHJ).
