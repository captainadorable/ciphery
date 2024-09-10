package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	homeView = iota
	vaultsView
	createVaultView
	enterVaultView
	vaultView
	createSecretView
)

const VAULTSPATH = "vaults/"

var Program *tea.Program

func main() {
	Program = tea.NewProgram(initialMainModel())
	if _, err := Program.Run(); err != nil {
		fmt.Printf("There is been an error: %v", err)
		os.Exit(1)
	}
}

type mainModel struct {
	viewState        int
	homeView         tea.Model
	vaultsView       tea.Model
	createVaultView  tea.Model
	enterVaultView   tea.Model
	vaultView        tea.Model
	createSecretView tea.Model
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.WindowSize())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.viewState {
	default:
		model, cmd := m.homeView.Update(msg)
		return model, cmd
	case vaultsView:
		model, cmd := m.vaultsView.Update(msg)
		return model, cmd
	case createVaultView:
		model, cmd := m.createVaultView.Update(msg)
		return model, cmd
	case enterVaultView:
		model, cmd := m.enterVaultView.Update(msg)
		return model, cmd
	case createSecretView:
		model, cmd := m.createSecretView.Update(msg)
		return model, cmd

	}
}

func (m mainModel) View() string {
	switch m.viewState {
	default:
		return m.homeView.View()
	case vaultsView:
		return m.vaultsView.View()
	case createVaultView:
		return m.createVaultView.View()
	case enterVaultView:
		return m.enterVaultView.View()
	case createSecretView:
		return m.createSecretView.View()
	}
}

func initialMainModel() mainModel {
	var m mainModel
	m = mainModel{
		viewState:        homeView,
		homeView:         InitialHomeModel(&m),
		vaultsView:       InitialVaultsModel(&m),
		createVaultView:  InitialCreateVaultModel(&m),
		enterVaultView:   InitialEnterVaultModel(&m),
		vaultView:        InitialVaultModel(&m),
		createSecretView: InitialCreateSecretModel(&m)}

	return m
}
