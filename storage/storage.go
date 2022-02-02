package storage

import (
	"errors"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

var InvalidTokenType = errors.New("Token type is invalid")

const TOKEN_KEY = "msgraphtoken"

type TokenStatus int

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
)

//TODO: we might need to store this type of field in AppOptions
type Storer interface {
	CheckToken() (oauthv2.OAuthToken, TokenStatus)
	StoreToken() error
	SaveAlias(a onenote.AliasName, n onenote.NotebookName, s onenote.SectionName) error
	GetAlias(a onenote.AliasName) (onenote.Alias, bool)
}
