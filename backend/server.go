package backend

type Server struct {
	rooms []*Room
}

func CreateServer() *Server {
	return &Server{
		rooms: make([]*Room, 0),
	}
}
