package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LokeshShelva/codenamessh/backend"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
)

const (
	host = "localhost"
	port = "6767"
)

var server = backend.NewServer()

func bubbleTeaHandler() wish.Middleware {
	teaHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		playerId := s.Context().SessionID()
		player := backend.NewPlayer(playerId)

		p := tea.NewProgram(player, bubbletea.MakeOptions(s)...)

		roomID := "124"
		room, err := server.CreateOrJoinRoom(playerId, roomID, p)
		if err != nil {
			log.Error("Failed to join room", "roomId", roomID, "playerId", playerId, "error", err)
			return nil
		}
		player.SetRoom(room)

		go func() {
			<-s.Context().Done()
			server.LeaveRoom(playerId, room.Id)
		}()

		log.Info("new player connection", "playerId", playerId, "roomId", room.Id)
		return p
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler)
}

func main() {
	s, err := wish.NewServer(
		ssh.AllocatePty(),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbleTeaHandler(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("could not start ssh server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("starting ssh server", "host", host, "port", port)

	// Garbase collect rooms with no users
	go func() {
		for {
			<-time.After(10 * time.Second)
			server.CleanRooms()
		}

	}()

	// Start SSH server
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server :(", "error", err)
			done <- nil
		}
	}()

	// Wait for Interrupt signal
	<-done

	log.Info("stopping ssh server. See you again soon :)")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	if err = s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}
