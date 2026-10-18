package components

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestAvatarInitials(t *testing.T) {
	for name, want := range map[string]string{
		"Maya Chen": "MC", "  leo  park  ": "LP", "Sara Ruiz Ortega": "SR", "ines": "I", "陈美雅": "陈", "": "",
	} {
		if got := avatarInitials(name); got != want {
			t.Errorf("%q: got %q, want %q", name, got, want)
		}
	}
}

func TestAvatarContentAndLabels(t *testing.T) {
	photo := ui.NewBitmap(image.NewNRGBA(image.Rect(0, 0, 8, 8)))
	var missing *ui.Bitmap
	tt := ui.NewTester(func(c *ui.Context) {
		c.SetTheme(theme.Palette(false))
		ui.Row(c).Gap(8).Children(func() {
			Avatar(c, "Maya Chen", AvatarOptions{Status: AvatarOnline, Size: 88})
			Avatar(c, "Leo Park", AvatarOptions{Image: photo, Status: AvatarOffline})
			Avatar(c, "Sara Ruiz", AvatarOptions{Image: missing})
		})
	}, 300, 120)
	for _, label := range []string{"Maya Chen, online", "Leo Park, offline", "Sara Ruiz"} {
		if _, ok := tt.Find(label); !ok {
			t.Fatalf("missing %q in %v", label, tt.Texts())
		}
	}
	if r, _ := tt.Find("Maya Chen, online"); r.W != 88 || r.H != 88 {
		t.Fatalf("explicit size not applied: %+v", r)
	}
	if !tt.HasText("MC") || !tt.HasText("SR") || tt.HasText("LP") {
		t.Fatalf("initials should replace only missing photos: %v", tt.Texts())
	}
	tt.SetLocale("zh-Hans")
	if _, ok := tt.Find("Maya Chen，在线"); !ok {
		t.Fatal("status label not localized")
	}
}

func avatarTeam(names ...string) []AvatarGroupMember {
	members := make([]AvatarGroupMember, len(names))
	for i, name := range names {
		members[i] = AvatarGroupMember{Name: name}
	}
	return members
}

func TestAvatarGroupOverflow(t *testing.T) {
	members := avatarTeam("Maya Chen", "Leo Park", "Sara Ruiz", "Tom Hale", "Ines Ma")
	options := AvatarGroupOptions{Max: 3, Label: "Editors"}
	tt := ui.NewTester(func(c *ui.Context) {
		AvatarGroup(c, members, options)
	}, 400, 120)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	for _, label := range []string{"Editors", "Maya Chen", "Sara Ruiz", "2 more members", "+2"} {
		if _, ok := tt.Find(label); !ok {
			t.Fatalf("missing %q in %v", label, tt.Texts())
		}
	}
	if _, ok := tt.Find("Tom Hale"); ok {
		t.Fatal("overflowing member was built")
	}
	// Four slots overlapping by a sixth of the diameter, after an equal inset.
	if r, _ := tt.Find("Editors"); r.H < 30 || r.W != r.H/6+4*(r.H-r.H/6) {
		t.Fatalf("stack bounds %+v", r)
	}

	members = members[:4]
	tt.Frame()
	if !tt.HasText("+1") || tt.HasText("+2") {
		t.Fatalf("overflow count not updated: %v", tt.Texts())
	}
	options.Max = 0
	tt.Frame()
	if _, ok := tt.Find("Tom Hale"); !ok || tt.HasText("+1") {
		t.Fatalf("default maximum should show four members: %v", tt.Texts())
	}
	options = AvatarGroupOptions{Max: -1}
	tt.SetLocale("zh-Hans")
	if _, ok := tt.Find("团队成员"); !ok || !tt.HasText("+4") {
		t.Fatalf("negative maximum or localized label: %v", tt.Texts())
	}
	if _, ok := tt.Find("另外 4 位成员"); !ok {
		t.Fatal("overflow label not localized")
	}
}

func TestAvatarGroupHover(t *testing.T) {
	var lifted ui.Rect
	members := avatarTeam("Maya Chen", "Leo Park", "Sara Ruiz")
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(60).Children(func() {
			AvatarGroup(c, members, AvatarGroupOptions{})
		})
	}, 300, 200)
	r, _ := tt.Find("Leo Park")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Advance(650 * time.Millisecond)
	tt.Frame()
	if !tt.HasText("Leo Park") {
		t.Fatal("hovered member did not name itself")
	}
	lifted, _ = tt.Find("Leo Park")
	if lifted.Y >= r.Y {
		t.Fatalf("hovered member did not rise: %+v from %+v", lifted, r)
	}
	if first, _ := tt.Find("Maya Chen"); first.X >= 60 {
		t.Fatalf("preceding member did not ease aside: %+v", first)
	}
}
