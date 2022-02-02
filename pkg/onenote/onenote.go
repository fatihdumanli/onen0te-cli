package onenote

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fatihdumanli/cnote/pkg/oauthv2"
)

type Notebook struct {
	Name string
}
type Section struct {
	Name string
}
type NotebookName string

func Authorize(opts oauthv2.OAuthParams, output io.Writer) {

	oauthv2.Authorize(opts, output)
	//TODO: store token

}

func GetNotebooks() ([]Notebook, error) {

	//Load token from local storage.
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/v1.0/me/onenote/notebooks", nil)

	tmpToken := "eyJ0eXAiOiJKV1QiLCJub25jZSI6Ii1BeklQQ1BxaU96WmhWbzhEZlN3dGVhX3RkUEJESC0zVzZuUmpUc2dBeWsiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik1yNS1BVWliZkJpaTdOZDFqQmViYXhib1hXMCIsImtpZCI6Ik1yNS1BVWliZkJpaTdOZDFqQmViYXhib1hXMCJ9.eyJhdWQiOiIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC8zMTk4NmVlOS04ZDBkLTRhOGUtOGM2ZC0xZDc2M2I2NmQ2YzIvIiwiaWF0IjoxNjQzNzU1NDcwLCJuYmYiOjE2NDM3NTU0NzAsImV4cCI6MTY0Mzc2MDIzOSwiYWNjdCI6MCwiYWNyIjoiMSIsImFpbyI6IkFVUUF1LzhUQUFBQURUcUtxcGs2WXZsMXkvZ0ZWV0U3Tlk5TnArcVJONVc1UEozU1lhWmMyNW5BOXY4YzJydG0xMGhYRWlmZGxwTzZDM2M5bmhvVXJXL1VhNjAvSkVZOEVnPT0iLCJhbHRzZWNpZCI6IjE6bGl2ZS5jb206MDAwM0JGRkQ4RDAyMUM3MSIsImFtciI6WyJwd2QiXSwiYXBwX2Rpc3BsYXluYW1lIjoiY25vdGUgLSB5b3VyIG5vdGVib29rIG9uIHRlcm1pbmFsIiwiYXBwaWQiOiIzMmRkYzJhYS1iMWQ0LTRhN2YtYTgxZi1mYjNkOWNhZjI0YWMiLCJhcHBpZGFjciI6IjAiLCJlbWFpbCI6ImZhdGloZHVtYW5saUBsaXZlLmNvbSIsImZhbWlseV9uYW1lIjoiRHVtYW5sxLEiLCJnaXZlbl9uYW1lIjoiRmF0aWgiLCJpZHAiOiJsaXZlLmNvbSIsImlkdHlwIjoidXNlciIsImlwYWRkciI6Ijc4LjE3OC43Ni4yOSIsIm5hbWUiOiJGYXRpaCBEdW1hbmzEsSIsIm9pZCI6IjZlM2JiNDY1LTc0MDctNGM1ZS05NDY2LWVlMTYyNWQ1MTlkOSIsInBsYXRmIjoiMyIsInB1aWQiOiIxMDAzN0ZGRTg4Qjg5NERGIiwicmgiOiIwLkFRc0E2VzZZTVEyTmprcU1iUjEyTzJiV3dnTUFBQUFBQUFBQXdBQUFBQUFBQUFBTEFMby4iLCJzY3AiOiJNYWlsLlJlYWQgTm90ZXMuQ3JlYXRlIE5vdGVzLlJlYWQgTm90ZXMuUmVhZFdyaXRlIE5vdGVzLlJlYWRXcml0ZS5BbGwgVXNlci5SZWFkIHByb2ZpbGUgb3BlbmlkIGVtYWlsIiwic2lnbmluX3N0YXRlIjpbImttc2kiXSwic3ViIjoiTmh4dkx4UVNBbTlzRXBhcnBOelk5dU1COXJIX2N3S2hMOVZCWmVudUVEZyIsInRlbmFudF9yZWdpb25fc2NvcGUiOiJFVSIsInRpZCI6IjMxOTg2ZWU5LThkMGQtNGE4ZS04YzZkLTFkNzYzYjY2ZDZjMiIsInVuaXF1ZV9uYW1lIjoibGl2ZS5jb20jZmF0aWhkdW1hbmxpQGxpdmUuY29tIiwidXRpIjoidDVtSXdMZ2pWay02eXJhQkZmOU5BQSIsInZlciI6IjEuMCIsIndpZHMiOlsiNjJlOTAzOTQtNjlmNS00MjM3LTkxOTAtMDEyMTc3MTQ1ZTEwIiwiYjc5ZmJmNGQtM2VmOS00Njg5LTgxNDMtNzZiMTk0ZTg1NTA5Il0sInhtc19zdCI6eyJzdWIiOiJuYW1FN1pUQmkyTndmRDhKbDgxdXVidVU2NWJSYzFOcUpyVnM1ZVk3Zk1BIn0sInhtc190Y2R0IjoxMzkxOTUxMzc3fQ.SZDsW_p3nFjoM-CYBEICadeKSrZv2lKSAhpiSG7YrgxqeliEuy5lTg6saLhwgojI013QSNsqOSeeMAMzYSJcuCpDqaalY7Xg6lcFVeWfhQANJJ0pu0duX3Rw1HMbr4fPXtyhv_LHEMDG8r_gBMdy_3g0AWOKbrvL4xPeW4LIfddwrOVj2WDHaWlwfetJmRKuw6KgEOfvGIQkLixftj-m1rNzYyLDSgB3ZQ0pFSWmf5XvtvE8VKLKkyb1wrR4RwDoYEB0Aftk3CClB8n3EnhM7p_mAzSVCYilqVq0o_MIc51HC2Y6Nnbtdx7e6Ip57jjzSYLJf40PQT29ATg4bYPDQg"

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tmpToken))
	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)

	bodyStr := string(body)
	fmt.Println(bodyStr)

	if err != nil {
		return []Notebook{}, err
	}

	_ = request
	_ = client
	_ = resp

	return getDummyNotebooks(), nil
}

func GetSections(n NotebookName) ([]Section, error) {
	return getDummySections(), nil
}

func getDummyNotebooks() []Notebook {
	return []Notebook{
		{"Fatih's Notebook"},
		{"Domain Driven Design"},
		{"Microservices"},
		{"Golang"},
	}
}

func getDummySections() []Section {
	return []Section{
		{"Quick notes"},
		{"Go"},
		{"Projects"},
		{"Todos"},
	}
}
