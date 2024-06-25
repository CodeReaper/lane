package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newDocumentationCommand(root *cobra.Command) *cobra.Command {
	var output string
	var cmd = &cobra.Command{
		Use:    "documentation",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(root, output)
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "Path to save documentation files (Required)")
	cmd.MarkFlagRequired("output")
	return cmd
}
