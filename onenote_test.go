package onenote

import "testing"

type ApiStub struct{}

func TestGetNotebooks(t *testing.T) {

	notebooks, err := GetNotebooks()
	_ = notebooks
	_ = err

}
