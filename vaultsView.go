package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type VaultsModel struct {
	keys            keyMap
	help            help.Model
	cursor          int
	w, h            int
	vaults          []Vault
	errorMsg        string
	confirmationMsg string
	mainModel       *mainModel
}

func InitialVaultsModel(mainMdl *mainModel) VaultsModel {
	m := VaultsModel{keys: keysVaults,
		help:      help.New(),
		mainModel: mainMdl}
	m.vaults = m.GetVaults()

	return m
}

type UpdateVaultsMsg struct {
	Vaults []Vault
}

func UpdateVaultsCmd(vaults []Vault) tea.Cmd {
	return func() tea.Msg {
		return UpdateVaultsMsg{Vaults: vaults}
	}
}

func (m VaultsModel) Init() tea.Cmd {
	return UpdateVaultsCmd(m.GetVaults())
}

func (m VaultsModel) View() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("Vaults that you've %s", highlightStyle.Render("created.")))
	s += "\n"

	// rendering vaults list
	if len(m.vaults) > 0 {
		v := ""
		for i, vault := range m.vaults {
			style := listItemStyle
			if m.cursor == i {
				style = listItemHighlightStyle
			}
			v += style.Render(fmt.Sprintf("%s\n%s", vault.Name, listItemDescriptionStyle.Render(vault.Description)))
			v += "\n"
		}
		s += listStyle.Render(v)
		s += "\n"
	}
	s += errorStyle.Render(m.errorMsg)
	s += "\n"
	s += confirmationStyle.Render(m.confirmationMsg)
	s += "\n"

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

// Sending selected vault to the enterVaultView.
type SendVaultMsg struct {
	VaultSended Vault
}

func SendVaultCmd(vault Vault) tea.Cmd {
	return func() tea.Msg {
		return SendVaultMsg{VaultSended: vault}
	}
}

func (m VaultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.vaults)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Delete):
			return m.handleDelete()
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = homeView
			return m.mainModel.homeView, tea.WindowSize()
		case key.Matches(msg, m.keys.Enter):
			if len(m.vaults) == 0 {
				m.errorMsg = "There is no vaults created. Go back and create one!"
				return m, nil
			}
			m.mainModel.viewState = enterVaultView
			return m.mainModel.enterVaultView, tea.Batch(tea.WindowSize(), m.mainModel.enterVaultView.Init(), SendVaultCmd(m.vaults[m.cursor]))
		}

	case UpdateVaultsMsg:
		m.vaults = msg.Vaults
		if len(m.vaults) == 0 {
			m.errorMsg = "There is no vaults created. Go back and create one!"
		} else {
			m.errorMsg = ""
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}
	return m, nil
}

func (m VaultsModel) GetVaults() []Vault {
	vaults := []Vault{}

	files, err := os.ReadDir(VAULTSPATH)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			// Skip directories
			continue
		}
		// Reading all the .jsons and pushing them to a list
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(VAULTSPATH, file.Name())
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatal(err)
			}

			var vault Vault
			if err := json.Unmarshal(fileData, &vault); err != nil {
				log.Fatal(err)
			}
			vaults = append(vaults, vault)
		}
	}
	return vaults
}
func (m VaultsModel) handleDelete() (tea.Model, tea.Cmd) {
	err := os.Remove(fmt.Sprintf("%s%s.json", VAULTSPATH, m.vaults[m.cursor].Name))
	if err != nil {
		m.errorMsg = fmt.Sprintf("Error deleting vault: %v", err)
	} else {
		m.errorMsg = ""
		m.confirmationMsg = fmt.Sprintf("Vault %s deleted successfully", m.vaults[m.cursor].Name)
	}

	m.vaults = append(m.vaults[:m.cursor], m.vaults[m.cursor+1:]...)

	if len(m.vaults) == 0 {
		m.errorMsg = "There is no vaults created. Go back and create one!"
	}
	return m, nil
}
