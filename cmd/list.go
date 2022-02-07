package main

import (
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "display alias list",
	Run: func(c *cobra.Command, args []string) {
		os.Exit(displayAliasList())
	},
}

func displayAliasList() int {
	var aliasList = cnote.GetAliases()

	if aliasList == nil {
		fmt.Println(pterm.Red("Your alias data couldn't be loaded."))
		return 1
	}

	if len(*aliasList) == 0 {
		fmt.Println(pterm.Yellow("You haven't added any alias yet."))
		return 2
	}

	for _, a := range *aliasList {
		fmt.Printf("-%s=%s (%s)\n", pterm.Cyan(a.Short), a.Section.Name, a.Notebook.DisplayName)
	}

	return 0
}

func init() {
	rootCmd.AddCommand(listCmd)
}
