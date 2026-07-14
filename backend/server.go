package backend

import (
	"sync"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
)

type Server struct {
	mu    sync.Mutex
	rooms map[string]*Room // rood id -> room
}

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

func (s *Server) LeaveRoom(playerId, roomId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.rooms[roomId].leave(playerId)
}

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
