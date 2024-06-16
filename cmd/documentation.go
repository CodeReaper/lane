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
			err := doc.GenMarkdownTree(root, output)
			if err != nil {
				return err
			}
			err = doc.GenReSTTree(root, output)
			if err != nil {
				return err
			}
			err = doc.GenYamlTree(root, output)
			if err != nil {
				return err
			}
			err = doc.GenManTree(root, nil, output)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "o", "", "Path to save documentation files (Required)")
	cmd.MarkFlagRequired("output")
	return cmd
}
