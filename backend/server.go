package backend

import (
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
)

// Server manages connections and rooms.
type Server struct {
	// shared mutex lock to access the rooms map
	mu sync.Mutex

	// mapping for all the rooms present in the server
	rooms map[string]*Room // rood id -> room
}

// Create or join a room.`tea.Program` is the reference server has to the client. This is used
// by the server to send updates to the client. The returned room should be updated on the player
// so that player can send actions to the server
func (s *Server) CreateOrJoinRoom(playerId string, roomId string, program *tea.Program) (*Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := s.rooms[roomId]
	if room != nil {
		room.join(playerId, program)
		return room, nil
	}

	r := NewRoom(roomId)
	if err := r.join(playerId, program); err != nil {
		return nil, err
	}
	s.rooms[roomId] = r

	go r.run()

	return r, nil
}

// Leave the room
func (s *Server) LeaveRoom(playerId, roomId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rooms[roomId].leave(playerId)
}

// Check and cleans all the rooms.
// If room is empty its reference is deleted
func (s *Server) CleanRooms() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, r := range s.rooms {
		if r.clean() {
			delete(s.rooms, id)
			log.Info("cleaning room", "roomId", id)
		}
	}
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}
