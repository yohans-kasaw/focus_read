package explearn

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const url = "https://charm.sh/"

var spinners = []spinner.Spinner{
	spinner.Line,
	spinner.Dot,
	spinner.MiniDot,
	spinner.Jump,
	spinner.Pulse,
	spinner.Points,
	spinner.Globe,
	spinner.Moon,
	spinner.Monkey,
}

type smodel struct {
	textInput   textinput.Model
	style       lipgloss.Style
	color_index int
	status      int
	err         error
	index       int
	spinner     spinner.Model
}

type StatusCodeMsg int

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := http.Client{Timeout: 10 * time.Second}

		// artificial sleep
		time.Sleep(100 * time.Second)

		r, err := c.Get(url)

		if err != nil {
			return err
		}

		return StatusCodeMsg(r.StatusCode)
	}
}

func (m smodel) Init() tea.Cmd {
	return tea.Batch(
		checkServer(url),
		textinput.Blink,
		m.spinner.Tick,
	)

}

func (m smodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatusCodeMsg:
		m.status = int(msg)
		return m, nil
	case error:
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.Type == tea.KeyCtrlUp {
			m.index = (m.index + 1) % len(spinners)
			m.resetSpinner()
			return m, nil
		}

		if msg.Type == tea.KeyUp {
			m.color_index = (m.color_index + 1) % 256
			m.style = m.style.Foreground(lipgloss.Color(fmt.Sprintf("%d", m.color_index)))
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.style = m.style.Width(msg.Width)
		m.style = m.style.Height(msg.Height/2)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(m.spinner.Tick())
		return m, cmd
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m smodel) View() string {

	if m.err != nil {
		return fmt.Sprintf("some thing went wrong %v\n\n", m.err)
	}

	s := fmt.Sprintf("checking ... %s ", url)
	s += fmt.Sprintf(m.textInput.View())

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	} else {
		s += fmt.Sprintf("  %s", m.spinner.View())
	}

	s += fmt.Sprintf("\n\n=n %s", m.textInput.Value())
	s += fmt.Sprintf("\n\n %s %d", m.style.Render("color Testing"), m.color_index)

	return s
}

func (m *smodel) resetSpinner() {
	m.spinner.Spinner = spinners[m.index]
}

func TestStatusCli() {
	s := spinner.New()
	s.Spinner = spinner.Line
	ti := textinput.New()
	ti.Placeholder = "input url"
	ti.Focus()
	ti.Width = 40

	if _, err := tea.NewProgram(smodel{
		spinner:     s,
		textInput:   ti,
		color_index: 1,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center),
	}).Run(); err != nil {
		panic(err)
	}
}
