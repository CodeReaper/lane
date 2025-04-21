![version](https://img.shields.io/github/v/release/CodeReaper/lane)
![tests](https://github.com/CodeReaper/lane/actions/workflows/tests.yaml/badge.svg)
![license](https://img.shields.io/github/license/CodeReaper/lane.svg)

# lane

`lane` is a task automation helper.

## Installation

The preferred method of installation is through [asdf](http://asdf-vm.com/).

A lane plugin to install has been set up at [asdf-lane](https://github.com/CodeReaper/asdf-lane).

Alternatively this tool can be run directly:
```go
go run github.com/codereaper/lane@1
```

## Containerized

You can run lane using docker by running:
```sh
docker run -it --rm ghcr.io/codereaper/lane
```

## Completion

You can set up auto completion by adding the following to your dot rc file:

> ~/.zshrc
```
source <(lane completion zsh)
```

> ~/.bashrc
```
source <(lane completion bash)
```

## Documentation

[Auto-generated documentation](docs/generated/lane.md) is available, but is also included in `lane`.
