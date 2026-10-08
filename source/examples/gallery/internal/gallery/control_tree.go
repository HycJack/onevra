package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/components"
)

func controlTree(c *ui.Context, s *controlState, variant string) {
	if variant == "TreeVirtualized" {
		controlTreeVirtualized(c, s, variant)
		return
	}
	components.Tree(c, func() {
		if components.TreeItem(c, "source", &s.on, s.choice == "source", func() {
			for _, name := range []string{"main.go", "go.mod"} {
				if components.TreeItem(c, name, nil, s.choice == name, nil).Clicked() {
					s.choice = name
				}
			}
		}).Clicked() {
			s.choice = "source"
		}
		if components.TreeItem(c, "README.md", nil, s.choice == "README.md", nil).Clicked() {
			s.choice = "README.md"
		}
	}).Label("Project")
}

var treeExampleFiles = func() []string {
	files := make([]string, 10000)
	for i := range files {
		files[i] = fmt.Sprintf("File %05d.go", i+1)
	}
	return files
}()

func controlTreeVirtualized(c *ui.Context, s *controlState, _ string) {
	s.tree.List.Selected = &s.index
	components.TreeVirtualized(c, &s.tree, []string{"source", "README.md"}, func(name string) []string {
		if name == "source" {
			return treeExampleFiles
		}
		return nil
	}, func(name string) {
		components.Text(c, name).SingleLine()
	}).Height(240).FillWidth().Label("Large project")
}
