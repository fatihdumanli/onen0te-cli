package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/style"
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

	sort.Slice(*aliasList, func(i, j int) bool {
		return (*aliasList)[i].Short < (*aliasList)[j].Short
	})

	if aliasList == nil {
		fmt.Println(style.Error("Your alias data couldn't be loaded."))
		return 1
	}

	if len(*aliasList) == 0 {
		fmt.Println(style.Error("You haven't added any alias yet."))
		return 2
	}

	var tableData [][]string
	tableData = append(tableData, []string{"Alias", "Section", "Notebook"})

	for _, a := range *aliasList {
		tableData = append(tableData, []string{a.Short, a.Section.Name, a.Notebook.DisplayName})
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	return 0
}

func init() {
	rootCmd.AddCommand(listCmd)
}
