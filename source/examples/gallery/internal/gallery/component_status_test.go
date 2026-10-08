package gallery

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

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

func TestComponentStatusCards(t *testing.T) {
	entries := mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled)
	cards := componentStatusCards(entries, mygoReplaced)
	if len(cards) != 29 {
		t.Fatalf("got %d replacement cards", len(cards))
	}
	for _, card := range cards {
		if card.Name == "Button" && !reflect.DeepEqual(card.Replaces, []string{"Button", "Link", "MenuButton", "PrimaryButton"}) {
			t.Fatalf("button variants not grouped: %+v", card)
		}
		if card.Name == "TextBox" && !reflect.DeepEqual(card.Replaces, []string{"SearchField", "TextArea", "TextInput"}) {
			t.Fatalf("TextBox modes not grouped: %+v", card)
		}
		if card.Name == "TextBoxBuffer" && !reflect.DeepEqual(card.Replaces, []string{"TextAreaBuffer", "TextInputBuffer"}) {
			t.Fatalf("TextBoxBuffer modes not grouped: %+v", card)
		}
		if card.Name == "Select" && !reflect.DeepEqual(card.Replaces, []string{"Autocomplete", "Combobox", "Select"}) {
			t.Fatalf("Select replacements not grouped: %+v", card)
		}
	}
	foundSplit := false
	for _, card := range cards {
		if card.Name == "Split" {
			foundSplit = true
			if !reflect.DeepEqual(card.Replaces, []string{"Split", "SplitVertical"}) {
				t.Fatalf("split orientations not grouped: %+v", card)
			}
		}
	}
	if !foundSplit {
		t.Fatal("missing Split replacement")
	}
	foundScrollView := false
	for _, card := range cards {
		if card.Name == "ScrollView" {
			foundScrollView = true
			if !reflect.DeepEqual(card.Replaces, []string{"Scroll", "ScrollBoth", "ScrollHorizontal"}) {
				t.Fatalf("scroll variants not grouped: %+v", card)
			}
		}
	}
	if !foundScrollView {
		t.Fatal("missing ScrollView replacement")
	}
	for _, status := range []mygoStatus{mygoExcluded, mygoUpstream, mygoUnreviewed, mygoPrimitive} {
		for _, card := range componentStatusCards(entries, status) {
			if len(card.Replaces) != 0 {
				t.Fatalf("unexpected replacement names: %+v", card)
			}
		}
	}
}

func TestComponentStatusPage(t *testing.T) {
	for _, width := range []int{820, 360} {
		for _, dark := range []bool{false, true} {
			g := &Gallery{page: componentStatusPage}
			tt := ui.NewTester(g.View, width, 900)
			tt.SetDark(dark)
			texts := tt.Texts()
			excluded := slices.Index(texts, string(mygoExcluded))
			if excluded < 0 {
				t.Fatal("missing Not needed group")
			}
			for _, status := range []mygoStatus{mygoReplaced, mygoUpstream, mygoUnreviewed, mygoPrimitive} {
				if index := slices.Index(texts, string(status)); index >= excluded {
					t.Fatalf("%s appears after Not needed", status)
				}
			}
			for _, label := range []string{"Replaced by Onevra", "Button", "Button, Link, MenuButton, PrimaryButton", "Use MyGo directly", "Box", "Row", "Column", "Spacer", "Grid", "Overlay"} {
				if !tt.HasText(label) {
					t.Fatalf("width=%d dark=%t missing %q", width, dark, label)
				}
			}
			for _, label := range []string{"Component status", "Not needed (0)", "No components in this group.", "Replacement: components.Button", "Onevra button styles", "ui/widgets.go"} {
				if tt.HasText(label) {
					t.Fatalf("unexpected metadata: %s", label)
				}
			}
			button, _ := tt.Find("Badge")
			checkbox, _ := tt.Find("Breadcrumbs")
			// Wrapped replacement names can shift the headings within a row.
			if width == 820 && (checkbox.X <= button.X || checkbox.Y >= button.Y+button.H || button.Y >= checkbox.Y+checkbox.H) {
				t.Fatalf("wide layout is not a grid: button=%+v checkbox=%+v", button, checkbox)
			}
			if width == 360 && checkbox.Y <= button.Y {
				t.Fatal("narrow layout did not wrap")
			}
			if tt.Image() == nil {
				t.Fatal("missing rendered status page")
			}
			title := ui.NewTester(g.Titlebar, width, 60)
			if !title.HasText("Component status") {
				t.Fatal("wrong title")
			}
		}
	}
	g := &Gallery{}
	nav := ui.NewTester(g.Left, 220, 1800)
	if err := nav.Click("Component status"); err != nil {
		t.Fatal(err)
	}
	if g.page != componentStatusPage || pageFromKey(pageKey(g.page)) != componentStatusPage {
		t.Fatal("status route not selected or persistent")
	}
	if _, ok := controlExampleFor(componentStatusPage); ok {
		t.Fatal("status route overlaps controls")
	}
	if _, ok := mygoExampleFor(componentStatusPage); ok {
		t.Fatal("status route overlaps MyGo")
	}
}
