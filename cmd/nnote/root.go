package main

import (
	"log"
	"os"

	onenote "github.com/fatihdumanli/onen0te-cli"
	"github.com/fatihdumanli/onen0te-cli/internal/style"
	"github.com/spf13/cobra"
)

var (
	cachePath string
)

var rootCmd = &cobra.Command{
	Long:                  "Take notes on your Onenote notebooks from terminal",
	Use:                   "nnote",
	DisableFlagsInUseLine: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cachePath, "cache", "c", "$HOME/.config/nnote", "Cache path")
}

func Execute() {
	onenote.Init(onenote.Options{CachePath: os.ExpandEnv(cachePath)})
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("\n" + style.Error(err.Error()))
		os.Exit(1)
	}

}
