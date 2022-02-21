package authentication

import (
	"testing"

	"github.com/fatihdumanli/onenote/internal/config"
	"github.com/fatihdumanli/onenote/internal/storage"
	"github.com/fatihdumanli/onenote/pkg/oauthv2"
)

type StorerStub struct {
	storage.Storer
}

type OAuthClientStub struct {
	*oauthv2.OAuthClient
}

func Test_AuthenticateUser(t *testing.T) {

	data := []struct {
		name            string
		storerStub      storage.Storer
		oauthClientStub *oauthv2.OAuthClient
		token           oauthv2.OAuthToken
		errMsg          string
	}{
		{"authenticate-happy-path", StorerStub{}, &oauthv2.OAuthClient{}, oauthv2.OAuthToken{}, ""},
	}

	//TODO: Complete
	for _, d := range data {
		t, err := AuthenticateUser(d.oauthClientStub, config.AppOptions{}, d.storerStub)

		_ = t
		_ = err

	}

}

func Test_RefreshToken(t *testing.T) {
}
