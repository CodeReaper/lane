# lane

`lane` is a task automation helper.

## Installation

The prefered method of installation is through [asdf](http://asdf-vm.com/).

A lane plugin to install has been set up at [asdf-lane](https://github.com/CodeReaper/asdf-lane).

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

## Manuals

[lane](lane.d/help.md)

### Google APIs

- [lane google-api-docs-sheets-download](lane.d/google-api-docs-sheets-download/help.md)
- [lane google-api-jwt-generate](lane.d/google-api-jwt-generate/help.md)

### Mobile

- [lane mobile-static-resources-images](lane.d/mobile-static-resources-images/help.md)
- [lane mobile-update-translations](lane.d/mobile-update-translations/help.md)

### Shell

- [lane shell-github-action-semver-compare](lane.d/shell-github-action-semver-compare/help.md)
- [lane shell-run-github-workflow-tests](lane.d/shell-run-github-workflow-tests/help.md)
