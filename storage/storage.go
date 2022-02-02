package storage

import (
	"errors"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

var InvalidTokenType = errors.New("Token type is invalid")

const TOKEN_KEY = "msgraphtoken"
const BUCKET = "cnote"

type TokenStatus int

const (
	DoesntExist TokenStatus = iota
	Expired
	Valid
)

//Friendly names
type Alias = onenote.Alias
type NotebookName = onenote.NotebookName
type SectionName = onenote.SectionName

type Storer interface {
	CheckToken() (oauthv2.OAuthToken, TokenStatus)
	StoreToken() error
	SaveAlias(a Alias, n NotebookName, s SectionName) error
}
