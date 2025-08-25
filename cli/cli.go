package cli

import (
	"fmt"
	"focusRead/epub"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	texts []epub.Text
	index int
	style lipgloss.Style
}

func (m model) View() string {
	s := fmt.Sprintln(m.style.Render(m.texts[m.index].Text))
	s += fmt.Sprintln("\n", m.index)
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right", " ":
			if m.index < len(m.texts) {
				m.index++
			}
			return m, nil
		case "left":
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.style = m.style.
			Width(msg.Width).
			Height(msg.Height).
			PaddingLeft(msg.Width/4).
			PaddingRight(msg.Width/4)
		return m, nil
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func Run(texts []epub.Text) {
	p := tea.NewProgram(model{
		texts: texts,
		style: lipgloss.NewStyle().
			Background(lipgloss.Color("#121212")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center),
	})

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
