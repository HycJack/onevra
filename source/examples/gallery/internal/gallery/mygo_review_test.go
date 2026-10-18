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
	if !reflect.DeepEqual(names, []string{"ANew", "ZNew"}) {
		t.Fatalf("pending = %v", names)
	}
	for _, name := range []string{"Tooltip", "Toast", "Button", "PrimaryButton", "Link", "MenuButton", "Combobox", "Autocomplete", "Tabs", "Segmented"} {
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

func TestDirectUseComponentsHaveOnlyExplicitExamples(t *testing.T) {
	plugins := map[string]bool{}
	for _, example := range mygoExamples {
		plugins[example.name] = example.group == "Plugins"
	}
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Status != mygoPrimitive {
			continue
		}
		if _, direct := mygoDirectExamples[entry.Name]; plugins[entry.Name] || direct {
			if !entry.HasExample || pageFromKey("mygo/"+entry.Name) == lucidePage {
				t.Fatalf("direct-use demonstration missing: %s", entry.Name)
			}
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

func TestMyGoPluginsAreDirectUse(t *testing.T) {
	examples := []mygoExample{
		{name: "Glass", group: "Plugins"},
		{name: "Terminal", group: "Plugins"},
		{name: "FuturePlugin", group: "Plugins"},
	}
	entries := mygoReviewEntries(nil, examples, nil)
	for _, example := range examples {
		found := false
		for _, entry := range entries {
			if entry.Name == example.name {
				found = true
				if entry.Status != mygoPrimitive || !entry.HasExample || entry.Reason == "" {
					t.Fatalf("plugin is not direct-use: %+v", entry)
				}
			}
		}
		if !found {
			t.Fatalf("missing plugin: %s", example.name)
		}
	}
	if pending := pendingMyGo(nil, examples); len(pending) != 0 {
		t.Fatalf("plugins pending review: %v", pending)
	}
	if got := unhandledMyGoExamples(examples); len(got) != len(examples) {
		t.Fatalf("plugin examples filtered: %v", got)
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
	for _, name := range []string{"Divider", "CheckboxGroup", "RadioGroup", "Meter", "BackButton", "ForwardButton", "Form", "Field", "Fieldset", "FindBar", "Toolbar", "OutlineTable", "TokenField", "Installed apps"} {
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
	for _, name := range []string{"Toast", "Tooltip", "Divider", "Button", "PrimaryButton", "Link", "MenuButton", "Toggle", "ToggleGroup", "CheckboxGroup", "RadioGroup", "Box", "Combobox", "Autocomplete", "Avatar", "Meter", "BackButton", "ForwardButton", "Form", "Field", "Fieldset", "FindBar", "Toolbar", "OutlineTable", "TokenField"} {
		if strings.Contains(output, "pending: "+name+" (") {
			t.Fatalf("unexpected reminder: %s", name)
		}
	}
	for _, name := range []string{"ContextMenu", "Menu", "Glass", "Terminal", "WebView", "Flyout"} {
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
	// Method-based components are supplemented by the catalog, not the scanner.
	for _, example := range mygoExamples {
		seen[example.name] = true
	}
	for name := range mygoSubpackages {
		seen[name] = true
	}
	for name := range mygoHandled {
		if !seen[name] {
			t.Fatalf("missing scanned component: %s", name)
		}
	}
}

func TestComponentStatuses(t *testing.T) {
	decisions := map[string]mygoDecision{
		"Replaced": {Replacement: "components.Button", Reason: "Button variant"},
		"Excluded": {Reason: "Not needed in Onevra"},
	}
	entries := mygoReviewEntries([]mygoComponent{
		{Name: "Replaced"}, {Name: "Excluded"}, {Name: "Future"}, {Name: "Box"},
	}, []mygoExample{{name: "Replaced"}, {name: "Excluded"}, {name: "Upstream"}}, decisions)
	want := map[string]mygoStatus{"Replaced": mygoReplaced, "Excluded": mygoExcluded, "Future": mygoUnreviewed, "Box": mygoPrimitive, "Upstream": mygoUpstream, "ContextMenu": mygoPrimitive, "Menu": mygoPrimitive}
	for i, entry := range entries {
		if entry.Status != want[entry.Name] {
			t.Fatalf("wrong status: %+v", entry)
		}
		if i > 0 && entries[i-1].Name >= entry.Name {
			t.Fatal("entries not sorted")
		}
		if entry.Status == mygoExcluded || entry.Status == mygoReplaced {
			if entry.HasExample || entry.Reason == "" {
				t.Fatalf("invalid handled entry: %+v", entry)
			}
		}
		delete(want, entry.Name)
	}
	if len(want) != 0 {
		t.Fatalf("missing statuses: %v", want)
	}
	stale := mygoReviewEntries(nil, nil, decisions)
	if len(stale) != 4 || stale[1].Name != "Excluded" || stale[1].Source != "Not found in current inventory" {
		t.Fatal("lost stale decisions")
	}
	for name, decision := range mygoHandled {
		if strings.TrimSpace(decision.Reason) == "" {
			t.Fatalf("missing reason for %s", name)
		}
		for _, example := range mygoCatalog {
			if example.name == name {
				t.Fatalf("handled example remains: %s", name)
			}
		}
	}
}

func TestTerminalUsesMyGoDirectly(t *testing.T) {
	found := false
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Name == "Terminal" {
			found = true
			if entry.Status != mygoPrimitive || !entry.HasExample || entry.Reason == "" {
				t.Fatalf("unexpected terminal status: %+v", entry)
			}
		}
	}
	if !found {
		t.Fatal("Terminal missing from review inventory")
	}
	for _, entry := range pendingMyGo(mygoDiscovered, mygoExamples) {
		if entry.Name == "Terminal" {
			t.Fatal("Terminal still awaiting review")
		}
	}
	for _, example := range mygoCatalog {
		if example.name == "Terminal" && example.group == "Plugins" {
			return
		}
	}
	t.Fatal("Terminal plugin demonstration missing")
}

func TestFocusedValueUsesMyGoDirectly(t *testing.T) {
	entries := mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled)
	found := false
	for _, entry := range entries {
		if entry.Name == "FocusedValue" {
			found = true
			if entry.Status != mygoPrimitive || entry.HasExample {
				t.Fatalf("unexpected focus API status: %+v", entry)
			}
		}
	}
	if !found {
		t.Fatal("FocusedValue missing from review inventory")
	}
	for _, entry := range pendingMyGo(mygoDiscovered, mygoExamples) {
		if entry.Name == "FocusedValue" {
			t.Fatal("FocusedValue still awaiting review")
		}
	}
	for _, example := range mygoCatalog {
		if example.name == "FocusedValue" {
			t.Fatal("FocusedValue still has a standalone example")
		}
	}
}

func TestTextfUsesMyGoDirectly(t *testing.T) {
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Name == "Textf" {
			if entry.Status != mygoPrimitive || entry.Replacement != "" || entry.HasExample {
				t.Fatalf("unexpected formatted text status: %+v", entry)
			}
			return
		}
	}
	t.Fatal("Textf missing from review inventory")
}
