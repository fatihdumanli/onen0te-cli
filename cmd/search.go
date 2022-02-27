package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/fatihdumanli/onenote"
	"github.com/fatihdumanli/onenote/internal/style"
	"github.com/fatihdumanli/onenote/internal/survey"
	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <phrase-to-search>",
	Short: "do a search in your notes",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the phrase argument to do a search")
		}
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = search(args[0])
		os.Exit(code)
		return err
	},
}

func search(query string) (_ int, err error) {

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "in search() (search.go)\n")
		}
	}()

	pages, err := onenote.Search(url.QueryEscape(query))
	if err != nil {
		return 1, err
	}

	if len(pages) == 0 {
		fmt.Println(style.Warning("couldn't find any match"))
		return 0, nil
	}

	options := []string{"◀️ Back to the search results", "❌ Exit"}
	var answer = 0
	var notepage msftgraph.NotePage
	var content []byte
	var section msftgraph.Section

	for answer != 1 {
		switch answer {
		case 0:
			notepage, err = survey.AskPage(pages)
			//fmt.Println(notepage)
			if err != nil {
				return 1, err
			}
			//page content - section data, concurrent
			content, err = onenote.GetPageContent(notepage)
			if err != nil {
				return 1, err
			}

			section, err = onenote.GetSection(notepage.ParentSection.ID)
			notepage.Section = section
			if err != nil {
				return 1, err
			}
		default:
		}

		answer, _ = displayContent(options, notepage, &content)
	}

	return 0, nil
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
