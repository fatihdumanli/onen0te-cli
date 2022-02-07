package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/style"
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
	if err != nil {
		panic(err)
	}

	notebooks := cnote.GetNotebooks()
	fmt.Println("Getting your notebooks...")
	n, err := survey.AskNotebook(notebooks)

	fmt.Println("Getting sections...")
	sections := cnote.GetSections(n)
	section, err := survey.AskSection(n, sections)

	//Saving the note to the section
	link, err := cnote.SaveNotePage(onenote.NotePage{
		SectionId: section.ID,
		Content:   noteContent,
	})
	//TODO: We might turn this responsibility over to cnote package. (Printing the outcome)
	//Sİnce we're creating aliases, notes and doing pretty much the same thing everywhere
	//It doesn't add up to write the same code again and again.
	if err != nil {
		log.Fatal("couldn't save the note.")
		return 1
	}

	//The note has been saved
	var msg = fmt.Sprintf("Your note has saved to the section %s (%s)\n\n", style.Section(section.Name), time.Now())
	//TODO: We might turn this responsibility over to cnote package. (Printing the outcome)
	//Sİnce we're creating aliases, notes and doing pretty much the same thing everywhere
	//It doesn't add up to write the same code again and again.
	fmt.Println(style.Success(msg))

	fmt.Println(fmt.Sprintf("%s\n", link))
	a, err := survey.AskAlias(onenote.NotebookName(n.DisplayName), onenote.SectionName(section.Name))

	if a != "" {
		err := cnote.SaveAlias(a, n, section)
		if err == nil {
			//TODO: We might turn this responsibility over to cnote package. (Printing the outcome)
			//Sİnce we're creating aliases, notes and doing pretty much the same thing everywhere
			//It doesn't add up to write the same code again and again.
			var msg = fmt.Sprintf("Alias '%s' has been saved. (%s)\n", style.Alias(a), style.Section(section.Name))
			fmt.Println(style.Success(msg))
			var infoMsg = "Now you can quickly add notes with the following command:\n\n"
			fmt.Println(style.Info(infoMsg))
			fmt.Println(fmt.Sprintf("$ cnote new <path-to-input> -a %s", a))

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
