package config

import (
	"io"
)

type AppOptions struct {
	Out io.Writer
	In  io.Reader
}

type MicrosoftGraphConfig struct {
	ClientId    string
	TenantId    string
	RedirectUrl string
}
