package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/pterm/pterm"
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
		if errors.Is(err, storage.KeyNotFound) {
			//TODO: Do not forget to use standardized color for aliases.
			fmt.Printf(" ❌ The alias %s has not found.\n", pterm.Blue(args[0]))
			return 2
		}

		return 3
	}

	fmt.Printf(" ✅The alias %s has been deleted\n", pterm.Blue(args[0]))
	return 0
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
