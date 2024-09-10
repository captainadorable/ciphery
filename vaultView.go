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
type keyMapVault struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Create key.Binding
}

func (k keyMapVault) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMapVault) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Back, k.Help},
		{k.Create},
	}
}

var keysVault = keyMapVault{
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
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create new secret"),
	),
}

type VaultModel struct {
	keys                  keyMapVault
	help                  help.Model
	w, h                  int
	mainModel             *mainModel
	vault                 Vault
	decryptedVaultSecrets []DecryptedSecret
	decryptedVaultKey     []byte
}

func InitialVaultModel(mainmdl *mainModel) VaultModel {
	m := VaultModel{
		keys:      keysVault,
		help:      help.New(),
		mainModel: mainmdl,
	}
	return m
}

func (m VaultModel) Init() tea.Cmd {
	return nil
}

func (m VaultModel) View() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("Viewing %s.", highlightStyle.Render(m.vault.Name)))
	s += "\n"

	if len(m.vault.Secrets) == 0 {
		s += errorStyle.Render("You haven't create any secrets yet. Press c to create one.")
		s += "\n"
	} else {
		v := ""
		for _, secret := range m.decryptedVaultSecrets {
			style := listItemStyle
			v += style.Render(fmt.Sprintf("%s\n%s", secret.SecretName, listItemDescriptionStyle.Render(secret.SecretText)))
			v += "\n"
		}
		s += listStyle.Render(v)
		s += "\n"
	}

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

func (m VaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = vaultsView
			return m.mainModel.vaultsView, tea.Batch(tea.WindowSize(), m.mainModel.vaultsView.Init())
		case key.Matches(msg, m.keys.Create):
			m.mainModel.viewState = createSecretView
			return m.mainModel.createSecretView, tea.Batch(tea.WindowSize(), textinput.Blink, m.mainModel.createSecretView.Init(), SendDecryptedVaultKeyCmd(m.decryptedVaultKey), SendVaultCmd(m.vault))
		}
	case SendVaultMsg:
		m.vault = msg.VaultSended
		m.decryptedVaultSecrets = make([]DecryptedSecret, len(m.vault.Secrets))
		return m, nil
	case SendDecryptedVaultKeyMsg:
		m.decryptedVaultKey = msg
		m.decryptVaultSecrets()
		return m, nil
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}

type DecryptedSecret struct {
	SecretName string
	SecretText string
}

func (m VaultModel) decryptVaultSecrets() {
	for i := range m.vault.Secrets {
		m.decryptedVaultSecrets[i].SecretName, m.decryptedVaultSecrets[i].SecretText = DecryptSecretData(m.vault.Secrets[i].EncodedEncryptedName, m.vault.Secrets[i].EncodedEncryptedText, m.decryptedVaultKey)
	}
}
