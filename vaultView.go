package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type VaultModel struct {
	keys                  keyMap
	help                  help.Model
	w, h                  int
	mainModel             *mainModel
	vault                 Vault
	decryptedVaultSecrets []DecryptedSecret
	decryptedVaultKey     []byte
	errorMsg              string
	confirmationMsg       string
	cursor                int
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
	s += titleStyle.Render(fmt.Sprintf("Vault: %s", highlightStyle.Render(m.vault.Name)))
	s += "\n"

	if len(m.vault.Secrets) == 0 {
		s += errorStyle.Render("You haven't create any secrets yet. Press c to create one.")
		s += "\n"
	} else {
		v := ""
		for i, secret := range m.decryptedVaultSecrets {
			style := listItemStyle
			if m.cursor == i {
				style = listItemHighlightStyle
			}
			v += style.Render(fmt.Sprintf("%s\n%s", secret.SecretName, listItemDescriptionStyle.Render(secret.SecretText)))
			v += "\n"
		}
		s += listStyle.Render(v)
		s += "\n"
	}

	s += errorStyle.Render(fmt.Sprintf("%s\n", m.errorMsg))
	s += confirmationStyle.Render(fmt.Sprintf("%s\n", m.confirmationMsg))

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
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.vault.Secrets)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = vaultsView
			return m.mainModel.vaultsView, tea.Batch(tea.WindowSize(), m.mainModel.vaultsView.Init())
		case key.Matches(msg, m.keys.Create):
			m.mainModel.viewState = createSecretView
			return m.mainModel.createSecretView, tea.Batch(tea.WindowSize(), textinput.Blink, m.mainModel.createSecretView.Init(), SendDecryptedVaultKeyCmd(m.decryptedVaultKey), SendVaultCmd(m.vault))
		case key.Matches(msg, m.keys.Delete):
			if len(m.vault.Secrets) == 0 {
				return m, nil
			}
			return m.handleDelete()
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

func (m VaultModel) handleDelete() (tea.Model, tea.Cmd) {
	m.vault.Secrets = append(m.vault.Secrets[:m.cursor], m.vault.Secrets[m.cursor+1:]...)
	m.decryptedVaultSecrets = append(m.decryptedVaultSecrets[:m.cursor], m.decryptedVaultSecrets[m.cursor+1:]...)
	updateVaultByte, err := json.Marshal(m.vault)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Error deleting secret: %v", err)
		return m, nil
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.json", VAULTSPATH, m.vault.Name), updateVaultByte, 0644)
	if err != nil {
		m.errorMsg = fmt.Sprintf("Error deleting secret: %v", err)
	}
	// reset cursor
	m.cursor = 0
	m.confirmationMsg = "Secret deleted successfully"
	return m, nil
}
