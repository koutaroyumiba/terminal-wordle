// main.go
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

type cell struct {
	r rune
	s game.CellState
}

type model struct {
	gameState       game.GameState
	secret          string
	guesses         [][]cell
	current         []rune
	row             int
	done            bool
	win             bool
	message         string
	knownLetters    map[rune]game.CellState // for keyboard hints
	allowDictionary bool                    // if true, only allow guesses in wordList
}

func initialModel() model {
	wordle := game.InitGame(wordLength, maxGuesses)
	guesses := make([][]cell, maxGuesses)
	for i := range maxGuesses {
		line := make([]cell, wordLength)
		for j := range wordLength {
			line[j] = cell{r: ' ', s: game.StateEmpty}
		}
		guesses[i] = line
	}
	return model{
		gameState:       wordle,
		secret:          wordle.GetAnswer(),
		guesses:         guesses,
		current:         []rune{},
		row:             0,
		done:            false,
		win:             false,
		message:         "Type letters, Backspace to delete, Enter to submit.",
		knownLetters:    make(map[rune]game.CellState),
		allowDictionary: true, // set to false to skip dictionary check
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func runeSliceToString(rs []rune) string {
	return string(rs)
}

func containsWord(list []string, word string) bool {
	for _, w := range list {
		if strings.EqualFold(w, word) {
			return true
		}
	}
	return false
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		// respond to q to quit or r to restart, or any key to exit
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "r", "R":
				return initialModel(), nil
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
				for i := range len(m.current) {
					m.guesses[m.row][i].r = m.current[i]
				}
				m.message = ""
			}
			return m, nil
		case tea.KeyBackspace:
			if len(m.current) > 0 {
				m.current = m.current[:len(m.current)-1]
				for i := range len(m.current) {
					m.guesses[m.row][i].r = m.current[i]
				}
			}
			m.message = ""
			return m, nil
		case tea.KeyEnter:
			// submit guess
			if len(m.current) != wordLength {
				m.message = fmt.Sprintf("Guess must be %d letters.", wordLength)
				return m, nil
			}
			guess := runeSliceToString(m.current)

			// TODO: need to be refactored - evaluation should be done in GameState
			choiceWords := []string{"apple"}

			if m.allowDictionary && !containsWord(choiceWords, guess) {
				m.message = "Not in word list."
				return m, nil
			}
			// evaluate
			states, _ := m.gameState.EvaluateGuess(guess)
			for i := range wordLength {
				m.guesses[m.row][i].r = m.current[i]
				m.guesses[m.row][i].s = states[i]
				updateKnownLetter(m.knownLetters, m.current[i], states[i])
			}
			if allCorrect(states) {
				m.done = true
				m.win = true
				m.message = ""
				return m, nil
			}
			m.row++
			m.current = []rune{}
			if m.row >= maxGuesses {
				m.done = true
				m.win = false
				m.message = ""
				return m, nil
			}
			m.message = ""
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func updateKnownLetter(known map[rune]game.CellState, r rune, s game.CellState) {
	prev, ok := known[r]
	// priority: correct > present > absent
	if !ok || s > prev {
		known[r] = s
	}
}

func allCorrect(states []game.CellState) bool {
	for _, s := range states {
		if s != game.StateCorrect {
			return false
		}
	}
	return true
}

func renderCell(c cell) string {
	ch := ' '
	if c.r != ' ' && c.r != 0 {
		ch = c.r
	}
	switch c.s {
	case game.StateCorrect:
		return greenStyle.Render(string(ch))
	case game.StatePresent:
		return yellowStyle.Render(string(ch))
	case game.StateAbsent:
		return grayStyle.Render(string(ch))
	default:
		return emptyStyle.Render(" ")
	}
}

func renderRow(cells []cell) string {
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
	b.WriteString(headerStyle.Render("Wordle"))
	b.WriteString("\n")

	// render guesses so far
	for i := range maxGuesses {
		// if current row, display current typed letters and empties
		if i == m.row && !m.done {
			line := make([]cell, wordLength)
			for j := range wordLength {
				if j < len(m.current) {
					line[j] = cell{r: m.current[j], s: game.StateEmpty}
				} else {
					line[j] = cell{r: 'x', s: game.StateEmpty}
				}
			}
			b.WriteString(renderRow(line))
		} else {
			b.WriteString(renderRow(m.guesses[i]))
		}
		b.WriteString("\n\n")
	}

	// keyboard
	b.WriteString("Keyboard:\n")
	b.WriteString(renderKeyboard(m.knownLetters))
	b.WriteString("\n\n")

	b.WriteString("input: ")
	b.WriteString(string(m.current))
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
			b.WriteString(losingStyle.Render(fmt.Sprintf("\ngg u suck, word: %s\n", m.secret)))
		}
		b.WriteString("\nPress r to play again, q to quit.\n")
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
}
