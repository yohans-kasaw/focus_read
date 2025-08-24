package explearn

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type smodel struct {
	status int
	err    error
}

type StatusCodeMsg int

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := http.Client{Timeout: 10 * time.Second}

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
	}

	return s
}

func TestStatusCli() {
	if _, err := tea.NewProgram(smodel{}).Run(); err != nil {
		panic(err)
	}
}
