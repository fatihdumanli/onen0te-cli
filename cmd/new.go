package main

import (
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/spf13/cobra"
)

var (
	alias string
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"add", "save"},
	Short:   "Create a new note",
	Long:    "Create a note on one of your Onenote sections",
	Run: func(c *cobra.Command, args []string) {

		if len(args) != 1 {
			c.Usage()
			return
		}

		noteContent := args[0]
		_ = noteContent

		a, ok := storage.GetAlias(alias)

		var appOptions = cnote.GetOptions()
		if !ok {
			fmt.Fprintf(appOptions.Out, "The alias %s couldn't be found.\n", alias)
			os.Exit(1)
		}

		fmt.Println(a)

	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	newCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias for the target onenote section")
	rootCmd.AddCommand(newCmd)
}
