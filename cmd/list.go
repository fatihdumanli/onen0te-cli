package main

import (
	"os"
	"regexp"

	s "github.com/AlecAivazis/survey/v2"
	"github.com/fatihdumanli/onenote"
	"github.com/fatihdumanli/onenote/internal/style"
	"github.com/fatihdumanli/onenote/internal/survey"
	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	errors "github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type Notebook = msftgraph.Notebook
type Section = msftgraph.Section
type NotePage = msftgraph.NotePage

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

	var answer int = 1
	var content []byte
	var s Section
	var p NotePage

	//0: Back to the section
	//1: Notebooks
	//2: Exit
	//TODO: add display html in web browser
	options := []string{"◀️ Back to the section", "📃 Notebooks", "❌ Exit"}

	for answer != 2 {
		switch answer {
		case 1:
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
			p, err = survey.AskPage(pages)
			if err != nil {
				return 1, err
			}
			content, err = onenote.GetPageContent(p)
		case 0:
			pages, err := onenote.GetPages(s)
			if err != nil {
				return 1, err
			}
			p, err = survey.AskPage(pages)
			if err != nil {
				return 1, err
			}
			content, err = onenote.GetPageContent(p)
		}
		answer, _ = displayContent(&options, p, &content)
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

//Returns the option index
func displayContent(options *[]string, page msftgraph.NotePage, content *[]byte) (int, error) {

	//Title, section and date saved
	var contentString = stripHtmlRegex(string(*content))
	var tableData [][]string
	tableData = append(tableData, []string{style.TableHeader("Title"), style.TableHeader("Section"), style.TableHeader("Last Modified Date/time")})
	tableData = append(tableData, []string{page.Title, page.ParentSection.Name, page.LastModifiedDateTime.Format("2006-01-02 15:04:05 Monday")})

	tableStr, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
	tableStr = pterm.DefaultBox.Sprint(tableStr)
	contentString = "\n\n" + tableStr + "\n" + page.Links.OneNoteWebURL.Href + "\n" + contentString

	var navPrompt = &s.Select{
		Message: contentString,
		Options: *options,
	}

	var answer string
	err := s.AskOne(navPrompt, &answer)
	if err != nil {
		return -1, errors.Wrap(err, "couldn't start the navigation prompt survey")
	}

	var findAnswerIndex = func() int {
		for i := 0; i < len(*options); i++ {
			if (*options)[i] == answer {
				return i
			}
		}
		return -1
	}

	return findAnswerIndex(), nil
}

func stripHtmlRegex(s string) string {
	r := regexp.MustCompile(`<.*?>`)
	return r.ReplaceAllString(s, "")
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
