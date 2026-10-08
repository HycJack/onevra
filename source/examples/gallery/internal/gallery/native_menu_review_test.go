package gallery

import "testing"

func TestNativeMenuReview(t *testing.T) {
	seen := map[string]bool{}
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Name != "ContextMenu" && entry.Name != "Menu" {
			continue
		}
		seen[entry.Name] = true
		if entry.Status != mygoPrimitive || entry.Replacement != "" || entry.Reason == "" || entry.Source != "ui/menu.go" {
			t.Fatalf("native menu must use MyGo directly: %+v", entry)
		}
	}
	if len(seen) != 2 {
		t.Fatalf("missing native menu decisions: %v", seen)
	}
	found := false
	for _, example := range mygoCatalog {
		if example.name == "ContextMenu" {
			found = true
		}
	}
	if !found {
		t.Fatal("native context menu example was removed")
	}
}
