package cnote

import (
	"github.com/fatihdumanli/cnote/internal/authentication"
	"github.com/fatihdumanli/cnote/internal/storage"
	"github.com/fatihdumanli/cnote/pkg/onenote"
)

type Cnote struct {
	storage storage.Storer
	auth    authentication.Authenticator
	api     onenote.Api
}

func (cnote *Cnote) SaveAlias(aname, nname, sname string) error {
	return nil
}

func (cnote *Cnote) GetAlias(n string) onenote.Alias {
	return onenote.Alias{}
}
