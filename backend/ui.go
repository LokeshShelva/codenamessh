package backend

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
)

type TimeMsg time.Time

func (m UIModel) Init() tea.Cmd {
	return nil
}

func Quit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return tea.Quit
		}

	}
	return nil
}

func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TimeMsg:
		m.Time = time.Time(msg)

	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			m.Count++

		}

	}
	return m, nil
}

func (m UIModel) View() tea.View {
	s := "Hello... Welcome\n"
	s += "Term: %s\n"
	s += "Width: %d, Height: %d\n"
	s += "Time: " + m.Time.Format(time.RFC1123) + "\n\n"
	s += "Press ctrl+c or q to quit"

	content := fmt.Sprintf(s, m.Term, m.Width, m.Height)
	v := tea.NewView(content)
	v.AltScreen = true

	return v
}
