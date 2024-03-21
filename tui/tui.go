package tui

import (
	"fmt"
	"os"

	"github.com/SamD2021/boba-break/tui/breakmanagerui"
	"github.com/SamD2021/boba-break/tui/mainmenuui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingTop(2).
	PaddingLeft(4).
	Width(22)

type sessionState int

const (
	mainMenuView sessionState = iota
	breakManagerView
	NotesView
)

type MainModel struct {
	mainMenu     tea.Model
	breakManager tea.Model
	state        sessionState
}

// View implements tea.Model.
func (m MainModel) View() string {
	switch m.state {
	case mainMenuView:
		return m.mainMenu.View()
	case breakManagerView:
		return m.breakManager.View()
	default:
		panic("Not implemented yet")
	}
}

func initialModel() MainModel {
	return MainModel{
		state:        mainMenuView,
		mainMenu:     mainmenuui.NewModel(),
		breakManager: breakmanagerui.InitialModel(),
	}
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg.(type) {
	case mainmenuui.SelectedBreakManagerMsg:
		m.state = breakManagerView
	case breakmanagerui.BackMsg:
		m.state = mainMenuView
	}
	switch m.state {
	case mainMenuView:
		model, newCmd := m.mainMenu.Update(msg)
		newModel, ok := model.(mainmenuui.Model)
		if !ok {
			panic("Couldn't assert return value needed")
		}
		m.mainMenu = newModel
		cmd = newCmd
	case breakManagerView:
		newModel, newCmd := m.breakManager.Update(msg)
		model, ok := newModel.(breakmanagerui.BreakModel)
		if !ok {
			panic("Couldn't assert newModel is of type BreakModel")
		}
		m.breakManager = model
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func Start() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
