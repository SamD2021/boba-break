package breakmanagerui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// const timeout = time.Second * 5
const timeout = time.Hour * 1

type BreakModel struct {
	help   help.Model
	keymap keymap
	timer  timer.Model
	done   bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
	back  key.Binding
}

func (m BreakModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m BreakModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		// var cmd tea.Cmd
		m.done = true
		// m.timer, cmd = m.timer.Update(msg)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = timeout
			m.done = false
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		case key.Matches(msg, m.keymap.back):
			return m,
				func() tea.Msg {
					return BackMsg{}
				}
		}
	}

	return m, nil
}

func (m BreakModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
		m.keymap.back,
	})
}

func (m BreakModel) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	// s := m.timer.View()
	ms := m.timer.Timeout.Milliseconds()
	hours := ms / (3.6e+6)
	ms -= hours
	minutes := ms / 60000
	ms -= minutes
	seconds := ms / 1000
	s := fmt.Sprintf("%v h %v m %v s", hours%60, minutes%60, seconds%60)

	if m.timer.Timedout() {
		s = "All done!"
		s += m.helpView()
	}
	s += "\n"
	if !m.done {
		// s = "Exiting in " + s
		s += m.helpView()
	}
	return s
}

func InitialModel() BreakModel {
	m := BreakModel{
		timer: timer.NewWithInterval(timeout, time.Millisecond),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			back: key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "back"),
			),
		},
		help: help.New(),
	}
	m.keymap.start.SetEnabled(false)
	return m
}

func Start() {
	m := InitialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
