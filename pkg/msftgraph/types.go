package msftgraph

type NotebookName string
type SectionName string
type AliasName string

const (
	AliasesKey = "aliases"
)

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
	//Need the pointer in runtime. It's not gonna be saved.
	Notebook *Notebook
}

//Represents a section alias
type Alias struct {
	Short    string   `json:"a"`
	Notebook Notebook `json:"n"`
	Section  Section  `json:"s"`
}

//Represents a note page
type NotePage struct {
	Section Section
	Title   string
	Content string
}

func NewNotePage(s Section, t string, c string) *NotePage {
	return &NotePage{
		Section: s,
		Title:   t,
		Content: c,
	}
}
