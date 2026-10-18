package gallery

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// avatarPeople are sample members; person names stay untranslated.
var avatarPeople = []string{"Maya Chen", "Leo Park", "Sara Ruiz", "Tom Hale", "Ines Ma", "Omar Diaz", "Ava Kim", "Noah Berg"}

func controlAvatar(c *ui.Context, _ *controlState, variant string) {
	photo, _ := galleryImage()
	ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
		switch variant {
		case "Initials":
			for i, size := range []float32{28, 0, 48, 88} {
				components.Avatar(c, avatarPeople[i], components.AvatarOptions{Size: size})
			}
		case "Photo":
			for _, size := range []float32{28, 0, 48, 88} {
				components.Avatar(c, avatarPeople[0], components.AvatarOptions{Image: photo, Size: size})
			}
		case "Status":
			components.Avatar(c, avatarPeople[0], components.AvatarOptions{Image: photo, Status: components.AvatarOnline, Size: 48})
			components.Avatar(c, avatarPeople[1], components.AvatarOptions{Status: components.AvatarOffline, Size: 48})
			components.Avatar(c, avatarPeople[2], components.AvatarOptions{Status: components.AvatarOnline})
			components.Avatar(c, avatarPeople[3], components.AvatarOptions{Status: components.AvatarOnline, Size: 88})
		}
	})
}

func controlAvatarGroup(c *ui.Context, s *controlState, variant string) {
	photo, _ := galleryImage()
	members := func(n int) []components.AvatarGroupMember {
		list := make([]components.AvatarGroupMember, n)
		for i := range list {
			list[i] = components.AvatarGroupMember{Name: avatarPeople[i]}
		}
		if n > 0 {
			list[0].Image, list[0].Status = photo, components.AvatarOnline
		}
		return list
	}
	switch variant {
	case "Default":
		components.AvatarGroup(c, members(3), components.AvatarGroupOptions{})
	case "Overflow":
		ui.Column(c).Gap(12).Children(func() {
			for _, size := range []float32{28, 0, 48} {
				components.AvatarGroup(c, members(len(avatarPeople)), components.AvatarGroupOptions{Max: 3, Size: size, Label: tr(c, "avatarEditors")})
			}
		})
	case "Changing members":
		// index counts the members shown, starting from the first visit.
		if s.index == 0 {
			s.index = 5
		}
		ui.Column(c).Gap(12).Children(func() {
			components.AvatarGroup(c, members(s.index), components.AvatarGroupOptions{Max: 4, Label: tr(c, "avatarEditors")})
			ui.Row(c).Gap(8).Children(func() {
				if components.Button(c, tr(c, "addMember"), components.ButtonSecondary).Disabled(s.index >= len(avatarPeople)).Clicked() {
					s.index++
				}
				if components.Button(c, tr(c, "removeMember"), components.ButtonSecondary).Disabled(s.index <= 1).Clicked() {
					s.index--
				}
			})
		})
	}
}
