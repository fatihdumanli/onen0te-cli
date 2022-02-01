package onenote

type Notebook struct {
	Name string
}
type Section struct {
	Name string
}
type NotebookName string

func GetNotebooks() ([]Notebook, error) {
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
