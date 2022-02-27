package main

import (
	"log"
	"os"

	"github.com/fatihdumanli/onenote-cli/internal/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long:                  "Take notes on your Onenote notebooks from terminal",
	Use:                   "nnote",
	DisableFlagsInUseLine: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("\n" + style.Error(err.Error()))
		os.Exit(1)
	}

}
