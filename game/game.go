package game

import (
	"fmt"
	"math/rand"
	"slices"
	"time"
)

type CellState int

var (
	dictionaryInputFile   = "data/valid-wordle-words.txt"
	dictionary            = ProcessFile(dictionaryInputFile)
	validAnswersInputFile = "data/wordle-answers-alphabetical.txt"
	validAnswers          = ProcessFile(validAnswersInputFile)
)

const (
	StateEmpty CellState = iota
	StateCorrect
	StatePresent
	StateAbsent
)

type Cell struct {
	char  rune
	state CellState
}

func (c Cell) GetInfo() (rune, CellState) {
	return c.char, c.state
}

type GameState struct {
	stats           Stats
	answer          string
	guessesResults  [][]Cell
	knownLetters    map[rune]CellState
	wordLength      int
	maxGuesses      int
	allowDictionary bool
	currentRow      int
	finished        bool
}

func InitGame(wordLength, maxGuesses int) GameState {
	secret := pickRandomWord(validAnswers)
	board := initialiseEmptyBoard(wordLength, maxGuesses)

	return GameState{
		stats:           loadStats(),
		answer:          secret,
		guessesResults:  board,
		knownLetters:    make(map[rune]CellState),
		wordLength:      wordLength,
		maxGuesses:      maxGuesses,
		allowDictionary: true,
		currentRow:      0,
		finished:        false,
	}
}

func pickRandomWord(words []string) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomWord := "lmfao"
	if len(words) > 0 {
		randomWord = words[rng.Intn(len(words))]
	} else {
		fmt.Println("err: no words found")
	}

	return randomWord
}

func InitGameWithWord(wordLength, maxGuesses int, correctWord string) GameState {
	board := initialiseEmptyBoard(wordLength, maxGuesses)

	return GameState{
		stats:           loadStats(),
		answer:          correctWord,
		guessesResults:  board,
		knownLetters:    make(map[rune]CellState),
		wordLength:      wordLength,
		maxGuesses:      maxGuesses,
		allowDictionary: true,
		currentRow:      0,
		finished:        false,
	}
}

func initialiseEmptyBoard(wordLength, maxGuesses int) [][]Cell {
	board := make([][]Cell, maxGuesses)
	for i := range maxGuesses {
		line := make([]Cell, wordLength)
		for j := range wordLength {
			line[j] = Cell{char: ' ', state: StateEmpty}
		}
		board[i] = line
	}

	return board
}

func (g GameState) ValidateWord(word string) (bool, string) {
	if len(word) != g.wordLength {
		return false, fmt.Sprintf("guess must be %d letters", g.wordLength)
	}

	if g.allowDictionary && !slices.Contains(dictionary, word) {
		return false, "not in word list"
	}

	return true, ""
}

func (g *GameState) ApplyGuess(guess string) (bool, bool) {
	guessResult := EvaluateGuess([]rune(g.answer), []rune(guess))
	g.updateKnownLetter(guess, guessResult)
	g.updateState(guess, guessResult)

	won := false
	if isCorrectGuess(guessResult) {
		g.finished = true
		won = true

		g.stats.GamesPlayed++
		g.stats.Wins++
		g.stats.CurrentStreak++
		if g.stats.CurrentStreak > g.stats.MaxStreak {
			g.stats.MaxStreak = g.stats.CurrentStreak
		}

		_, ok := g.stats.GuessFrequency[g.currentRow]
		if ok {
			g.stats.GuessFrequency[g.currentRow]++
		} else {
			g.stats.GuessFrequency[g.currentRow] = 1
		}

		saveStats(g.stats)
	} else if g.currentRow >= g.maxGuesses {
		g.finished = true
		g.stats.GamesPlayed++
		g.stats.CurrentStreak = 0

		saveStats(g.stats)
	}

	return g.finished, won
}

func EvaluateGuess(answer, guess []rune) []CellState {
	result := make([]CellState, len(guess))
	counts := map[rune]int{}

	// find all green
	for i := range guess {
		if answer[i] == guess[i] {
			result[i] = StateCorrect
		} else {
			counts[answer[i]]++
		}
	}

	// second pass (for yellow)
	for i := range guess {
		if result[i] == StateCorrect {
			continue
		}

		if counts[guess[i]] > 0 {
			result[i] = StatePresent
			counts[guess[i]]--
		} else {
			result[i] = StateAbsent
		}
	}

	return result
}

func isCorrectGuess(guess []CellState) bool {
	for _, s := range guess {
		if s != StateCorrect {
			return false
		}
	}

	return true
}

func (g *GameState) updateKnownLetter(guess string, states []CellState) {
	for i := range g.wordLength {
		char := rune(guess[i])
		prev, ok := g.knownLetters[char]
		if !ok || states[i] < prev {
			g.knownLetters[char] = states[i]
		}
	}
}

func (g *GameState) updateState(guess string, states []CellState) {
	for i := range g.wordLength {
		g.guessesResults[g.currentRow][i].char = rune(guess[i])
		g.guessesResults[g.currentRow][i].state = states[i]
	}

	g.currentRow++
}

func (g GameState) GetCurrentBoardRow(currentWord []rune, index int) []Cell {
	if index != g.currentRow || g.finished {
		return g.guessesResults[index]
	}

	line := make([]Cell, g.wordLength)
	for i := range g.wordLength {
		char := ' '
		if i < len(currentWord) {
			char = currentWord[i]
		}
		line[i] = Cell{char: char, state: StateEmpty}
	}

	return line
}

func (g GameState) GetAnswer() string {
	return g.answer
}

func (g GameState) GetKnown() map[rune]CellState {
	return g.knownLetters
}

func (g GameState) GetStats() Stats {
	return g.stats
}

func (g GameState) GetGuesses() [][]Cell {
	return g.guessesResults
}
