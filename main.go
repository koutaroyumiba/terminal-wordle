package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"koutaroyumiba/wordle/tui"
)

func main() {
	p := tea.NewProgram(tui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
}
