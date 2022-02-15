package main

import (
	"os"

	"github.com/fatihdumanli/onenote"
	"github.com/fatihdumanli/onenote/internal/survey"
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

func browse() (_ int, err error) {

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "in browse()\n")
		}
	}()

	notebooks, err := onenote.GetNotebooks()
	if err != nil {
		return 1, err
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1, err
	}

	sections, err := onenote.GetSections(n)
	if err != nil {
		return 1, err
	}
	s, err := survey.AskSection(n, sections)
	_ = s
	if err != nil {
		return 1, err
	}

	return 0, nil
}
func init() {
	rootCmd.AddCommand(browseCmd)
}
