package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var versionString string
var rootCmd = &cobra.Command{
	Use:   "lane",
	Short: "Automates common tasks",
	Long:  `lane is a task automation helper that works well with tools like make.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = versionString
	rootCmd.AddCommand(newDocumentationCommand(rootCmd))
	rootCmd.AddCommand(newTranslationsCommand())
	rootCmd.AddCommand(newWorkflowCommand())
}
