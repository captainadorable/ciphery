package main

import "github.com/charmbracelet/bubbles/key"

// Key bindings for the home view.
var keysHome = HomeKeyMap()

func HomeKeyMap() keyMap {
	keys := newKeyMap()

	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down},
		{keys.Quit, keys.Enter, keys.Help},
	}
	return keys
}

// Key bindings for the vault view.
var keysVault = VaultKeyMap()

func VaultKeyMap() keyMap {
	keys := newKeyMap()
	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down, keys.Back},
		{keys.Quit, keys.Enter, keys.Help},
		{keys.Create, keys.Delete},
	}
	return keys
}

// Key bindings for the create secret view.
var keysCreateSecret = CreateSecretKeyMap()

func CreateSecretKeyMap() keyMap {
	keys := newKeyMap()
	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down, keys.Back},
		{keys.Quit, keys.Enter, keys.Help},
	}
	return keys
}

// Key bindings for the create vault view.
var keysCreateVault = CreateVaultKeyMap()

func CreateVaultKeyMap() keyMap {
	keys := newKeyMap()
	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down, keys.Back},
		{keys.Quit, keys.Enter, keys.Help},
	}
	return keys
}

// Key bindings for the enter vault view.
var keysEnterVault = EnterVaultKeyMap()

func EnterVaultKeyMap() keyMap {
	keys := newKeyMap()
	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down, keys.Back},
		{keys.Quit, keys.Enter, keys.Help},
	}
	return keys
}

// Key bindings for the vaults view.
var keysVaults = VaultsKeyMap()

func VaultsKeyMap() keyMap {
	keys := newKeyMap()
	keys.Full = [][]key.Binding{
		{keys.Up, keys.Down, keys.Back},
		{keys.Quit, keys.Enter, keys.Help},
		{keys.Delete},
	}
	return keys
}

// Define a generic key map struct
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Create key.Binding
	Delete key.Binding
	Full   [][]key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return k.Full
}

// General bindings
func newKeyMap() keyMap {
	return keyMap{
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
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete secret"),
		),
		Create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create secret"),
		),
	}
}
