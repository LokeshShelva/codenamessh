package backend

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var textInputStyle = lipgloss.NewStyle().PaddingTop(2).PaddingBottom(2)

func (u *User) Init() tea.Cmd {
	return textinput.Blink
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := Quit(msg); cmd != nil {
		return u, cmd
	}

	var cmd tea.Cmd
	u.textinput, cmd = u.textinput.Update(msg)

	return u, cmd
}

func (u *User) View() tea.View {
	var c *tea.Cursor

	// if !u.textinput.VirtualCursor() {
	// 	c = u.textinput.Cursor()
	// 	c.Y += lipgloss.Height(u.headerView())
	// }

	textInputStr := textInputStyle.Render("Name: " + u.textinput.View())

	str := lipgloss.JoinVertical(lipgloss.Top, u.headerView(), textInputStr, u.footerView())

	v := tea.NewView(str)
	v.Cursor = c
	v.AltScreen = true
	return v
}

func (u *User) headerView() string { return "Welcome to CodeName SSH\n" }
func (u *User) footerView() string { return "\nctrl+c to quit" }
