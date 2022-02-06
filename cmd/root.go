package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/pterm/pterm"
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

//The function gets executed once the application starts without any commands/arguments.
func startNoteSurvey() int {

	noteContent, err := survey.AskNoteContent()
	if err != nil {
		panic(err)
	}

	notebooks := cnote.GetNotebooks()
	fmt.Println("Getting your notebooks... This might take a while...")

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		log.Fatalf("An error has occured while starting notebook survey: %s", err.Error())
		return 1
	}

	sections := cnote.GetSections(n)
	section, err := survey.AskSection(n, sections)
	if err != nil {
		log.Fatalf("An error has occured while starting section survey: %s", err.Error())
		return 1
	}

	//Saving the note to the section
	err = cnote.SaveNotePage(onenote.NotePage{
		SectionId: section.ID,
		Content:   noteContent,
	})
	if err != nil {
		log.Fatalf("An error has occured while trying to save your note. %s", err.Error())
		return 1
	}

	//The note has been saved
	fmt.Println(pterm.Green(fmt.Sprintf("✅ Your note has saved to the section %s (%s)", section.Name, time.Now())))

	a, err := survey.AskAlias(onenote.NotebookName(n.DisplayName), onenote.SectionName(section.Name))
	if err != nil {
		log.Fatalf("An error has occured while trying to start alias survey. %s", err.Error())
		return 1
	}

	if a != "" {
		err := cnote.SaveAlias(a, n.DisplayName, section.Name)
		if err == nil {
			fmt.Println(pterm.Green(fmt.Sprintf("✅ Alias '%s' has been saved. (%s)", a, section.Name)))
		} else {
			return 2
		}
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
