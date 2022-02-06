package onenote

type NotebookName string
type SectionName string
type AliasName string

//Represents a onenote notebook
type Notebook struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	SectionsUrl string `json:"sectionsUrl"`
}

//Represents a section in a notebook
type Section struct {
	Name string `json:"displayName"`
	ID   string `json:"id"`
}

//Represents a section alias
type Alias struct {
	Notebook NotebookName `json:"n"`
	Section  SectionName  `json:"s"`
}

//Represents a note page
type NotePage struct {
	SectionId string
	Content   string
}
