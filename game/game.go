package game

import (
	"fmt"
	"math/rand"
	"time"
)

type CellState int

var inputFile string = "data/wordle-answers-alphabetical.txt"

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

type GameState struct {
	answer         string
	guessesResults [][]Cell
	knownLetters   map[rune]CellState
}

func InitGame(wordLength, maxGuesses int) GameState {
	words := ProcessFile(inputFile)
	secret := pickRandomWord(words)
	board := initialiseEmptyBoard(wordLength, maxGuesses)

	return GameState{
		answer:         secret,
		guessesResults: board,
		knownLetters:   make(map[rune]CellState),
	}
}

func pickRandomWord(words []string) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomWord := "lmfao"
	if len(words) > 0 {
		randomWord = words[rng.Intn(len(words))]
		fmt.Println("random word chosen")
	} else {
		fmt.Println("no words found")
	}

	return randomWord
}

func InitGameWithWord(wordLength, maxGuesses int, correctWord string) GameState {
	board := initialiseEmptyBoard(wordLength, maxGuesses)

	return GameState{
		answer:         correctWord,
		guessesResults: board,
		knownLetters:   make(map[rune]CellState),
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

func (g *GameState) Guess(guess string) ([]CellState, bool) {
	// g.guesses = append(g.guesses, guess)

	guessResult := make([]CellState, len(guess))
	answerRunes := []rune(g.answer)
	guessRunes := []rune(guess)
	counts := map[rune]int{}

	// find all green
	for i := range len(guess) {
		if answerRunes[i] == guessRunes[i] {
			guessResult[i] = StateCorrect
		} else {
			counts[answerRunes[i]]++
		}
	}

	// second pass
	for i := range len(guess) {
		if guessResult[i] == StateCorrect {
			continue
		}

		if counts[guessRunes[i]] > 0 {
			guessResult[i] = StatePresent
			counts[guessRunes[i]]--
		} else {
			guessResult[i] = StateAbsent
		}
	}

	// g.guessesResults = append(g.guessesResults, guessResult)

	for _, s := range guessResult {
		if s != StateCorrect {
			return guessResult, false
		}
	}

	return guessResult, true
}

// func (g GameState) GetAttempts() int {
// 	return len(g.guesses)
// }
//
// func (g GameState) GetLetters() string {
// 	return string(g.alphabet)
// }
//
// func (g GameState) PrintBoard() {
// 	fmt.Printf("== board (attempt %d) ==\n", len(g.guesses)+1)
// 	for index, _ := range g.guesses {
// 		fmt.Printf("\t%s\n", g.guesses[index])
// 		fmt.Printf("\t%p\n", g.guessesResults[index])
// 	}
// }

func (g GameState) GetAnswer() string {
	return g.answer
}
