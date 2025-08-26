package cli

import (
	"fmt"
	"focusRead/epub"
	"os"
	"strconv"
	"strings"

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
	s += fmt.Sprintln("\n", m.index, "/", len(m.texts))
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
			PaddingLeft(msg.Width / 4).
			PaddingRight(msg.Width / 4)
		return m, nil
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func Run(texts []epub.Text) {
	indexPath := "./cache/index_cache.bin"

	savedIndex, err := readIndexFromFile(indexPath)
	if err != nil {
		fmt.Printf("Could not read saved index, starting from beginning: %v", err)
		savedIndex = 0
	}

	if savedIndex >= len(texts) || savedIndex < 0 {
		savedIndex = 0
	}

	p := tea.NewProgram(model{
		texts: texts,
		index: savedIndex,
		style: lipgloss.NewStyle().
			Background(lipgloss.Color("#121212")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center),
	})

	finalModel, err := p.Run()
	if err != nil {
		panic(err)
	}

	if m, ok := finalModel.(model); ok {
		err := writeIndexToFile(indexPath, m.index)
		if err != nil {
			fmt.Printf("Error saving final index: %v", err)
		} else {
			fmt.Println("Progress saved.")
		}
	}
}

func readIndexFromFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read index file: %w", err)
	}

	index, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		fmt.Printf("Warning: could not parse index from file, defaulting to 0: %v", err)
		return 0, nil
	}

	return index, nil
}

func writeIndexToFile(path string, index int) error {
	data := []byte(strconv.Itoa(index))
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}
	return nil
}
