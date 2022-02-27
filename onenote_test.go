package onenote

import (
	"testing"

	"github.com/fatihdumanli/onenote-cli/pkg/msftgraph"
	"github.com/google/go-cmp/cmp"
)

//Ergonomic alias
type Alias = msftgraph.Alias
type Notebook = msftgraph.Notebook
type Section = msftgraph.Section

func Test_GetAlias(t *testing.T) {
	data := []struct {
		name           string
		expectedAlias  Alias
		expectedErrMsg string
	}{
		{"a", Alias{Short: "a", Notebook: Notebook{DisplayName: "recursing-bohr"}, Section: Section{Name: "angry-payne"}}, ""},
		{"b", Alias{Short: "b", Notebook: Notebook{DisplayName: "beautiful-maxwell"}, Section: Section{Name: "beautiful-maxwell"}}, ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {

			//1. Create the alias. (Arrange)
			err := SaveAlias(d.name, d.expectedAlias.Notebook, d.expectedAlias.Section)
			if err != nil {
				t.Error(err)
			}

			//2. Get the alias (Act)
			alias, err := GetAlias(d.name)

			//3. Assert
			if diff := cmp.Diff(*alias, d.expectedAlias); diff != "" {
				t.Error("the result was different from the expected one")
				t.Error(diff)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.expectedErrMsg {
				t.Errorf("expected error msg %s, got %s", d.expectedErrMsg, errMsg)
			}

			//4. Remove the alias (Cleanup)
			err = RemoveAlias(d.name)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

//TODO: Complete
func Test_checkTokenPresented(t *testing.T) {
}
