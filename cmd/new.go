package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal/survey"
	"github.com/fatihdumanli/cnote/pkg/onenote"
	"github.com/spf13/cobra"
)

var (
	alias string
	//template string
	title         string
	inline        bool
	inlineContent string
)

var newCmd = &cobra.Command{
	Use:     "new <path to input>",
	Aliases: []string{"add", "save"},
	Short:   "Create a new note",
	Long:    "Create a note on one of your Onenote sections",
	Run: func(c *cobra.Command, args []string) {
		os.Exit(saveNote(c, args))

	},
}

func saveNote(c *cobra.Command, args []string) int {

	if len(args) != 1 && !inline {
		c.Usage()
		return 1
	}

	//Input validations
	if inline {
		if inlineContent == "" {
			log.Fatal("text flag cannot be empty")
			return 2
		}
	} else {
		//it cannot be empty
		var inputPath = args[0]
		_ = inputPath

	}

	var section onenote.Section

	if alias == "" {
		var notebooks = cnote.GetNotebooks()
		n, _ := survey.AskNotebook(notebooks)

		var sections = cnote.GetSections(n)
		section, _ = survey.AskSection(n, sections)
	} else {
		var a = cnote.GetAlias(alias)
		section = a.Section
	}

	if inline {

		err := cnote.SaveNotePage(onenote.NotePage{
			SectionId: section.ID,
			Title:     title,
			Content:   inlineContent,
		})

		if err != nil {
			log.Fatal("couldn't save the note.")
			return 3
		}

		fmt.Printf(" ✅ Your note has been saved. (%s)\n", section.Name)
	} else {

		//TODO: we should read from the file

	}

	return 0
}

func init() {
	newCmd.PersistentFlags().BoolVarP(&inline, "inline", "i", false, "specify this flag along with --text flag to save an inline note.")
	newCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias for the target onenote section.")
	newCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "title for the note page.")
	newCmd.PersistentFlags().StringVarP(&inlineContent, "content", "c", "", "inline content for the note")

	//newCmd.PersistentFlags().StringVarP(&template, "template", "t", "vanilla", "template for the note page that will be saved")
	rootCmd.AddCommand(newCmd)
}
