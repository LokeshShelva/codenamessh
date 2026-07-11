package backend

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var textInputStyle = lipgloss.NewStyle().PaddingTop(2)

type keyMap struct {
	Up   key.Binding
	Down key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑ arrow", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓ arrow", "move down"),
	),
}

func (u *User) Init() tea.Cmd {
	return textinput.Blink
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := Quit(msg); cmd != nil {
		return u, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if !u.ui.nameInput.Focused() {
				u.ui.nameInput.Focus()
				u.ui.roomInput.Blur()
			}

		case key.Matches(msg, keys.Down):
			if !u.ui.roomInput.Focused() {
				u.ui.roomInput.Focus()
				u.ui.nameInput.Blur()
			}
		}

	}

	var cmd tea.Cmd
	if u.ui.nameInput.Focused() {
		u.ui.nameInput, cmd = u.ui.nameInput.Update(msg)
	} else if u.ui.roomInput.Focused() {
		u.ui.roomInput, cmd = u.ui.roomInput.Update(msg)
	}

	return u, cmd
}

func (u *User) View() tea.View {
	var c *tea.Cursor

	nameInputStr := textInputStyle.Render("Name: " + u.ui.nameInput.View())
	roomInputStr := textInputStyle.Render("Room: " + u.ui.roomInput.View())

	str := lipgloss.JoinVertical(lipgloss.Top, u.headerView(), nameInputStr, roomInputStr, u.ui.help.View(keys), u.footerView())

	v := tea.NewView(str)
	v.Cursor = c
	v.AltScreen = true
	return v
}

func (u *User) headerView() string { return "Welcome to CodeName SSH\n" }
func (u *User) footerView() string { return "\nctrl+c to quit" }

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
	}
}
