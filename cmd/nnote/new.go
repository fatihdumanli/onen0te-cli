package main

import (
	"fmt"

	errors "github.com/pkg/errors"

	onenote "github.com/fatihdumanli/onen0te-cli"
	"github.com/fatihdumanli/onen0te-cli/internal/style"
	"github.com/fatihdumanli/onen0te-cli/internal/survey"
	"github.com/fatihdumanli/onen0te-cli/internal/util/file"
	"github.com/fatihdumanli/onen0te-cli/pkg/msftgraph"
	"github.com/spf13/cobra"
)

var (
	alias         string
	title         string
	inlineText    string
	inputFilePath string
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"add", "save"},
	Short:   "Create a new note",
	Long:    "Create a note on one of your Onenote sections",
	RunE: func(c *cobra.Command, args []string) error {
		var _, err = saveNote(c, args)
		return err
	},
}

//Three methods to save a note
//1. Via default editor (Nano, vim or whatever.
//2. Inline text
//3. From a file
func saveNote(c *cobra.Command, args []string) (int, error) {

	var noteContent *string

	if inputFilePath != "" {
		//File specified
		if !file.Exists(inputFilePath) {
			return 3, fmt.Errorf("the file %s not found\n", inputFilePath)
		}
		fileContent, err := file.ReadString(inputFilePath)
		if err != nil {
			return 4, errors.Wrap(err, "in  saveNote()\n")
		}

		noteContent = &fileContent

	} else if inlineText != "" {
		//Inline text
		noteContent = &inlineText
	} else {
		//Launch the editor
		content, err := survey.AskNoteContent()
		if err != nil {
			return 1, errors.Wrap(err, "in saveNote()")
		}
		noteContent = &content
	}

	//Get confirmation on adding a note without a title.
	if title == "" {
		var tAnswer, err = survey.AskTitle()
		if err != nil {
			return 3, errors.Wrap(err, "in saveNote()\n")
		}
		title = tAnswer
	}

	var section msftgraph.Section

	if alias == "" {
		var notebooks, err = onenote.GetNotebooks()
		if err != nil {
			return 1, errors.Wrap(err, "getNotebooks operation has failed\n")
		}
		n, err := survey.AskNotebook(notebooks)
		if err != nil {
			return 1, errors.Wrap(err, "askNotebook operation has failed")
		}

		sections, err := onenote.GetSections(n)
		if err != nil {
			return 1, errors.Wrap(err, "getSections operation has failed")
		}
		section, err = survey.AskSection(n, sections)
		if err != nil {
			return 1, errors.Wrap(err, "askSection operation has failed")
		}

	} else {
		var a, err = onenote.GetAlias(alias)
		if err != nil {
			return 1, errors.Wrap(err, "alias data couldn't be loaded")
		}

		if a == nil {
			var errMsg = fmt.Sprintf("the alias %s does not exist", alias)
			fmt.Println(style.Error(errMsg))
			return 1, fmt.Errorf(errMsg)
		}

		section = a.Section
	}

	//Save the note. Show alias instructions only if the user could've used an alias for the section.
	_, err := onenote.SaveNotePage(*msftgraph.NewNotePage(section, title, *noteContent), alias == "")
	if err != nil {
		return 1, errors.Wrap(err, "saveNotePage operation has failed")
	}

	return 0, nil
}

func init() {
	newCmd.PersistentFlags().StringVarP(&inputFilePath, "file", "f", "", "use this flag to send a file to onenote")
	newCmd.PersistentFlags().StringVarP(&inlineText, "inline", "i", "", "use this flag to save an inline note")
	newCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias for the target onenote section.")
	newCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "title for the note page.")

	rootCmd.AddCommand(newCmd)
}
