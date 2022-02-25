package main

import (
	"os"
	"regexp"

	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/onenote"
	"github.com/fatihdumanli/onenote/internal/survey"
	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Notebook = msftgraph.Notebook
type Section = msftgraph.Section

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse the pages within a onenote section",
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = list()
		os.Exit(code)
		return err
	},
}

func list() (_ int, err error) {

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "in browse()\n")
		}
	}()

	var notebooks []Notebook
	notebooks, err = onenote.GetNotebooks()
	if err != nil {
		return 1, err
	}

	var answer string = "1"
	var content []byte
	var s Section

	for answer != "exit" {
		switch answer {
		case "1":
			n, err := survey.AskNotebook(notebooks)
			if err != nil {
				return 1, err
			}
			s, err = askSection(n)
			if err != nil {
				return 1, err
			}

			pages, err := onenote.GetPages(s)
			if err != nil {
				return 1, err
			}
			p, err := survey.AskPage(pages)
			if err != nil {
				return 1, err
			}
			content, err = onenote.GetPageContent(p)
		case "2":
			pages, err := onenote.GetPages(s)
			if err != nil {
				return 1, err
			}
			p, err := survey.AskPage(pages)
			if err != nil {
				return 1, err
			}
			content, err = onenote.GetPageContent(p)
		}
		answer, _ = displayContent(&content)
	}

	return 0, nil
}

func askSection(n msftgraph.Notebook) (msftgraph.Section, error) {
	sections, err := onenote.GetSections(n)
	if err != nil {
		return msftgraph.Section{}, err
	}

	s, err := survey.AskSection(n, sections)
	if err != nil {
		return msftgraph.Section{}, err
	}
	return s, nil
}

func displayContent(content *[]byte) (string, error) {

	var contentString = stripHtmlRegex(string(*content))

	//TODO: add display html in web browser
	options := []string{"◀️ Back to the section", "📃 Notebooks", "❌ Exit"}

	var navPrompt = &s.Select{
		Message: contentString,
		Options: options,
	}

	var answer string
	s.AskOne(navPrompt, &answer)

	if answer == "◀️ Back to the section" {
		return "2", nil
	} else if answer == "📃 Notebooks" {
		return "1", nil
	} else {
		return "exit", nil
	}

}

func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(`<.*?>`)
	return r.ReplaceAllString(s, "")
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
