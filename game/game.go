package game

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type GameState struct {
	answer         string
	guesses        []string // len = attempts
	guessesResults []string
	alphabet       []rune
}

func InitGame(words []string) GameState {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWord := "lmfao"
	if len(words) > 0 {
		randomWord = words[rng.Intn(len(words))]
		fmt.Println("random word chosen")
	} else {
		fmt.Println("no words found")
	}

	return GameState{
		answer:   randomWord,
		alphabet: []rune("abcdefghijklmnopqrstuvwxyz"),
	}
}

func InitGameWithWord(correctWord string) GameState {
	return GameState{
		answer:   correctWord,
		alphabet: []rune("abcdefghijklmnopqrstuvwxyz"),
	}
}

func (g *GameState) Guess(guess string) (string, bool) {
	g.guesses = append(g.guesses, guess)

	green := "x"
	yellow := "^"
	black := "-"

	guessResult := []string{black, black, black, black, black}
	guessIndexUsed := []int{}
	answerIndexUsed := []int{}

	// find all green
	for index, char := range guess {
		alphabetIndex := char - 'a'
		g.alphabet[alphabetIndex] = '-'
		if guess[index] == g.answer[index] && !slices.Contains(guessIndexUsed, index) {
			guessIndexUsed = append(guessIndexUsed, index)
			answerIndexUsed = append(answerIndexUsed, index)
			guessResult[index] = green
			g.alphabet[alphabetIndex] = char
		}
	}
	// find all yellow
	for gIndex, char := range guess {
		alphabetIndex := char - 'a'
		if !slices.Contains(guessIndexUsed, gIndex) {
			for aIndex, _ := range g.answer {
				if !slices.Contains(answerIndexUsed, aIndex) {
					if guess[gIndex] == g.answer[aIndex] {
						guessIndexUsed = append(guessIndexUsed, gIndex)
						answerIndexUsed = append(answerIndexUsed, aIndex)
						guessResult[gIndex] = yellow
						g.alphabet[alphabetIndex] = char
						break
					}
				}
			}
		}
	}

	completeResult := strings.Join(guessResult, "")
	g.guessesResults = append(g.guessesResults, completeResult)

	if completeResult == strings.Repeat(green, len(guess)) {
		return completeResult, true
	}

	return completeResult, false
}

func (g GameState) GetAttempts() int {
	return len(g.guesses)
}

func (g GameState) GetLetters() string {
	return string(g.alphabet)
}

func (g GameState) PrintBoard() {
	fmt.Printf("== board (attempt %d) ==\n", len(g.guesses)+1)
	for index, _ := range g.guesses {
		fmt.Printf("\t%s\n", g.guesses[index])
		fmt.Printf("\t%s\n", g.guessesResults[index])
	}
}
