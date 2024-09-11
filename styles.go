package main

import lg "github.com/charmbracelet/lipgloss"

const (
	primaryBg          = "#3D5A80"
	primaryFg          = "#FFFFFF"
	secondaryFg        = "#6C6C6C"
	primaryHighlight   = "#B5F5F8"
	secondaryHighlight = "#7FBDD7"
	err                = "#FE5F55"
	confirm            = "#70C1B3"
)

var (
	titleStyle          = lg.NewStyle().Bold(true).MarginBottom(1).Background(lg.Color(primaryBg)).Foreground(lg.Color(primaryFg)).Padding(1, 2)
	highlightStyle      = lg.NewStyle().Bold(true).Foreground(lg.Color(primaryHighlight)).Italic(true)
	choicesStyle        = lg.NewStyle()
	choicesFocusedStyle = lg.NewStyle().Bold(true).Foreground(lg.Color(secondaryHighlight))
	helpStyle           = lg.NewStyle().AlignHorizontal(lg.Center).MarginTop(1)

	errorStyle        = lg.NewStyle().Bold(true).Foreground(lg.Color(err))
	confirmationStyle = lg.NewStyle().Foreground(lg.Color(confirm))
	focusedStyle      = lg.NewStyle().Foreground(lg.Color(primaryHighlight))
	cursorStyle       = focusedStyle
	noStyle           = lg.NewStyle()
	formBorderStyle   = lg.NewStyle().Border(lg.RoundedBorder()).Padding(1)

	listStyle                = lg.NewStyle().Width(20)
	listItemStyle            = lg.NewStyle().Border(lg.NormalBorder(), false, false, false, true).MarginTop(1).PaddingLeft(1)
	listItemHighlightStyle   = listItemStyle.BorderForeground(lg.Color(primaryHighlight))
	listItemDescriptionStyle = lg.NewStyle().Italic(true).Foreground(lg.Color(secondaryFg))
)
