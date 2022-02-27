package msftgraph

import "time"

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
	Name     string    `json:"displayName"`
	ID       string    `json:"id"`
	Notebook *Notebook `json:"parentNotebook"`
}

//Represents a section alias
type Alias struct {
	Short    string   `json:"a"`
	Notebook Notebook `json:"n"`
	Section  Section  `json:"s"`
}

//Represents a note page
type NotePage struct {
	Section    Section
	Title      string `json:"title"`
	Content    string
	ContentUrl string `json:"contentUrl"`
	Links      struct {
		OneNoteClientURL struct {
			Href string `json:"href"`
		} `json:"oneNoteClientUrl"`
		OneNoteWebURL struct {
			Href string `json:"href"`
		} `json:"oneNoteWebUrl"`
	} `json:"links"`
	ParentSection        Section   `json:"parentSection"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}

func NewNotePage(s Section, t string, c string) *NotePage {
	return &NotePage{
		Section: s,
		Title:   t,
		Content: c,
	}
}
