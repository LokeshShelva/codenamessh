package backend

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textinput"
	"github.com/google/uuid"
)

type focusedInput int

type User struct {
	Id   string
	Name string

	Room *Room

	ui *UserUI
}

type UserUI struct {
	nameInput textinput.Model
	roomInput textinput.Model
	focused   focusedInput

	help help.Model
}

func (u *User) JoinRoom(room *Room) error {
	if u.Room != nil {
		return fmt.Errorf("connot join new room '%s'. user is already part room '%s'", room.Id, u.Room.Id)
	}
	u.Room = room
	return nil
}

func newInput(placeholder string, focus bool) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 50
	ti.SetWidth(20)

	if focus {
		ti.Focus()
	}

	return ti
}

func CreateUser(name string) *User {
	h := help.New()
	h.SetWidth(40)

	return &User{
		Id:   uuid.New().String(),
		Name: name,
		ui: &UserUI{
			nameInput: newInput(name, true),
			roomInput: newInput("", false),
			help:      h,
		},
	}
}
