package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
type keyMapTemplate struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Create key.Binding
}

func (k keyMapTemplate) ShortHelp() []key.Binding {
	return nil
}
func (k keyMapTemplate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Back},
		{k.Create},
	}
}

var keysCreateSecret = keyMapTemplate{
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
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "enter"),
	),
	Back: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go back"),
	),
}

type TemplateModel struct {
	keys      keyMapTemplate
	help      help.Model
	w, h      int
	mainModel *mainModel
}

func InitialTemplateModel(mainmdl *mainModel) TemplateModel {
	m := TemplateModel{
		keys:      keysCreateSecret,
		help:      help.New(),
		mainModel: mainmdl,
	}
	m.help.ShowAll = true
	return m
}

func (m TemplateModel) Init() tea.Cmd {
	return nil
}

func (m TemplateModel) View() string {
	s := ""
	s += titleStyle.Render("Create secret")
	s += "\n"

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

func (m TemplateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = vaultsView
			return m.mainModel.vaultsView, tea.WindowSize()
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}
