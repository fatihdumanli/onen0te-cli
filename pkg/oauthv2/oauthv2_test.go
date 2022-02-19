package oauthv2

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestOauth_GetToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/token", httpmock.NewStringResponder(200,
		`{"token_type":"Bearer","scope":"Notes.ReadWrite.All Notes.Create Notes.Read Notes.ReadWrite","expires_in":3600,"ext_expires_in":3600,"access_token":"dummy-access-token","refresh_token":"refresh-token"}`))
	httpmock.GetTotalCallCount()

	var token, err = getToken(getTokenParams{})
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}
