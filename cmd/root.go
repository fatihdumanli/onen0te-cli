package main

import (
	"io"
	"log"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/spf13/cobra"
)

var out io.Writer

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
		return 1, err
	}

	notebooks, err := cnote.GetNotebooks()
	if err != nil {
		return 2, err
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1, err
	}
	sections, err := cnote.GetSections(n)
	if err != nil {
		return 3, err
	}
	section, err := survey.AskSection(n, sections)

	title, err := survey.AskTitle()
	if err != nil {
		return 4, err
	}

	//Saving the note to the section
	_, err = cnote.SaveNotePage(*onenote.NewNotePage(section, title, noteContent), false)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func Execute() {
	out = os.Stdout
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
