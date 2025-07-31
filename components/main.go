package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	gateOptions list.Model
	keys        *listKeyMap
	results     bool
	track       track
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.gateOptions.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.gateOptions.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.selectTrack):
			x := m.gateOptions.Index()
			m.track = tracks[x]
			m.results = true
			return m, nil
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		}

	}
	var cmd tea.Cmd
	m.gateOptions, cmd = m.gateOptions.Update(msg)
	return m, cmd
}

func (m model) View() string {

	switch {

	case m.results:
		results := fmt.Sprintf("Track:%s\nBy:%s,Gates: %v", m.track.name, m.track.author, m.track.gates)
		return docStyle.Render(results)

	default:
		return docStyle.Render(m.gateOptions.View())
	}

}

func main() {
	m := initTrackList()
	// Initialize our program
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
