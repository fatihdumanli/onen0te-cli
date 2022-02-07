package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/internal/style"
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
			var msg = fmt.Sprintf("The alias %s has not found.\n", style.Alias(args[0]))
			fmt.Println(style.Error(msg))
			return 2
		}

		return 3
	}

	var msg = fmt.Sprintf("The alias %s has been deleted.\n", style.Alias(args[0]))
	fmt.Println(style.Success(msg))
	return 0
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
