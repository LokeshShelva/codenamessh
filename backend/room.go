package backend

import (
	"fmt"
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
)

type RoomUpdateMss struct {
	State any
}

type Room struct {
	Id      string
	Actions chan Action

	mu       sync.Mutex
	programs map[string]*tea.Program // player id -> their programs
	done     chan struct{}
}

// Clean a room if not user is present
// This is close the channel
func (r *Room) clean() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	isUserPresent := len(r.programs) != 0

	if !isUserPresent {
		close(r.done)
	}

	return !isUserPresent
}

func (r *Room) join(playerId string, program *tea.Program) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.programs[playerId] != nil {
		return fmt.Errorf("player %s already present in room", playerId)
	}
	r.programs[playerId] = program
	return nil
}

func (r *Room) leave(playerId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.programs[playerId] == nil {
		return fmt.Errorf("player %s is not in room.", playerId)
	}
	delete(r.programs, playerId)
	return nil
}

func (r *Room) run() {
	for {
		select {
		case action := <-r.Actions:
			log.Info("New Action", "playerId", action.PlayerID, "kind", action.Kind, "playload", action.Payload)

			// process action from user
			// ...

			// send updates to all connected clients
			r.broadcast()

		case <-r.done:
			return
		}
	}
}

func (r *Room) broadcast() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, p := range r.programs {
		p.Send(RoomUpdateMss{State: "helllo"})
		_ = id // TODO: remove
	}
}

func NewRoom(id string) *Room {
	return &Room{
		Id:       id,
		Actions:  make(chan Action),
		programs: make(map[string]*tea.Program),
		done:     make(chan struct{}),
	}
}
