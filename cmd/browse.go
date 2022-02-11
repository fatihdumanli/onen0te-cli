package main

import (
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse the pages within a onenote section",
	Run: func(c *cobra.Command, args []string) {
		os.Exit(browse())
	},
}

func browse() int {
	notebooks, err := cnote.GetNotebooks()
	if err != nil {
		return 1
	}

	n, err := survey.AskNotebook(notebooks)
	sections, err := cnote.GetSections(n)
	if err != nil {
		return 1
	}
	if err != nil {
		return 2
	}
	s, err := survey.AskSection(n, sections)
	_ = s

	return 0
}
func init() {
	rootCmd.AddCommand(browseCmd)
}
