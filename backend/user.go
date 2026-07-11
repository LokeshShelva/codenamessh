package backend

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	"github.com/google/uuid"
)

type User struct {
	Id   string
	Name string

	Room *Room

	textinput textinput.Model
	err       error
}

type UserUI struct {
	textInput textinput.Model
	err       error
}

func (u *User) JoinRoom(room *Room) error {
	if u.Room != nil {
		return fmt.Errorf("connot join new room '%s'. user is already part room '%s'", room.Id, u.Room.Id)
	}
	u.Room = room
	return nil
}

func CreateUser(name string) *User {
	ti := textinput.New()
	ti.Placeholder = name
	// ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 50
	ti.SetWidth(20)

	return &User{
		Id:        uuid.New().String(),
		Name:      name,
		err:       nil,
		textinput: ti,
	}
}
