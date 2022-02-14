package main

import (
	"fmt"
	"os"

	errors "github.com/pkg/errors"

	"github.com/fatihdumanli/cnote"
	"github.com/fatihdumanli/cnote/internal"
	"github.com/fatihdumanli/cnote/internal/style"
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
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = saveNote(c, args)
		os.Exit(code)
		return err
	},
}

func saveNote(c *cobra.Command, args []string) (int, error) {

	var noteContent *string

	if len(args) != 1 && !inline {
		c.Usage()
		return 1, fmt.Errorf("inline note requires exactly one argument")
	}

	//Input validations
	if inline && inlineContent == "" {
		return 2, fmt.Errorf("content flag cannot be empty")
	}

	if !inline {
		//Load the content from the file
		var inputPath = args[0]
		if !internal.Exists(inputPath) {
			return 3, fmt.Errorf("the file %s not found", inputPath)
		}

		var fileContent, ok = internal.ReadFile(inputPath)
		if !ok {
			return 4, fmt.Errorf("couldn't read the file %s", inputPath)
		}

		noteContent = &fileContent
	} else {
		noteContent = &inlineContent
	}

	//Get confirmation on adding a note without a title.
	if title == "" {
		var tAnswer, err = survey.AskTitle()
		if err != nil {
			return 3, errors.Wrap(err, "askTitle operation has failed")
		}
		title = tAnswer
	}

	var section onenote.Section

	if alias == "" {
		var notebooks, err = cnote.GetNotebooks()
		if err != nil {
			return 1, errors.Wrap(err, "getNotebooks operation has failed")
		}
		n, err := survey.AskNotebook(notebooks)
		if err != nil {
			return 1, errors.Wrap(err, "askNotebook operation has failed")
		}

		sections, err := cnote.GetSections(n)
		if err != nil {
			return 1, errors.Wrap(err, "getSections operation has failed")
		}
		section, err = survey.AskSection(n, sections)
		if err != nil {
			return 1, errors.Wrap(err, "askSection operation has failed")
		}

	} else {
		var a, err = cnote.GetAlias(alias)
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
	_, err := cnote.SaveNotePage(*onenote.NewNotePage(section, title, *noteContent), alias == "")
	if err != nil {
		return 1, errors.Wrap(err, "saveNotePage operation has failed")
	}

	return 0, nil
}

func init() {
	newCmd.PersistentFlags().BoolVarP(&inline, "inline", "i", false, "specify this flag along with --text flag to save an inline note.")
	newCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias for the target onenote section.")
	newCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "title for the note page.")
	newCmd.PersistentFlags().StringVarP(&inlineContent, "content", "c", "", "inline content for the note")

	//newCmd.PersistentFlags().StringVarP(&template, "template", "t", "vanilla", "template for the note page that will be saved")
	rootCmd.AddCommand(newCmd)
}
