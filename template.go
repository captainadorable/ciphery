package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type TemplateModel struct {
	keys      keyMap
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
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}
