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
	var ok bool
	notebooks, ok := cnote.GetNotebooks()
	if !ok {
		return 1
	}

	n, err := survey.AskNotebook(notebooks)
	sections, ok := cnote.GetSections(n)
	if !ok {
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
