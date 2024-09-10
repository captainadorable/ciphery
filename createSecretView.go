package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
type keyMapCreateSecret struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Create key.Binding
}

func (k keyMapCreateSecret) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMapCreateSecret) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Back, k.Help},
		{k.Create},
	}
}

var keysTemplate = keyMapCreateSecret{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
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
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Back: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go back"),
	),
}

type CreateSecretModel struct {
	keys      keyMapCreateSecret
	help      help.Model
	w, h      int
	mainModel *mainModel

	errorMsg   string
	focusIndex int
	inputs     []textinput.Model

	decryptedVaultKey []byte
	vault             Vault
}

const (
	secretName = iota
	secretText
)

func InitialCreateSecretModel(mainmdl *mainModel) CreateSecretModel {
	m := CreateSecretModel{
		keys:      keysTemplate,
		help:      help.New(),
		mainModel: mainmdl,
		inputs:    make([]textinput.Model, 2)}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case secretName:
			t.Placeholder = "Secret name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 16
		case secretText:
			t.Placeholder = "Secret text"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		}
		m.inputs[i] = t
	}
	return m
}

func (m CreateSecretModel) Init() tea.Cmd {
	return nil
}

func (m CreateSecretModel) View() string {
	s := ""
	s += titleStyle.Render("Create new secret")
	s += "\n"

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}
	s += formBorderStyle.Render(b.String())
	s += "\n"
	s += fmt.Sprintf("Press %s to create secret. \n", highlightStyle.Render("enter"))
	s += errorStyle.Render(fmt.Sprintf("%s\n", m.errorMsg))

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

func (m CreateSecretModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = vaultsView
			return m.mainModel.vaultsView, tea.WindowSize()
		case key.Matches(msg, m.keys.Enter):
			return m.handleCreate()
		case key.Matches(msg, m.keys.Up) || key.Matches(msg, m.keys.Down):
			// Cycle indexes
			if key.Matches(msg, m.keys.Up) {
				m.focusIndex--
			} else if key.Matches(msg, m.keys.Down) {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)-1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	case SendDecryptedVaultKeyMsg:
		m.decryptedVaultKey = msg
		return m, nil
	case SendVaultMsg:
		m.vault = msg.VaultSended
		return m, nil
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *CreateSecretModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m CreateSecretModel) handleCreate() (tea.Model, tea.Cmd) {
	for i := range m.inputs {
		if len(m.inputs[i].Value()) == 0 {
			m.errorMsg = fmt.Sprintf("[%s] option can't be empty!", m.inputs[i].Placeholder)
			return m, nil
		}
	}
	// Encrypt the data
	cryptedName, cryptedText := EncryptSecretData(m.inputs[secretName].Value(), m.inputs[secretText].Value(), m.decryptedVaultKey)
	newSecret := Secret{
		EncodedEncryptedName: cryptedName,
		EncodedEncryptedText: cryptedText,
	}

	// Append new secret
	m.vault.Secrets = append(m.vault.Secrets, newSecret)

	// Write the json
	updateVaultByte, err := json.Marshal(m.vault)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.json", VAULTSPATH, m.vault.Name), updateVaultByte, 0644)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Reset the view
	m.mainModel.createSecretView = InitialCreateSecretModel(m.mainModel)

	return m.mainModel.homeView, tea.WindowSize()
}
