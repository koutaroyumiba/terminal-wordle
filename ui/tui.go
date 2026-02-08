// Package ui : package for the TUI
package ui

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

func InitialModel() model {
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
				return InitialModel(), tea.ClearScreen
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
			finished, won := m.gameState.ApplyGuess(guess)
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

func (m model) renderGame() string {
	var lines []string
	header := headerStyle.Render("Terminal Wordle (ctrl+c to exit)")
	lines = append(lines, header, "")

	guesses := m.gameState.GetGuesses()
	bot := bot.InitBot(wordLength, maxGuesses)
	length, words := bot.Analysis(guesses)

	for i := range maxGuesses {
		row := renderRow(m.gameState.GetCurrentBoardRow(m.current, i))
		row += fmt.Sprintf("  no. of words left: %d", length[i])
		if m.done && len(words[i]) > 0 && len(words[i]) < 8 {
			row += " ["
			for _, word := range words[i] {
				row += fmt.Sprintf(" %s ", word)
			}
			row += "]"
		}
		lines = append(lines, row)
	}

	// keyboard
	lines = append(lines, "", "Keyboard:")
	lines = append(lines, strings.Split(renderKeyboard(m.gameState.GetKnown()), "\n")...)

	// end game message
	if m.done {
		if m.win {
			lines = append(lines, "", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6aaa64")).Render("congrats"))
		} else {
			lines = append(lines, "", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5f87")).Render(fmt.Sprintf("gg u suck, word: %s", m.gameState.GetAnswer())))
		}
		lines = append(lines, "Press r to play again, q to quit.")
	}

	return strings.Join(lines, "\n")
}

func (m model) renderStats() string {
	stats := m.gameState.GetStats()
	var lines []string
	lines = append(lines, headerStyle.Render("Statistics"), "")
	lines = append(lines, fmt.Sprintf("Games Played: %d", stats.GamesPlayed))
	lines = append(lines, fmt.Sprintf("Wins: %d", stats.Wins))
	lines = append(lines, fmt.Sprintf("Win Rate: %.1f%%", stats.WinRate()))
	lines = append(lines, fmt.Sprintf("Current Streak: %d", stats.CurrentStreak))
	lines = append(lines, fmt.Sprintf("Max Streak: %d", stats.MaxStreak))
	lines = append(lines, fmt.Sprintf("Avg Guesses (wins): %.2f", stats.AverageGuesses()))

	distribution := stats.GuessFrequency
	total := 0
	for _, c := range distribution {
		total += c
	}
	lines = append(lines, "", "Guess Distribution:")
	for i := range maxGuesses {
		count, ok := distribution[i+1]
		if !ok {
			count = 0
		}
		bar := strings.Repeat("#", int(float64(count)/float64(total)*30))
		lines = append(lines, fmt.Sprintf("%d : %s [%d]", i+1, bar, count))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderMessage() string {
	if m.message == "" {
		return ""
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Render("msg: " + m.message)
}

func (m model) View() string {
	gameView := m.renderGame()
	statsView := m.renderStats()
	messageView := m.renderMessage()

	// Combine game and stats horizontally
	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		gameView,
		lipgloss.NewStyle().MarginLeft(3).Render(statsView),
	)

	// Then put the message at the bottom
	if messageView != "" {
		return lipgloss.JoinVertical(lipgloss.Left, mainContent, "", messageView)
	}

	return mainContent
}
