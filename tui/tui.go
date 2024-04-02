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
package tui

import (
	"fmt"
	"os"

	"github.com/SamD2021/boba-break/tui/breakmanagerui"
	"github.com/SamD2021/boba-break/tui/mainmenuui"
	"github.com/SamD2021/boba-break/tui/noteui"
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
	notesView
)

type MainModel struct {
	mainMenu     tea.Model
	breakManager tea.Model
	notes        tea.Model
	state        sessionState
}

// View implements tea.Model.
func (m MainModel) View() string {
	switch m.state {
	case mainMenuView:
		return m.mainMenu.View()
	case breakManagerView:
		return m.breakManager.View()
	case notesView:
		return m.notes.View()
	default:
		panic("Not implemented yet")
	}
}

func initialModel() MainModel {
	return MainModel{
		state:        mainMenuView,
		mainMenu:     mainmenuui.NewModel(),
		breakManager: breakmanagerui.InitialModel(),
		notes:        noteui.InitialModel(),
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
	case mainmenuui.SelectedNoteMsg:
		m.state = notesView
	case noteui.GoBackMsg:
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
	case notesView:
		newModel, newCmd := m.notes.Update(msg)
		model, ok := newModel.(noteui.NotesModel)
		if !ok {
			panic("Couldn't assert newModel is of type BreakModel")
		}
		m.notes = model
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
	// breakmanagerui.Start()
}
