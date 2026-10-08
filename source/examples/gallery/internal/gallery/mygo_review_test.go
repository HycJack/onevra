package gallery

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMyGoReviewDecisions(t *testing.T) {
	entries := []mygoComponent{{Name: "ZNew"}, {Name: "Button"}, {Name: "Link"}, {Name: "PrimaryButton"}, {Name: "Box"}, {Name: "ANew"}}
	pending := pendingMyGo(entries, []mygoExample{{name: "ZNew"}, {name: "Toast"}})
	var names []string
	for _, entry := range pending {
		names = append(names, entry.Name)
	}
	if !reflect.DeepEqual(names, []string{"ANew", "Toast", "ZNew"}) {
		t.Fatalf("pending = %v", names)
	}
	for _, name := range []string{"Button", "PrimaryButton", "Link", "MenuButton", "Combobox", "Autocomplete"} {
		if mygoHandled[name].Replacement == "" {
			t.Fatalf("missing decision for %s", name)
		}
		for _, entry := range mygoCatalog {
			if entry.name == name {
				t.Fatalf("handled example remains: %s", name)
			}
		}
		if pageFromKey("mygo/"+name) != lucidePage {
			t.Fatalf("handled route remains: %s", name)
		}
	}
	examples := unhandledMyGoExamples([]mygoExample{{name: "Button"}, {name: "Box"}, {name: "Menu"}, {name: "ContextMenu"}, {name: "Future"}})
	if len(examples) != 1 || examples[0].name != "Future" {
		t.Fatalf("filtered = %v", examples)
	}
}

func TestDirectUseComponentsHaveNoExamples(t *testing.T) {
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Status != mygoPrimitive {
			continue
		}
		if entry.HasExample {
			t.Fatalf("direct-use component has an example: %s", entry.Name)
		}
		for _, example := range mygoCatalog {
			if example.name == entry.Name {
				t.Fatalf("direct-use example remains: %s", entry.Name)
			}
		}
		if pageFromKey("mygo/"+entry.Name) != lucidePage {
			t.Fatalf("direct-use route remains: %s", entry.Name)
		}
	}
}

func TestToggleDecisions(t *testing.T) {
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		switch entry.Name {
		case "Toggle":
			if entry.Status != mygoReplaced || entry.Replacement != "components.ToggleButton" {
				t.Fatalf("wrong Toggle decision: %+v", entry)
			}
		case "ToggleGroup":
			if entry.Status != mygoExcluded {
				t.Fatalf("wrong ToggleGroup decision: %+v", entry)
			}
		}
	}
	for _, name := range []string{"Toggle", "ToggleGroup"} {
		if pageFromKey("mygo/"+name) != lucidePage {
			t.Fatalf("handled route remains: %s", name)
		}
	}
}

func TestExcludedComponents(t *testing.T) {
	for _, name := range []string{"Divider", "CheckboxGroup", "RadioGroup", "Segmented", "Avatar", "Meter", "BackButton", "ForwardButton", "Form", "Field", "Fieldset", "FindBar", "Toolbar", "Table", "OutlineTable", "TokenField"} {
		decision, ok := mygoHandled[name]
		if !ok || decision.Replacement != "" || decision.Reason == "" {
			t.Fatalf("missing exclusion for %s", name)
		}
		for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
			if entry.Name == name && (entry.Status != mygoExcluded || entry.HasExample) {
				t.Fatalf("incorrect exclusion: %+v", entry)
			}
		}
		for _, example := range mygoCatalog {
			if example.name == name {
				t.Fatalf("excluded example remains: %s", name)
			}
		}
		if pageFromKey("mygo/"+name) != lucidePage {
			t.Fatalf("excluded route remains: %s", name)
		}
	}
}

func TestMyGoReviewLog(t *testing.T) {
	var lines []string
	LogMyGoReview(func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) })
	output := strings.Join(lines, "\n")
	if !strings.Contains(output, mygoScanVersion) || !strings.Contains(output, "mygo_review.go") {
		t.Fatal(output)
	}
	for _, name := range []string{"Divider", "Button", "PrimaryButton", "Link", "MenuButton", "Toggle", "ToggleGroup", "CheckboxGroup", "RadioGroup", "Segmented", "Box", "Combobox", "Autocomplete", "Avatar", "Meter", "BackButton", "ForwardButton", "Form", "Field", "Fieldset", "FindBar", "Toolbar", "Table", "OutlineTable", "TokenField"} {
		if strings.Contains(output, "pending: "+name+" (") {
			t.Fatalf("unexpected reminder: %s", name)
		}
	}
	for _, name := range []string{"ContextMenu", "Menu"} {
		if strings.Contains(output, "pending: "+name+" (") {
			t.Fatalf("native menu should not need review: %s", name)
		}
	}
	seen := map[string]bool{}
	for _, entry := range mygoDiscovered {
		if seen[entry.Name] || entry.Source == "" {
			t.Fatalf("invalid inventory: %+v", entry)
		}
		seen[entry.Name] = true
	}
	for name := range mygoHandled {
		if !seen[name] {
			t.Fatalf("missing scanned component: %s", name)
		}
	}
}
