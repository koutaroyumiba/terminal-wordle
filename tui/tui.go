package main

import (
	"fmt"
	"strings"

	"koutaroyumiba/wordle/bot"
	"koutaroyumiba/wordle/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	return tea.ClearScreen
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
	b.WriteString(headerStyle.Render("Terminal Wordle (ctrl+c to exit)"))
	b.WriteString("\n")

	guesses := m.gameState.GetGuesses()
	bot := bot.InitBot(wordLength, maxGuesses)
	length := bot.Analysis(guesses)

	// render guesses so far
	for i := range maxGuesses {
		b.WriteString(fmt.Sprintf("%s  no. of words left: %d", renderRow(m.gameState.GetCurrentBoardRow(m.current, i)), length[i]))
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

	stats := m.gameState.GetStats()

	if m.done {
		if m.win {
			b.WriteString(winningStyle.Render("\ncongrats\n"))
		} else {
			b.WriteString(losingStyle.Render(fmt.Sprintf("\ngg u suck, word: %s\n", m.gameState.GetAnswer())))
		}
		b.WriteString("\nPress r to play again, q to quit.\n")

	}

	b.WriteString("\n--- Statistics ---\n")
	b.WriteString(fmt.Sprintf("Games Played: %d\n", stats.GamesPlayed))
	b.WriteString(fmt.Sprintf("Wins: %d\n", stats.Wins))
	b.WriteString(fmt.Sprintf("Win Rate: %.1f%%\n", stats.WinRate()))
	b.WriteString(fmt.Sprintf("Current Streak: %d\n", stats.CurrentStreak))
	b.WriteString(fmt.Sprintf("Max Streak: %d\n", stats.MaxStreak))
	b.WriteString(fmt.Sprintf("Avg Guesses (wins): %.2f\n", stats.AverageGuesses()))

	distribution := stats.GuessFrequency
	total := 0
	for _, c := range distribution {
		total += c
	}

	for i := range maxGuesses {
		count, ok := distribution[i+1]
		if !ok {
			count = 0
		}
		b.WriteString(fmt.Sprintf("%d : %s[%d]\n", i+1, strings.Repeat("#", int(float64(count)/float64(total)*30)), count))
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
	}
}
