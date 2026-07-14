package backend

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type Actionkind int

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

func (p *Player) SetRoom(room *Room) {
	p.room = room
}

func (u *Player) Init() tea.Cmd {
	return nil
}

func (u *Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch m.String() {
		case "ctrl+c":
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

func NewPlayer(playerId string) *Player {
	return &Player{
		id:                playerId,
		broadCastMsgCount: 0,
	}
}
