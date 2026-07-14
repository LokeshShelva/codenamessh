package backend

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
)

type RoomUpdateMss struct {
	State any
}

type Room struct {
	Actions  chan Action
	programs map[string]*tea.Program // player id -> their programs
}

func (r *Room) broadcast() {
	for id, p := range r.programs {
		p.Send(RoomUpdateMss{State: "helllo"})
		_ = id // TODO: remove
	}
}

func (r *Room) Run() {
	for action := range r.Actions {
		log.Info("New Action", "playerId", action.PlayerID, "kind", action.Kind, "playload", action.Payload)

		// process action from user
		switch action.Kind {
		case LeaveRoomAct:
			if err := r.leave(action.PlayerID); err != nil {
				log.Error("failed to leave room", "playerId", action.PlayerID)
			}
			log.Info("player left room", "playerId", action.PlayerID)
		}

		// send updates to all connected clients
		r.broadcast()
	}
}

func (r *Room) Join(playerId string, program *tea.Program) error {
	if r.programs[playerId] != nil {
		return fmt.Errorf("player %s already present in room", playerId)
	}
	r.programs[playerId] = program
	log.Info("player added to room", "playerId", playerId)
	return nil
}

func (r *Room) leave(playerId string) error {
	if r.programs[playerId] == nil {
		return fmt.Errorf("player %s is not in room.", playerId)
	}
	delete(r.programs, playerId)
	return nil
}

func NewRoom() *Room {
	return &Room{
		Actions:  make(chan Action),
		programs: make(map[string]*tea.Program),
	}
}
