package backend

import (
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	Id    uuid.UUID
	Users []*User
}

func (r *Room) AddUser(user *User) error {
	for _, u := range r.Users {
		if u.Id == user.Id {
			return fmt.Errorf("user(id='%s', name='%s') already in room '%s'", user.Id, user.Name, r.Id)
		}
	}
	r.Users = append(r.Users, user)
	return nil
}

func CreateRoom() *Room {
	return &Room{
		Id: uuid.New(),
	}
}
