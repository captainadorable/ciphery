package main

import lg "github.com/charmbracelet/lipgloss"

var (
	titleStyle          = lg.NewStyle().Bold(true).Border(lg.DoubleBorder()).MarginBottom(1).Padding(1, 1)
	highlightStyle      = lg.NewStyle().Bold(true).Foreground(lg.Color("#88D7E5"))
	choicesStyle        = lg.NewStyle()
	choicesFocusedStyle = lg.NewStyle().Bold(true).Foreground(lg.Color("#F15757"))
	helpStyle           = lg.NewStyle().AlignHorizontal(lg.Center).MarginTop(1)

	errorStyle      = lg.NewStyle().Bold(true).Foreground(lg.Color("#E96379"))
	focusedStyle    = lg.NewStyle().Foreground(lg.Color("205"))
	cursorStyle     = focusedStyle
	noStyle         = lg.NewStyle()
	formBorderStyle = lg.NewStyle().Border(lg.RoundedBorder()).Padding(1)

	listStyle                = lg.NewStyle().Width(20)
	listItemStyle            = lg.NewStyle().Border(lg.NormalBorder(), false, false, false, true).MarginTop(1).PaddingLeft(1)
	listItemHighlightStyle   = listItemStyle.BorderForeground(lg.Color("205"))
	listItemDescriptionStyle = lg.NewStyle().Italic(true).Foreground(lg.Color("#6C6C6C"))
)
