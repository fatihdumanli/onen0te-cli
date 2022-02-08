package main

import (
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <alias>",
	Aliases: []string{"delete"},
	Short:   "remove an alias",
	Args:    cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		os.Exit(removeAlias(c, args))
	},
	DisableFlagsInUseLine: true,
}

func removeAlias(c *cobra.Command, args []string) int {
	if len(args) != 1 {
		c.Usage()
		return 1
	}

	err := cnote.RemoveAlias(args[0])
	if err != nil {
		return 2
	}

	return 0
}

func init() {
	cmdAlias.AddCommand(removeCmd)
}
