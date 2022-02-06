package main

import (
	"fmt"
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
	Run: func(c *cobra.Command, args []string) {
		os.Exit(startNoteSurvey())
	},
	Use: "cnote [command] [args] [flags]",
}

//TODO: return different integers
//The function gets executed once the application starts without any commands/arguments.
func startNoteSurvey() int {

	_, err := survey.AskNoteContent()
	if err != nil {
		panic(err)
	}

	notebooks := cnote.GetNotebooks()
	fmt.Fprintln(out, "Getting your notebooks... This might take a while...")

	n, err := survey.AskNotebook(notebooks)
	sections := cnote.GetSections(n)
	if err != nil {
		panic(err)
	}

	section, err := survey.AskSection(n, sections)
	if err != nil {
		log.Fatalf("An error has occured while starting section survey: %s", err.Error())
		return -1
	}

	//TODO: save the note.
	fmt.Fprintf(out, "Your note has saved to the notebook %s and the section %s",
		n.DisplayName, section.Name)

	a, err := survey.AskAlias(onenote.NotebookName(n.DisplayName), onenote.SectionName(section.Name))
	if err != nil {
		panic(err)
	}

	if a != "" {
		cnote.SaveAlias(a, n.DisplayName, section.Name)
	}

	return 1
}

func Execute() {
	//TODO: figure out how to deal with output
	out = os.Stdout

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

}
