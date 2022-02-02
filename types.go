package main


type AuthenticationResult int

const (
	Successful AuthenticationResult = iota
	Failed	AuthenticationResult
)


type Notebook struct {
	Name string
}

type Section struct {
	Name string
}

type NotebookName string

