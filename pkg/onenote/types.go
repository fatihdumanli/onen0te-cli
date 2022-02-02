package onenote

type NotebookName string
type SectionName string
type AliasName string

type Notebook struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	SectionsUrl string `json:"sectionsUrl"`
	Sections    []Section
}

type Section struct {
	Name string `json:"displayName"`
}

type Alias struct {
	Notebook NotebookName `json:"n"`
	Section  SectionName  `json:"s"`
}
