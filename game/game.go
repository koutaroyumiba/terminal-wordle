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

type GameState struct {
	answer          string
	guessesResults  [][]Cell
	knownLetters    map[rune]CellState
	wordLength      int
	allowDictionary bool
}

func InitGame(wordLength, maxGuesses int) GameState {
	secret := pickRandomWord(validAnswers)
	board := initialiseEmptyBoard(wordLength, maxGuesses)

	return GameState{
		answer:          secret,
		guessesResults:  board,
		knownLetters:    make(map[rune]CellState),
		wordLength:      wordLength,
		allowDictionary: true,
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
		answer:          correctWord,
		guessesResults:  board,
		knownLetters:    make(map[rune]CellState),
		wordLength:      wordLength,
		allowDictionary: true,
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

func (g *GameState) EvaluateGuess(guess string) ([]CellState, bool) {
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

	// second pass (for yellow)
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

	won := false
	if isCorrectGuess(guessResult) {
		won = true
	}

	return guessResult, won
}

func isCorrectGuess(guess []CellState) bool {
	for _, s := range guess {
		if s != StateCorrect {
			return false
		}
	}

	return true
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
