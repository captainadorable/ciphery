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

// keyMap defines a set of keybindings. To work for help it must satisfy
type keyMapVaults struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Enter key.Binding
	Back  key.Binding
}

func (k keyMapVaults) ShortHelp() []key.Binding {
	return nil
}
func (k keyMapVaults) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Quit, k.Enter, k.Back},
	}
}

var keysVaults = keyMapVaults{
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

type VaultsModel struct {
	keys      keyMapVaults
	help      help.Model
	cursor    int
	w, h      int
	vaults    []Vault
	mainModel *mainModel
}

func InitialVaultsModel(mainMdl *mainModel) VaultsModel {
	m := VaultsModel{keys: keysVaults,
		help:      help.New(),
		mainModel: mainMdl}
	m.help.ShowAll = true
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
	s += titleStyle.Render(fmt.Sprintf("Vaults that %s created", highlightStyle.Render("you've")))
	s += "\n"

	// rendering vaults list
	if len(m.vaults) == 0 {
		s += errorStyle.Render("There is no vaults created. Go back and create one!")
		s += "\n"
	} else {
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

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

// Sending selected vault to the enterVaultView
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
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.vaults)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = homeView
			return m.mainModel.homeView, tea.WindowSize()
		case key.Matches(msg, m.keys.Enter):
			m.mainModel.viewState = enterVaultView
			return m.mainModel.enterVaultView, tea.Batch(tea.WindowSize(), m.mainModel.enterVaultView.Init(), SendVaultCmd(m.vaults[m.cursor]))
		}

	case UpdateVaultsMsg:
		m.vaults = msg.Vaults
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
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
	fmt.Println("qweqwe")
	return vaults
}
