package noteui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type keymap struct {
	back key.Binding
}

type NotesModel struct {
	textarea textarea.Model
	err      error
	help     help.Model
	keymap   keymap
}

func (m NotesModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.back,
	})
}

func InitialModel() NotesModel {
	ti := textarea.New()
	ti.Placeholder = "Brainstorm..."
	ti.Focus()

	return NotesModel{
		textarea: ti,
		err:      nil,
		keymap: keymap{
			back: key.NewBinding(
				key.WithKeys("ctrl", "shift", "<"),
				key.WithHelp("ctrl < ", "back"),
			),
		},
		help: help.New(),
	}
}

func (m NotesModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m NotesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
		switch {
		case key.Matches(msg, m.keymap.back):
			return m,
				func() tea.Msg {
					return GoBackMsg{}
				}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m NotesModel) View() string {
	return fmt.Sprintf(
		"Write down thoughts.\n\n%s\n\n%s",
		m.textarea.View(),
		m.helpView(),
	) + "\n\n"
}
