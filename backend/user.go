package backend

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Actionkind int

const (
	LeaveRoomAct Actionkind = 0
)

type Action struct {
	PlayerID string
	Kind     Actionkind
	Payload  any
}

type Player struct {
	id                string
	room              *Room
	broadCastMsgCount int
}

func (u *Player) Init() tea.Cmd {
	return nil
}

func (u *Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "ctrl+c":
			u.room.Actions <- Action{PlayerID: u.id, Kind: LeaveRoomAct}
			return u, tea.Quit
		case "a":
			u.room.Actions <- Action{PlayerID: u.id, Kind: 1}
		}

	case RoomUpdateMss:
		u.broadCastMsgCount += 1
	}

	return u, nil
}

func (u *Player) View() tea.View {
	return tea.NewView(fmt.Sprintf("Hello from %s\n\nBroadCast Count: %d", u.id, u.broadCastMsgCount))
}

func NewPlayer(playerId string, room *Room) *Player {
	return &Player{
		id:                playerId,
		room:              room,
		broadCastMsgCount: 0,
	}
}
