/*
 * Copyright (c) 2024 Samuel Dasilva
 *
 * This file is part of Boba Break.
 *
 * Boba Break is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Boba Break is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Boba Break. If not, see <https://www.gnu.org/licenses/>.
 */
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

const (
	// workTime = time.Second * 5
	workTime  = time.Minute * 25
	breakTime = time.Minute * 5
)

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
	return nil
}

func (m BreakModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		// FIXME Has to click twice to start, but seems to be a problem with how the program is called
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		var cmd tea.Cmd
		m.done = true
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = workTime
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
		s = "Go take a break!"
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
		timer: timer.NewWithInterval(workTime, time.Millisecond),
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
	m.keymap.stop.SetEnabled(false)
	// m.keymap.start.SetEnabled(false)
	m.done = false
	return m
}

func (m BreakModel) Start() {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
