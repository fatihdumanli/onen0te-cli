package main

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{}

//TODO:
//-list aliases
//-new alias
func init() {
	rootCmd.AddCommand(listCmd)
}
