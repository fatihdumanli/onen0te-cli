package main

import "github.com/spf13/cobra"

var cmdAlias = &cobra.Command{
	Use:   "alias <command>",
	Short: "add/list and remove alias",
	Long:  "aliases are used for quick access to onenote sections. you can quickly add a new note to any onenote section by specifiying an alias along with your command input.",
}

func init() {
	rootCmd.AddCommand(cmdAlias)
}
