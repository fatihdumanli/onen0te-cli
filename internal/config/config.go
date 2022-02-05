package config

import (
	"io"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type AppOptions struct {
	Out         io.Writer
	In          io.Reader
	OAuthParams oauthv2.OAuthParams
}

type MicrosoftGraphConfig struct {
	ClientId    string
	TenantId    string
	RedirectUrl string
}
