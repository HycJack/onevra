package gallery

import (
	"sort"
)

//go:generate go run ../cmd/mygocomponents

type mygoDecision struct {
	Replacement string // Empty means excluded, not awaiting implementation.
	Reason      string
}

// Record decisions here, never in the generated inventory. Both replacement and
// exclusion decisions hide upstream examples and suppress pending reminders.
// To exclude a component, add {Reason: "why this component is not needed"}.
var mygoHandled = map[string]mygoDecision{
	"Divider":          {Reason: "A plain theme-colored separator; compose a Box or container borders instead of a dedicated Onevra component"},
	"Badge":            {Replacement: "components.Badge", Reason: "Onevra compact status label"},
	"RichText":         {Replacement: "components.RichText", Reason: "Onevra inline text composition on the MyGo text primitive"},
	"Breadcrumbs":      {Replacement: "components.Breadcrumbs", Reason: "Onevra path navigation with application-owned activation"},
	"Rating":           {Replacement: "components.Rating", Reason: "Onevra star rating with matching selected fill and stroke"},
	"Split":            {Replacement: "components.Split", Reason: "Unified split panes with a neutral, hover-free divider"},
	"SplitVertical":    {Replacement: "components.Split", Reason: "Vertical Split orientation"},
	"Tree":             {Replacement: "components.Tree", Reason: "Onevra tree with neutral selection, no hover face and double-click disclosure"},
	"TreeItem":         {Replacement: "components.Tree", Reason: "Onevra TreeItem rows"},
	"Outline":          {Replacement: "components.TreeVirtualized", Reason: "Onevra virtualized tree with neutral selection and double-click disclosure"},
	"Progress":         {Replacement: "components.Progress", Reason: "Onevra determinate and indeterminate progress bar"},
	"Spinner":          {Replacement: "components.ProgressRing", Reason: "Onevra circular progress and indeterminate activity"},
	"Avatar":           {Reason: "No Onevra Avatar component is needed"},
	"Meter":            {Reason: "No Onevra Meter component is needed"},
	"BackButton":       {Reason: "No Onevra BackButton component is needed"},
	"ForwardButton":    {Reason: "No Onevra ForwardButton component is needed"},
	"Form":             {Reason: "No Onevra Form component is needed"},
	"Field":            {Reason: "No Onevra Field component is needed"},
	"Fieldset":         {Reason: "No Onevra Fieldset component is needed"},
	"FindBar":          {Reason: "No Onevra FindBar component is needed"},
	"Toolbar":          {Reason: "No Onevra Toolbar component is needed"},
	"Table":            {Reason: "No Onevra Table component is needed"},
	"OutlineTable":     {Reason: "No Onevra OutlineTable component is needed"},
	"TokenField":       {Reason: "No Onevra TokenField component is needed"},
	"Scroll":           {Replacement: "components.ScrollView", Reason: "Unified scroll container with vertical scrolling by default"},
	"ScrollHorizontal": {Replacement: "components.ScrollView", Reason: "Horizontal ScrollView axes"},
	"ScrollBoth":       {Replacement: "components.ScrollView", Reason: "Both ScrollView axes"},
	"List":             {Replacement: "components.List", Reason: "Onevra virtual list with neutral selection and no hover face"},
	"Sidebar":          {Replacement: "components.Sidebar", Reason: "Onevra grouped navigation with application-owned selection"},
	"SidebarSection":   {Replacement: "components.Sidebar", Reason: "Onevra sidebar section headings"},
	"SidebarItem":      {Replacement: "components.Sidebar", Reason: "Onevra sidebar navigation items"},

	"Segmented":       {Reason: "No Onevra segmented control component is needed"},
	"CheckboxGroup":   {Reason: "No separate Onevra checkbox group component is needed"},
	"RadioGroup":      {Reason: "No separate Onevra radio group component is needed; upstream grouping remains a composition building block"},
	"Toggle":          {Replacement: "components.ToggleButton", Reason: "Onevra toggle button styling"},
	"ToggleGroup":     {Reason: "Grouped toggle styling is not needed; compose independent ToggleButtons with layout primitives"},
	"Checkbox":        {Replacement: "components.Checkbox", Reason: "Onevra checkbox styling"},
	"Switch":          {Replacement: "components.Switch", Reason: "Onevra switch styling"},
	"Radio":           {Replacement: "components.Radio", Reason: "Onevra radio styling"},
	"NumberInput":     {Replacement: "components.NumberInput", Reason: "Onevra numeric text field with integrated vertical step buttons"},
	"Stepper":         {Reason: "Step buttons are integrated into NumberInput; no standalone Stepper is needed"},
	"TextInputBuffer": {Replacement: "components.TextBoxBuffer", Reason: "Onevra single-line input with application-owned indexed storage"},
	"TextAreaBuffer":  {Replacement: "components.TextBoxBuffer", Reason: "Onevra multiline input with application-owned indexed storage"},
	"TextInput":       {Replacement: "components.TextBox", Reason: "Onevra single-line text input"},
	"TextArea":        {Replacement: "components.TextBox", Reason: "Onevra multiline text input"},
	"SearchField":     {Replacement: "components.TextBox", Reason: "TextBox with a leading search icon; applications own search and clearing behavior"},
	"Select":          {Replacement: "components.Select", Reason: "Onevra select parts with application-owned options"},
	"Combobox":        {Replacement: "components.Select", Reason: "Searchable Select with application-owned filtering"},
	"Autocomplete":    {Replacement: "components.Select", Reason: "Searchable Select with application-owned suggestions"},
	"Slider":          {Replacement: "components.Slider", Reason: "Onevra continuous slider"},
	"StepSlider":      {Replacement: "components.Slider", Reason: "Slider Step property enables snapping and tick marks"},
	"RangeSlider":     {Replacement: "components.Slider", Reason: "Slider Range property enables two independently focusable thumbs"},
	"Accordion":       {Replacement: "components.Accordion", Reason: "Onevra disclosure styling without hover"},
	"AccordionItem":   {Replacement: "components.Accordion", Reason: "AccordionGroup.Item builds Onevra sections"},
	"Collapsible":     {Replacement: "components.Collapsible", Reason: "Onevra disclosure with configurable adjacent arrow"},
	"Tabs":            {Replacement: "components.Tabs", Reason: "Onevra tab styling"},
	"Text":            {Replacement: "components.Text", Reason: "Onevra text defaults"},
	"EditableText":    {Replacement: "components.Text", Reason: "Text with an Editable binding; undecorated inline editing"},
	"Textf":           {Replacement: "components.Textf", Reason: "Formatted Onevra text"},
	"Icon":            {Replacement: "components.Icon", Reason: "Onevra icon component with named asset accessors"},
	"Image":           {Replacement: "components.Image", Reason: "Onevra image component with application-owned decoded assets"},

	"Button":        {Replacement: "components.Button", Reason: "Onevra button styles"},
	"MenuButton":    {Replacement: "components.Button / Menu", Reason: "Compose components.Button(...).Menu(...) as demonstrated under Menus; no separate component"},
	"PrimaryButton": {Replacement: "components.Button / ButtonPrimary", Reason: "A button variant, not a separate component"},
	"Link":          {Replacement: "components.Button / ButtonLink", Reason: "An inline button variant"},
}

// These primitives remain composition building blocks, without standalone examples.
// Keep this list explicit so new APIs still prompt review. Base constructors are
// skipped by the scanner itself and are not inventoried as styled components.
var mygoPrimitives = map[string]bool{
	"Box": true, "Column": true, "Row": true, "Spacer": true,
	"Grid": true, "Overlay": true,
}

// Element methods are not found by the free-function scanner. These native menu
// APIs are deliberately used directly, including by the Onevra gallery examples.
var mygoNativeMenus = map[string]string{
	"ContextMenu": "ui/menu.go",
	"Menu":        "ui/menu.go",
}

type mygoComponent struct{ Name, Source string }

type mygoStatus string

const (
	mygoReplaced   mygoStatus = "Replaced by Onevra"
	mygoExcluded   mygoStatus = "Not needed"
	mygoUpstream   mygoStatus = "Still using MyGo"
	mygoUnreviewed mygoStatus = "Awaiting review"
	mygoPrimitive  mygoStatus = "Use MyGo directly"
)

type mygoReviewEntry struct {
	mygoComponent
	Status              mygoStatus
	Replacement, Reason string
	HasExample          bool
}

// The catalog supplements discovery with method-based UI and plugins. An API
// discovered in source is not proof of application usage: only catalog entries
// are labeled as upstream examples; others explicitly await review.
func mygoReviewEntries(discovered []mygoComponent, examples []mygoExample, decisions map[string]mygoDecision) []mygoReviewEntry {
	found := make(map[string]mygoComponent, len(discovered))
	for _, entry := range discovered {
		found[entry.Name] = entry
	}
	for name, source := range mygoNativeMenus {
		found[name] = mygoComponent{name, source}
	}
	demonstrated := map[string]bool{}
	for _, entry := range examples {
		demonstrated[entry.name] = true
		if _, ok := found[entry.name]; !ok {
			found[entry.name] = mygoComponent{entry.name, "gallery catalog"}
		}
	}
	for name := range decisions {
		if _, ok := found[name]; !ok {
			found[name] = mygoComponent{name, "Not found in current inventory"}
		}
	}
	entries := make([]mygoReviewEntry, 0, len(found))
	for name, component := range found {
		entry := mygoReviewEntry{mygoComponent: component, Status: mygoUnreviewed, HasExample: demonstrated[name]}
		if entry.HasExample {
			entry.Status = mygoUpstream
		}
		if mygoPrimitives[name] {
			entry.Status, entry.HasExample = mygoPrimitive, false
		}
		if _, native := mygoNativeMenus[name]; native {
			entry.Status, entry.HasExample = mygoPrimitive, false
			entry.Reason = "Use MyGo native menus directly; no Onevra menu implementation"
		}
		if decision, handled := decisions[name]; handled {
			entry.Status, entry.Replacement, entry.Reason, entry.HasExample = mygoExcluded, decision.Replacement, decision.Reason, false
			if decision.Replacement != "" {
				entry.Status = mygoReplaced
			}
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	return entries
}

func pendingMyGo(discovered []mygoComponent, examples []mygoExample) []mygoComponent {
	var pending []mygoComponent
	for _, entry := range mygoReviewEntries(discovered, examples, mygoHandled) {
		if entry.Status == mygoUpstream || entry.Status == mygoUnreviewed {
			pending = append(pending, entry.mygoComponent)
		}
	}
	return pending
}

// LogMyGoReview reports outstanding decisions once at application startup, never
// during rendering. Pending decisions are reminders, not build or runtime errors.
func LogMyGoReview(printf func(string, ...any)) {
	pending := pendingMyGo(mygoDiscovered, mygoExamples)
	printf("[MyGo review] %s: %d pending; see Component status or record decisions in mygo_review.go", mygoScanVersion, len(pending))
	for _, entry := range mygoReviewEntries(mygoDiscovered, mygoExamples, mygoHandled) {
		if entry.Status == mygoUpstream || entry.Status == mygoUnreviewed {
			detail := entry.Source
			if !entry.HasExample {
				detail += "; no gallery example"
			}
			printf("[MyGo review] pending: %s (%s)", entry.Name, detail)
		}
		if entry.Source == "Not found in current inventory" {
			printf("[MyGo review] decision no longer matches inventory: %s", entry.Name)
		}
	}
}

func unhandledMyGoExamples(examples []mygoExample) []mygoExample {
	result := make([]mygoExample, 0, len(examples))
	for _, entry := range examples {
		_, native := mygoNativeMenus[entry.name]
		if _, handled := mygoHandled[entry.name]; !handled && !mygoPrimitives[entry.name] && !native {
			result = append(result, entry)
		}
	}
	return result
}
