package main

import (
	"log"
	"os"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal"
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

	var noteContent *string

	if len(args) != 1 && !inline {
		c.Usage()
		return 1
	}

	//Input validations
	if inline && inlineContent == "" {
		log.Fatal("content flag cannot be empty")
		return 2
	}

	if !inline {
		//Load the content from the file
		var inputPath = args[0]
		if !internal.Exists(inputPath) {
			log.Fatalf(" ❌ the file %s not found.", inputPath)
			return 3
		}

		var fileContent, ok = internal.ReadFile(inputPath)
		if !ok {
			log.Fatalf("😣 couldn't read the file %s", inputPath)
			return 4
		}

		noteContent = &fileContent
	} else {
		noteContent = &inlineContent
	}

	//Get confirmation on adding a note without a title.
	if title == "" {
		var tAnswer, _ = survey.AskTitle()
		title = tAnswer
	}

	var section onenote.Section

	if alias == "" {
		var notebooks, _ = cnote.GetNotebooks()
		n, _ := survey.AskNotebook(notebooks)

		var sections, _ = cnote.GetSections(n)
		section, _ = survey.AskSection(n, sections)
	} else {
		var a = cnote.GetAlias(alias)
		section = a.Section
	}

	_, err := cnote.SaveNotePage(*onenote.NewNotePage(section, title, *noteContent))
	if err != nil {
		return 5
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
