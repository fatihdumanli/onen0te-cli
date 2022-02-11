package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/spf13/cobra"
)

var out io.Writer

var rootCmd = &cobra.Command{
	Long: "Take notes on your Onenote notebooks from terminal",
	Run: func(c *cobra.Command, args []string) {
		os.Exit(startNoteSurvey())
	},
	Use:                   "cnote",
	DisableFlagsInUseLine: true,
}

//The function gets executed once the application starts without any commands/arguments.
func startNoteSurvey() int {
	noteContent, err := survey.AskNoteContent()
	//TODO: consider creating constsns for error codes.
	if err != nil {
		return 1
	}

	notebooks, ok := cnote.GetNotebooks()
	if !ok {
		return 2
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1
	}
	sections, ok := cnote.GetSections(n)
	if !ok {
		return 3
	}
	section, err := survey.AskSection(n, sections)

	title, err := survey.AskTitle()
	if err != nil {
		return 4
	}

	//Saving the note to the section
	_, err = cnote.SaveNotePage(*onenote.NewNotePage(section, title, noteContent), false)
	if err != nil {
		return 1
	}

	return 0
}

func Execute() {
	out = os.Stdout
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
