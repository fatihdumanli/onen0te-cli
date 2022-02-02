package onenote

type GetNotebooksResponse struct {
	Notebooks []Notebook `json:"value"`
}

type Notebook struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type Section struct {
	Name string
}
type NotebookName string
