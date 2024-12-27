# Antithesis

[![Build](https://github.com/guergabo/antithesis-cli/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/guergabo/antithesis-cli/actions/workflows/build.yml)

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
[Releases](https://github.com/guergabo/antithesis/releases) page.

### Create an Account

To create an **Antithesis** account, or login to an existing account:

```console
antithesis auth login # signup for waitlist ? 
```

### Create a Test Run

To create your first **Antithesis** test run, see our
[Getting Started](https://www.antithesis.com/docs/getting_started/) guide.

```console 
antithesis run \
  --name='quickstart' \
  --tenant='tenat-name' \
  --description='Running a quick antithesis test.' \
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
