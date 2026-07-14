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

// This is the room in the server.
// Keeps track all the joined players (clients) in the room.
// Processing and sending updates to the connected clients
// Does not manage the game state. Deligates that to the game
// struct.
type Room struct {
	// Id of the room
	Id string

	// Listens on this for messagae from the clients
	Actions chan Action

	// shared mutext for map access
	mu sync.Mutex

	// Map of connected client `tea.Program` to the playerId
	programs map[string]*tea.Program

	// Listens on this for closing the room
	done chan struct{}
}

// Send updates to all the connected clients.
// TODO: decide on what each update contains (mostly a game state snapshot)
func (r *Room) broadcast() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, p := range r.programs {
		p.Send(RoomUpdateMss{State: "helllo"})
		_ = id // TODO: remove
	}
}

// Clean a room if no client is connected. Closes all the channels
func (r *Room) clean() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	isUserPresent := len(r.programs) != 0

	if !isUserPresent {
		close(r.done)
	}

	return !isUserPresent
}

// Join this room.
func (r *Room) join(playerId string, program *tea.Program) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.programs[playerId] != nil {
		return fmt.Errorf("player %s already present in room", playerId)
	}
	r.programs[playerId] = program
	return nil
}

// Leave this room
func (r *Room) leave(playerId string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.programs[playerId] == nil {
		return fmt.Errorf("player %s is not in room.", playerId)
	}
	delete(r.programs, playerId)
	return nil
}

// Main processing method that listens for updates and
// delegates to the game struct.
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

func NewRoom(id string) *Room {
	return &Room{
		Id:       id,
		Actions:  make(chan Action),
		programs: make(map[string]*tea.Program),
		done:     make(chan struct{}),
	}
}
