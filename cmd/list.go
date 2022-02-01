package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(listCmd)
}
