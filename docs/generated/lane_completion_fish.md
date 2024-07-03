## lane completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	lane completion fish | source

To load completions for every new session, execute once:

	lane completion fish > ~/.config/fish/completions/lane.fish

You will need to start a new shell for this setup to take effect.


```
lane completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [lane completion](lane_completion.md)	 - Generate the autocompletion script for the specified shell

