package explearn

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	status  int
	err     error
	index   int
	spinner spinner.Model
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
	return checkServer(url)
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
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(m.spinner.Tick())
		return m, cmd
	}

	return m, nil
}

func (m smodel) View() string {

	if m.err != nil {
		return fmt.Sprintf("some thing went wrong %v\n\n", m.err)
	}

	s := fmt.Sprintf("checking ... %s ", url)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	} else {
		s += fmt.Sprintf("  %s", m.spinner.View())
	}

	return s
}

func (m *smodel) resetSpinner(){
	m.spinner.Spinner = spinners[m.index]
}

func TestStatusCli() {
	s := spinner.New()
	s.Spinner = spinner.Line
	if _, err := tea.NewProgram(smodel{spinner: s}).Run(); err != nil {
		panic(err)
	}
}
