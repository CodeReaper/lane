package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// FIXME: lane shell-run-github-workflow-tests -i .github/workflows/tests.yaml

type RunConfig struct {
	input string
	kind  string
	job   string
}

func newWorkflowCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows locally",
	}
	cmd.AddCommand(newWorkflowRunCommand())
	return cmd
}

func newWorkflowRunCommand() *cobra.Command {
	var additionalHelp = `Reads a yaml file with the structure of a GitHub workflow and runs the 'run' steps.
Each step will reports if its exit code indicated success or failure and counts towards a tally.
The steps in each job is run on a fresh copy of the workspace.

The purpose is to enable unit-test-style tests for shell-based tooling.

If the contents of 'test.yaml' is:

    name: Tests

    on:
    workflow_dispatch: {}

    jobs:
    test-run:
        name: Test run
        runs-on: ubuntu-latest
        steps:
        - uses: actions/checkout@v3 # steps without the 'run' is ignored

        - name: Test that passes
          run: echo 'Test-in-tests'

        - name: Test that fails
          run: exit 1

The output would be:

Preparing runnner... done!
Preparing workspace...  done!
Executing runner...

Test run (test-run)
- Test that passes: Pass
- Test that fails: Failed with exit code 1

Tests; Total: 2 Passes: 1 Fails: 1
`
	var config = &RunConfig{
		kind: "github",
	}
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run workflow locally",
		Long:  additionalHelp,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.errors()
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("run it")
		},
	}
	cmd.Flags().StringVarP(&config.input, "input", "i", "", "Path to workflow file (Required)")
	cmd.Flags().StringVarP(&config.kind, "type", "t", "", "Type of workflow (defaults to GitHub)")
	cmd.Flags().StringVarP(&config.job, "job", "j", "", "A job ID in the workflow file")
	cmd.MarkFlagRequired("input")
	return cmd
}

func (c *RunConfig) errors() error {
	var validTypes = []string{"github"}

	for _, v := range validTypes {
		if v != strings.ToLower(c.kind) {
			return fmt.Errorf("invalid type: %s. Valid options are %v", c.kind, validTypes)
		}
	}

	return nil
}
