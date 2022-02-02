package onenote

type AuthenticationResult int

const (
	Successful AuthenticationResult = iota
	Failed
)

type Notebook struct {
	Name string
}
type Section struct {
	Name string
}
type NotebookName string
