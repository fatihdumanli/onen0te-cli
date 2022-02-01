package cmd

import (
	"fmt"
	"os"

	"github.com/fatihdumanli/cnote/config"
	"github.com/fatihdumanli/cnote/survey"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long:                  "Take notes on your Onenote notebooks from terminal",
	Run:                   runRoot,
	Use:                   "cnote [command] [args] [flags]",
	DisableFlagsInUseLine: true,
}

func runRoot(c *cobra.Command, args []string) {

	var defaultOptions = config.GetOptions()
	_ = defaultOptions

	noteContent, err := survey.AskNoteContent(defaultOptions)

	if err != nil {
		panic(err)
	}

	notebook, err := survey.AskNotebook(defaultOptions)
	section, err := survey.AskSection(defaultOptions, notebook)

	_ = notebook
	_ = noteContent
	_ = section
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
