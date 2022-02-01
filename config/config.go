package config

import (
	"io"
	"os"
)

type AppOptions struct {
	out io.Writer
	in  io.Reader
}

func GetOptions() AppOptions {
	return AppOptions{
		in:  os.Stdin,
		out: os.Stdout,
	}
}
