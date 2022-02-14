package main

import (
	"log"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "Take notes on your Onenote notebooks from terminal",
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = startNoteSurvey()
		os.Exit(code)
		return err
	},
	Use:                   "cnote",
	DisableFlagsInUseLine: true,
}

//The function gets executed once the application starts without any commands/arguments.
func startNoteSurvey() (int, error) {
	noteContent, err := survey.AskNoteContent()
	if err != nil {
		return 1, errors.Wrap(err, "askNoteContent operation has failed")
	}

	notebooks, err := cnote.GetNotebooks()
	if err != nil {
		return 2, errors.Wrap(err, "getNotebooks operation has failed")
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1, errors.Wrap(err, "askNotebook operation has failed")
	}
	sections, err := cnote.GetSections(n)
	if err != nil {
		return 3, errors.Wrap(err, "getSection operation has failed")
	}
	section, err := survey.AskSection(n, sections)
	if err != nil {
		return 4, errors.Wrap(err, "askSection operation has failed")
	}

	title, err := survey.AskTitle()
	if err != nil {
		return 4, errors.Wrap(err, "askTitle operation has failed")
	}

	//Saving the note to the section
	_, err = cnote.SaveNotePage(*onenote.NewNotePage(section, title, noteContent), false)
	if err != nil {
		return 1, errors.Wrap(err, "saveNote operation has failed")
	}

	return 0, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
