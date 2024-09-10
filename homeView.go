package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
type keyMapHome struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Help  key.Binding
	Enter key.Binding
}

func (k keyMapHome) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMapHome) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Help},
	}
}

var keysHome = keyMapHome{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit program"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "enter"),
	),
}

type HomeModel struct {
	keys      keyMapHome
	help      help.Model
	choices   []string
	cursor    int
	w, h      int
	mainModel *mainModel
}

func InitialHomeModel(mainmdl *mainModel) HomeModel {
	m := HomeModel{
		choices:   []string{"Enter to an existing vault.", "Create a new vault."},
		keys:      keysHome,
		help:      help.New(),
		mainModel: mainmdl,
	}
	return m
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) View() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("Welcome to %s", highlightStyle.Render("Ciphery!")))
	s += "\n"

	for i, choice := range m.choices {
		if m.cursor == i {
			s += choicesFocusedStyle.Render(fmt.Sprintf("> %s", choice))

		} else {
			s += choicesStyle.Render(fmt.Sprintf("  %s", choice))
		}
		s += "\n"
	}
	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)

	return s
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 {
				m.mainModel.viewState = vaultsView
				return m.mainModel.vaultsView, tea.Batch(tea.WindowSize(), m.mainModel.vaultsView.Init())
			} else if m.cursor == 1 {
				m.mainModel.viewState = createVaultView
				return m.mainModel.createVaultView, tea.Batch(tea.WindowSize(), m.mainModel.createVaultView.Init(), textinput.Blink)
			}
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}
