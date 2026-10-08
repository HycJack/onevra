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
	for _, name := range []string{"Button", "PrimaryButton", "Link", "MenuButton"} {
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
	examples := unhandledMyGoExamples([]mygoExample{{name: "Button"}, {name: "Box"}, {name: "Future"}})
	if len(examples) != 2 || examples[0].name != "Box" || examples[1].name != "Future" {
		t.Fatalf("filtered = %v", examples)
	}
}

func TestMyGoReviewLog(t *testing.T) {
	var lines []string
	LogMyGoReview(func(format string, args ...any) { lines = append(lines, fmt.Sprintf(format, args...)) })
	output := strings.Join(lines, "\n")
	if !strings.Contains(output, mygoScanVersion) || !strings.Contains(output, "mygo_review.go") {
		t.Fatal(output)
	}
	for _, name := range []string{"Button", "PrimaryButton", "Link", "MenuButton", "Box"} {
		if strings.Contains(output, "pending: "+name+" (") {
			t.Fatalf("unexpected reminder: %s", name)
		}
	}
	if !strings.Contains(output, "pending: Toggle (") {
		t.Fatal(output)
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
