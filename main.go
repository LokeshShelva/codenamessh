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

var room = backend.NewRoom()

func bubbleTeaHandler() wish.Middleware {
	teaHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		playerId := s.Context().SessionID()
		player := backend.NewPlayer(playerId, room)

		p := tea.NewProgram(player, bubbletea.MakeOptions(s)...)

		err := room.Join(playerId, p)
		if err != nil {
			log.Error("failed to join room", "error", err)
			return nil
		}

		log.Info("new player connection", "playerId", playerId)
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

	// Global server tick of 500 MS
	// TODO: this needs to be made to tick independantly for each room
	go func() {
		for {
			<-time.After(500 * time.Millisecond)
			room.Run()
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
