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
	"strings"
	"time"

	"github.com/SamD2021/boba-break/tui/mainmenuui"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/huh"
	"github.com/gen2brain/beeep"
)

const maxWidth = 80
const tickInterval = time.Second / 2

type tickMsg time.Time

func tickCmd(t time.Time) tea.Msg {
	return tickMsg(t)
}

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

const (
	focusColor = "#2EF8BB"
	breakColor = "#FF5F87"
)

var (
	focusTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(focusColor)).MarginRight(1).SetString("Focus Mode")
	breakTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(breakColor)).MarginRight(1).SetString("Break Mode")
	pausedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(breakColor)).MarginRight(1).SetString("Continue?")
	helpStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(2)
	sidebarStyle    = lipgloss.NewStyle().MarginLeft(3).Padding(1, 3).Border(lipgloss.RoundedBorder()).BorderForeground(helpStyle.GetForeground())
)

var baseTimerStyle = lipgloss.NewStyle().Padding(1, 2)

type sessionState int

const (
	Focusing sessionState = iota
	Relaxing
	Paused
)

// const (
// 	workTime = time.Second * 5
// 	// workTime  = time.Minute * 25
// 	breakTime = time.Minute * 5
// )

type BreakModel struct {
	help       help.Model
	keymap     keymap
	Timer      timer.Model
	done       bool
	workTime   time.Duration
	breakTime  time.Duration
	state      sessionState
	count      int8
	scribble   *huh.Form
	scribbling bool
	lg         *lipgloss.Renderer
	styles     *Styles
	width      int
}

type keymap struct {
	start    key.Binding
	stop     key.Binding
	reset    key.Binding
	quit     key.Binding
	back     key.Binding
	scribble key.Binding
}

func (m BreakModel) Init() tea.Cmd {

	return tea.Batch(m.Timer.Init(), m.scribble.Init())
}

func (m BreakModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case timer.TickMsg:
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		m.Timer, cmd = m.Timer.Update(msg)
		m.keymap.stop.SetEnabled(m.Timer.Running())
		m.keymap.start.SetEnabled(!m.Timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		var switchmsg tea.Cmd
		m.done = true
		m.Timer, cmd = m.Timer.Update(msg)
		switch m.state {
		case Focusing:
			icon := ""
			title := "Boba Time"
			message := "Time is up, Enjoy some Boba!"
			err := beeep.Alert(title, message, icon)
			if err != nil {
				fmt.Println("Error sending message: ", err)
			}
		case Relaxing:
			icon := ""
			title := "Get Working!"
			message := "Lets put the cup down and get busy!"
			err := beeep.Alert(title, message, icon)
			if err != nil {
				fmt.Println("Error sending message: ", err)
			}
		}
		m.keymap.stop.SetEnabled(m.Timer.Running())
		m.keymap.start.SetEnabled(!m.Timer.Running())
		switch m.state {
		case Focusing:
			switchmsg =
				func() tea.Msg {
					return SwitchBreakMsg{}
				}

		case Relaxing:
			m.count++
			switchmsg = func() tea.Msg {
				return SwitchWorkMsg{}
			}
		}
		return m, tea.Batch(cmd, switchmsg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			var cmds []tea.Cmd
			switch m.state {
			case Focusing:
				m.Timer.Timeout = m.workTime
			case Relaxing:
				m.Timer.Timeout = m.breakTime
			}
			m.done = false
			cmd = m.Timer.Stop()
			cmds = append(cmds, cmd)
			m.Timer, cmd = m.Timer.Update(timer.TickMsg{})
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.Timer.Toggle()
		case key.Matches(msg, m.keymap.back):
			return m,
				func() tea.Msg {
					return BackMsg{}
				}
		case key.Matches(msg, m.keymap.scribble):
			cmd = m.Timer.Toggle()
			return m, tea.Batch(cmd, func() tea.Msg {
				return ScribblingMsg{}
			})

		}
	case mainmenuui.SelectedBreakManagerMsg:
		// m.Timer, cmd = m.Timer.Update(timer.TickMsg{})
		return m, m.Timer.Init()
	case SwitchWorkMsg:
		m.state = Focusing
		m.Timer.Timeout = m.workTime
		m.done = false
		cmd = m.Timer.Stop()
		m.keymap.stop.SetEnabled(m.Timer.Running())
		m.keymap.start.SetEnabled(!m.Timer.Running())
		return m, cmd
	case SwitchBreakMsg:
		m.state = Relaxing
		m.Timer.Timeout = m.breakTime
		m.done = false
		cmd = m.Timer.Stop()
		m.keymap.stop.SetEnabled(m.Timer.Running())
		m.keymap.start.SetEnabled(!m.Timer.Running())
		return m, cmd
	case ScribblingMsg:
		m.scribbling = true
		m.scribble.Run()
		cmds = append(cmds,
			func() tea.Msg {
				return timer.TickMsg{}
			})
	}
	// Process the form
	form, cmd := m.scribble.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.scribble = f
		cmds = append(cmds, cmd)
	}
	if m.scribble.State == huh.StateCompleted {
		// Quit when the form is done.
		m.scribbling = false
		m.scribble = huh.NewForm(huh.NewGroup(huh.NewText().Title("Scribble").Placeholder("Current thoughts")))
		cmd = m.Timer.Toggle()
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m BreakModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
		m.keymap.back,
		m.keymap.scribble,
	})
}

func (m BreakModel) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	// s := m.timer.View()
	var timer string
	var scribble string

	v := strings.TrimSuffix(m.TimerView(), "\n\n")
	timer = m.lg.NewStyle().Margin(1, 0).Render(v)
	styles := m.styles
	// // var sb strings.Builder
	var body string
	var footer string
	header := m.appBoundaryView("Boba Break")
	if m.scribbling {
		sv := strings.TrimSuffix(m.scribble.View(), "\n\n")
		scribble = m.lg.NewStyle().Margin(1, 1).Render(sv)
		body = lipgloss.JoinVertical(lipgloss.Top, timer, scribble)
		footer = m.appBoundaryView(m.scribble.Help().ShortHelpView(m.scribble.KeyBinds()))
	} else {
		body = lipgloss.JoinVertical(lipgloss.Top, timer)
		footer = m.appBoundaryView(m.helpView())
	}
	// if len(errors) > 0 {
	// 	footer = m.appErrorBoundaryView("")
	// }
	return styles.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m BreakModel) TimerView() string {
	styles := m.styles
	ms := m.Timer.Timeout.Milliseconds()
	hours := ms / (3.6e+6)
	ms -= hours
	minutes := ms / 60000
	ms -= minutes
	seconds := ms / 1000
	// remainingMs := ms % 1000 // Capture the remainder milliseconds
	var s string
	switch m.state {
	case Focusing:
		s = fmt.Sprintf("Work Session: %v\n", m.count)
	case Relaxing:
		s = fmt.Sprintf("Break Session: %v\n", m.count)
	}
	s += styles.Highlight.Render(fmt.Sprintf("%v h %v m %v s", hours%60, minutes%60, seconds%60))

	// if m.Timer.Timedout() && m.state == Focusing {
	// 	s = "Go take a break!"
	// 	s += m.helpView()
	// } else if m.Timer.Timedout() && m.state == Relaxing {
	// 	s = "Get back to work!"
	// 	s += m.helpView()
	// }
	// s += "\n"
	return styles.Status.Copy().Margin(0, 1).Padding(1, 2).Width(48).Render(s) + "\n\n"
}
func (m BreakModel) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}
func (m BreakModel) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func InitialModel(workDuration time.Duration, breakDuration time.Duration) BreakModel {
	m := BreakModel{
		width: maxWidth,
		help:  help.New(),
		keymap: keymap{
			start:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start")),
			stop:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "stop")),
			reset:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset")),
			quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
			back:     key.NewBinding(key.WithKeys("backspace"), key.WithHelp("backspace", "back")),
			scribble: key.NewBinding(key.WithKeys("n"), key.WithHelp("scribble", "n")),
		},
		Timer:      timer.NewWithInterval(workDuration, tickInterval),
		done:       false,
		workTime:   workDuration,
		breakTime:  breakDuration,
		state:      Focusing,
		count:      1,
		scribble:   huh.NewForm(huh.NewGroup(huh.NewText().Title("Scribble").Placeholder("Current thoughts"))),
		lg:         lipgloss.DefaultRenderer(),
		styles:     NewStyles(lipgloss.DefaultRenderer()),
		scribbling: false,
	}
	m.keymap.stop.SetEnabled(true)
	m.keymap.start.SetEnabled(false)
	// m.keymap.scribble.SetEnabled(false)
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
