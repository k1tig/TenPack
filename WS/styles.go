package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

var live = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))
var start = lipgloss.NewStyle().
	Foreground(lipgloss.Color("46"))
var end = lipgloss.NewStyle().
	Foreground(lipgloss.Color("124"))

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("86")
	lightGray = lipgloss.Color("87")

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14).Align(lipgloss.Center)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
	re           = lipgloss.NewRenderer(os.Stdout)
	baseStyle    = re.NewStyle().Padding(0, 1)
)
