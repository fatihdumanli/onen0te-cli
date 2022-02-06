package main

import (
	"github.com/spf13/cobra"
)

var (
	alias    string
	template string
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"add", "save"},
	Short:   "Create a new note",
	Long:    "Create a note on one of your Onenote sections",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			c.Usage()
			return
		}

		noteContent := args[0]
		_ = noteContent

		//a, ok := storage.GetAlias(alias)

		//var appOptions = cnote.GetOptions()
		//if !ok {
		//	fmt.Fprintf(appOptions.Out, "The alias %s couldn't be found.\n", alias)
		//	os.Exit(1)
		//}

		//fmt.Println(a)

	},
	Args: cobra.MinimumNArgs(1),
}

func init() {
	newCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias for the target onenote section")
	newCmd.PersistentFlags().StringVarP(&template, "template", "t", "vanilla", "template for the note page that will be saved")
	rootCmd.AddCommand(newCmd)
}
