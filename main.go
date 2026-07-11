package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	count int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			m.count++

		}

	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Count: %d", m.count)
}

func main() {
	p := tea.NewProgram(model{count: 0})
	_, err := p.Run()

	if err != nil {
		fmt.Printf("Error while running app.\n%v", err)
		os.Exit(1)
	}
}
