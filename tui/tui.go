package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"koutaroyumiba/wordle/game"
)

const (
	wordLength = 5
	maxGuesses = 6
)

var (
	// styles
	greenStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#6aaa64")).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1)
	yellowStyle = lipgloss.NewStyle().Background(lipgloss.Color("#c9b458")).Foreground(lipgloss.Color("#000000")).Padding(0, 1)
	grayStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#787c7e")).Foreground(lipgloss.Color("#ffffff")).Padding(0, 1)
	emptyStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#121212")).Foreground(lipgloss.Color("#888888")).Padding(0, 1)

	keyStyle    = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder()).Margin(0, 1)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffffff")).MarginBottom(1)
)

type model struct {
	gameState game.GameState
	current   []rune
	done      bool
	win       bool
	message   string
}

func initialModel() model {
	wordle := game.InitGame(wordLength, maxGuesses)

	return model{
		gameState: wordle,
		current:   []rune{},
		done:      false,
		win:       false,
		message:   "Type letters, Backspace to delete, Enter to submit.",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		// respond to q to quit or r to restart, or any key to exit
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "r", "R":
				return initialModel(), tea.ClearScreen
			case "q", "Q", "ctrl+c":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			r := msg.Runes[0]
			if len(m.current) < wordLength && ((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				m.current = append(m.current, rune(strings.ToLower(string(r))[0]))
				m.message = ""
			}
			return m, nil
		case tea.KeyBackspace:
			if len(m.current) > 0 {
				m.current = m.current[:len(m.current)-1]
			}
			m.message = ""
			return m, nil
		case tea.KeyEnter:
			// submit guess
			if len(m.current) != wordLength {
				m.message = fmt.Sprintf("Guess must be %d letters.", wordLength)
				return m, nil
			}
			guess := string(m.current)

			validateRes, errMsg := m.gameState.ValidateWord(guess)

			if !validateRes {
				m.message = errMsg
				return m, nil
			}

			// evaluate
			finished, won := m.gameState.EvaluateGuess(guess)
			m.current = []rune{}
			m.message = ""

			if finished {
				m.done = true
			}
			if won {
				m.win = true
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func renderCell(c game.Cell) string {
	char, state := c.GetInfo()
	ch := ' '
	if char != ' ' && char != 0 {
		ch = char
	}
	switch state {
	case game.StateCorrect:
		return greenStyle.Render(string(ch))
	case game.StatePresent:
		return yellowStyle.Render(string(ch))
	case game.StateAbsent:
		return grayStyle.Render(string(ch))
	default:
		return emptyStyle.Render(string(ch))
	}
}

func renderRow(cells []game.Cell) string {
	parts := make([]string, len(cells))
	for i, c := range cells {
		parts[i] = renderCell(c)
	}
	return strings.Join(parts, " ")
}

func renderKeyboard(known map[rune]game.CellState) string {
	// simple QWERTY rows
	rows := []string{
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
	}
	outRows := make([]string, len(rows))
	for ri, row := range rows {
		parts := []string{}
		for _, ch := range row {
			s, ok := known[ch]
			cellRep := string(ch)
			switch {
			case ok && s == game.StateCorrect:
				parts = append(parts, greenStyle.Render(cellRep))
			case ok && s == game.StatePresent:
				parts = append(parts, yellowStyle.Render(cellRep))
			case ok && s == game.StateAbsent:
				parts = append(parts, grayStyle.Render(cellRep))
			default:
				parts = append(parts, emptyStyle.Render(cellRep))
			}
		}
		outRows[ri] = strings.Join(parts, " ")
	}
	return strings.Join(outRows, "\n")
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(headerStyle.Render("Terminal Wordle"))
	b.WriteString("\n")

	// render guesses so far
	for i := range maxGuesses {
		b.WriteString(renderRow(m.gameState.GetCurrentBoardRow(m.current, i)))
		b.WriteString("\n\n")
	}

	// keyboard
	b.WriteString("Keyboard:\n")
	b.WriteString(renderKeyboard(m.gameState.GetKnown()))
	b.WriteString("\n\n")

	// message
	if m.message != "" {
		b.WriteString("msg: ")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(m.message))
		b.WriteString("\n")
	}

	winningStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6aaa64"))
	losingStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5f87"))

	if m.done {
		if m.win {
			b.WriteString(winningStyle.Render("\ncongrats\n"))
		} else {
			b.WriteString(losingStyle.Render(fmt.Sprintf("\ngg u suck, word: %s\n", m.gameState.GetAnswer())))
		}
		b.WriteString("\nPress r to play again, q to quit.\n")
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
}
