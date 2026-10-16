package gallery

import (
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

type tableExampleFile struct {
	name string
	size int
}

var tableExampleFiles = []tableExampleFile{{"main.go", 12}, {"go.mod", 1}, {"README.md", 4}, {"LICENSE", 2}}

// Column IDs stay stable while titles follow the gallery language.
func tableExampleColumns(c *ui.Context) []ui.TableColumn {
	return []ui.TableColumn{
		{Title: tr(c, "name"), ID: "name", Sortable: true, MinWidth: 120},
		{Title: tr(c, "size"), ID: "size", Width: 90, Align: ui.End, Sortable: true},
		{Title: tr(c, "kind"), ID: "kind", Width: 90},
	}
}

func controlTable(c *ui.Context, s *controlState, variant string) {
	if variant == "TableVirtualized" {
		controlTableVirtualized(c, s, variant)
		return
	}
	s.table.Sort = &s.tableSort
	files := slices.Clone(tableExampleFiles)
	slices.SortStableFunc(files, func(a, b tableExampleFile) int {
		order := 0
		switch s.tableSort.Column {
		case "name":
			order = strings.Compare(a.name, b.name)
		case "size":
			order = a.size - b.size
		}
		if s.tableSort.Descending {
			return -order
		}
		return order
	})
	components.Table(c, &s.table.TableState, tableExampleColumns(c), func() {
		for _, file := range files {
			if components.TableRow(c, file.name, file.name, s.choice == file.name, func(col int) {
				tableExampleCell(c, col, file.name, file.size)
			}).Clicked() {
				s.choice = file.name
			}
		}
	}).Label(tr(c, "files"))
}

const tableExampleCount = 10000

// Retained rows are re-sorted only when the user changes the order.
func controlTableVirtualized(c *ui.Context, s *controlState, _ string) {
	size := func(i int) int { return i*7919%977 + 1 }
	if s.tableRows == nil {
		s.tableRows = make([]int, tableExampleCount)
		for i := range s.tableRows {
			s.tableRows[i] = i
		}
	}
	s.table.Sort = &s.tableSort
	s.table.OnSort = func() {
		slices.SortStableFunc(s.tableRows, func(a, b int) int {
			order := a - b
			if s.tableSort.Column == "size" {
				order = size(a) - size(b)
			}
			if s.tableSort.Descending {
				return -order
			}
			return order
		})
	}
	s.table.List.Selected = &s.index
	s.table.List.Key = func(row int) any { return s.tableRows[row] }
	s.table.List.KeepOffset = true
	components.TableVirtualized(c, &s.table, tableExampleColumns(c), len(s.tableRows), func(row, col int) {
		i := s.tableRows[row]
		tableExampleCell(c, col, msg(c).FileD(i+1), size(i))
	}).Height(240).FillWidth().Label(tr(c, "largeProject"))
}

func tableExampleCell(c *ui.Context, col int, name string, size int) {
	switch col {
	case 0:
		components.Text(c, name).SingleLine()
	case 1:
		components.Text(c, msg(c).SizeKb(size)).SingleLine()
	default:
		components.Text(c, tr(c, "file")).SingleLine().TextColor(c.Theme().TextMuted)
	}
}
