// Package main : main entrypoint
package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"koutaroyumiba/wordle/ui"
)

func main() {
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
}
