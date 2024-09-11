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
type keyMapCreateVault struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Help   key.Binding
	Create key.Binding
	Back   key.Binding
}

func (k keyMapCreateVault) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMapCreateVault) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Back},
		{k.Quit, k.Create, k.Help},
	}
}

var keysCreateVault = keyMapCreateVault{
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
	Create: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "create vault"),
	),
	Back: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "go back"),
	),
}

type CreateVaultModel struct {
	keys      keyMapCreateVault
	help      help.Model
	w, h      int
	mainModel *mainModel

	errorMsg   string
	focusIndex int
	inputs     []textinput.Model
}

const (
	name = iota
	description
	password
	rePassword
)

func InitialCreateVaultModel(mainMdl *mainModel) CreateVaultModel {
	m := CreateVaultModel{
		keys:      keysCreateVault,
		help:      help.New(),
		mainModel: mainMdl,
		inputs:    make([]textinput.Model, 4)}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case name:
			t.Placeholder = "Vault name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.CharLimit = 16
			t.TextStyle = focusedStyle
		case description:
			t.Placeholder = "Description"
			t.CharLimit = 64
		case password:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.CharLimit = 32
			t.EchoCharacter = '•'
		case rePassword:
			t.Placeholder = "Re-enter password"
			t.EchoMode = textinput.EchoPassword
			t.CharLimit = 32
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m CreateVaultModel) Init() tea.Cmd {
	// clearing inputs
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	return nil
}

func (m CreateVaultModel) View() string {
	s := ""
	s += titleStyle.Render(fmt.Sprintf("Create a %s", highlightStyle.Render("vault.")))
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
	s += fmt.Sprintf("Press %s to create vault.\n", highlightStyle.Render("enter"))
	s += errorStyle.Render(fmt.Sprintf("%s\n", m.errorMsg))

	helpView := m.help.View(m.keys)
	s += helpStyle.Render(helpView)
	s = lg.Place(m.w, m.h, lg.Center, lg.Center, s)
	return s
}

func (m CreateVaultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Back):
			m.mainModel.viewState = homeView
			return m.mainModel, tea.WindowSize()

		case key.Matches(msg, m.keys.Create):
			// Create vault
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

	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
		m.help.Width = msg.Width
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *CreateVaultModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

type Vault struct {
	Name                     string   `json:"Name"`
	Description              string   `json:"Description"`
	EncodedEncryptedVaultKey string   `json:"EncodedEncryptedVaultKey"`
	EncodedSalt              string   `json:"EncodedSalt"`
	EncodedNonce             string   `json:"EncodedNonce"`
	Secrets                  []Secret `json:"Secrets"`
}

type Secret struct {
	// array contains encryptedEncoded plaintext, and encodedNonce
	EncodedEncryptedText [2]string // password or any secret data
	EncodedEncryptedName [2]string // name of the secret
}

func (m CreateVaultModel) handleCreate() (tea.Model, tea.Cmd) {
	ok, errMsg := CreateVaultValidation(m.inputs)
	if !ok {
		m.errorMsg = errMsg
		return m, nil
	}

	key, salt, nonce := CreateAndEncryptVaultKey(m.inputs[3].Value())

	newVault := Vault{
		Name:                     m.inputs[name].Value(),
		Description:              m.inputs[description].Value(),
		EncodedEncryptedVaultKey: key,
		EncodedSalt:              salt,
		EncodedNonce:             nonce,
		Secrets:                  make([]Secret, 0),
	}

	newVaultByte, err := json.Marshal(newVault)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = os.WriteFile(fmt.Sprintf("%s%s.json", VAULTSPATH, newVault.Name), newVaultByte, 0644)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	m.mainModel.viewState = homeView
	return m.mainModel, tea.WindowSize()
}

func CreateVaultValidation(inputs []textinput.Model) (bool, string) {
	// Form validation
	// if it doesnt validate return false and set error message
	// Vault name cant contain spaces or special characters it's going to be used as a file name
	errorMsg := ""
	for i := range inputs {
		if len(inputs[i].Value()) == 0 {
			errorMsg = fmt.Sprintf("[%s] option can't be empty!", inputs[i].Placeholder)
			return false, errorMsg
		}
	}

	//  Trim whitespace from the beginning and end of input
	for i := range inputs {
		inputs[i].SetValue(strings.TrimSpace(inputs[i].Value()))
	}

	if _, err := os.Stat(fmt.Sprintf("%s%s.json", VAULTSPATH, inputs[name].Value())); err == nil {
		errorMsg = "Vault with that name already exists!"
	} else if len(inputs[password].Value()) < 8 {
		errorMsg = "Password must be at least 8 characters long!"
	} else if strings.ContainsAny(inputs[name].Value(), "/\\") {
		errorMsg = "Vault name can't contain special characters!"
	} else if strings.ContainsAny(inputs[name].Value(), " ") {
		errorMsg = "Vault name can't contain spaces!"
	} else if strings.ContainsAny(inputs[description].Value(), "/\\") {
		errorMsg = "Description can't contain special characters!"
	} else if strings.ContainsAny(inputs[password].Value(), "/\\") {
		errorMsg = "Password can't contain special characters!"
	} else if strings.ContainsAny(inputs[password].Value(), " ") {
		errorMsg = "Password can't contain spaces!"
	} else if inputs[password].Value() != inputs[rePassword].Value() {
		errorMsg = "Passwords don't match!"
	} else {
		errorMsg = ""
		return true, errorMsg
	}
	return false, errorMsg
}
