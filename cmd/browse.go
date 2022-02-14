package main

import (
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse the pages within a onenote section",
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = browse()
		os.Exit(code)
		return err
	},
}

func browse() (int, error) {
	notebooks, err := cnote.GetNotebooks()
	if err != nil {
		return 1, errors.Wrap(err, "getNotebooks operation has failed")
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1, errors.Wrap(err, "askNotebook operation has failed")
	}

	sections, err := cnote.GetSections(n)
	if err != nil {
		return 1, errors.Wrap(err, "getSections operation has failed")
	}
	s, err := survey.AskSection(n, sections)
	_ = s
	if err != nil {
		return 1, errors.Wrap(err, "askSection operation has failed")
	}

	return 0, nil
}
func init() {
	rootCmd.AddCommand(browseCmd)
}
