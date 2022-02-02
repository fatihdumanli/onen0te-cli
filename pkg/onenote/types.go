package onenote

type NotebookName string
type SectionName string

type Notebook struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	SectionsUrl string `json:"sectionsUrl"`
	Sections    []Section
}

type Section struct {
	Name string `json:"displayName"`
}
