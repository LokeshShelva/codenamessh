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

func bubbleTeaHandler() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)

		go func() {
			for {
				<-time.After(1 * time.Second)
				p.Send(backend.TimeMsg(time.Now()))
			}

		}()
		return p
	}

	teaHandler := func(s ssh.Session) *tea.Program {
		_, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		user := backend.CreateUser(s.User())
		log.Info("user connected", "id", user.Id, "name", user.Name)

		return newProg(user, bubbletea.MakeOptions(s)...)
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

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done

	log.Info("stopping ssh server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	if err = s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}
