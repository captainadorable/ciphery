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
type keyMapEnterVault struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Help  key.Binding
	Enter key.Binding
	Back  key.Binding
}

func (k keyMapEnterVault) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMapEnterVault) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Back, k.Help},
	}
}

var keysEnterVault = keyMapEnterVault{
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
	Back: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go back"),
	),
}

type EnterVaultModel struct {
	keys      keyMapEnterVault
	help      help.Model
	w, h      int
	mainModel *mainModel
	vault     Vault
	textInput textinput.Model
	errorMsg  string
}

func InitialEnterVaultModel(mainmdl *mainModel) EnterVaultModel {
	m := EnterVaultModel{
		keys:      keysEnterVault,
		help:      help.New(),
		mainModel: mainmdl,
		errorMsg:  "",
	}

	ti := textinput.New()
	ti.Placeholder = "Master password"
	ti.Focus()
	ti.CharLimit = 16
	ti.Width = 20
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	m.textInput = ti
	return m
}

func (m EnterVaultModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m EnterVaultModel) View() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("Entering vault %s", highlightStyle.Render(m.vault.Name)))
	s += "\n"

	s += focusedStyle.Render(m.textInput.View())
	s += "\n"

	s += errorStyle.Render(m.errorMsg)
	s += "\n"

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

func (m EnterVaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Enter):
			return m.handleEnterVault()
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = vaultsView
			return m.mainModel.vaultsView, tea.WindowSize()
		}
	case SendVaultMsg:
		m.vault = msg.VaultSended
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

type SendDecryptedVaultKeyMsg []byte

func SendDecryptedVaultKeyCmd(decryptedVaultKey []byte) tea.Cmd {
	return func() tea.Msg {
		return SendDecryptedVaultKeyMsg(decryptedVaultKey)
	}
}

func (m EnterVaultModel) handleEnterVault() (tea.Model, tea.Cmd) {
	decryptedVaultKey, auth := DecryptVaultKeyFromPassword(m.textInput.Value(), m.vault.EncodedSalt, m.vault.EncodedEncryptedVaultKey, m.vault.EncodedNonce)
	if !auth {
		m.errorMsg = "Wrong master password!"
		return m, nil
	}

	m.mainModel.viewState = vaultView
	return m.mainModel.vaultView, tea.Batch(tea.WindowSize(), m.mainModel.vaultView.Init(), SendVaultCmd(m.vault), SendDecryptedVaultKeyCmd(decryptedVaultKey))
}
