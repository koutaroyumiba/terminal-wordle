package bot

import "koutaroyumiba/wordle/game"

var (
	dictionaryInputFile = "data/valid-wordle-words.txt"
	dictionary          = game.ProcessFile(dictionaryInputFile)
)

type WordleBot struct {
	wordLength int
	maxGuesses int
}

func InitBot(wordLength, maxGuesses int) WordleBot {
	return WordleBot{
		wordLength: wordLength,
		maxGuesses: maxGuesses,
	}
}

func (w WordleBot) Analysis(guesses [][]game.Cell) []int {
	result := make([]int, w.maxGuesses)
	validWords := dictionary
	for rowIndex := range w.maxGuesses {
		newValidWords := []string{}
		currGuess := guesses[rowIndex]
		if rowIndex == 0 || rowIndex > 0 && len(validWords) != result[rowIndex-1] {
			result[rowIndex] = len(validWords)
		}

		for _, word := range validWords {
			if isValid(currGuess, word) {
				newValidWords = append(newValidWords, word)
			}
		}

		validWords = newValidWords
	}

	return result
}

func isValid(guess []game.Cell, word string) bool {
	guessResult := make([]game.CellState, len(guess))
	answerRunes := []rune(word)
	counts := map[rune]int{}

	// find all green
	for i := range guess {
		char, _ := guess[i].GetInfo()
		if answerRunes[i] == char {
			guessResult[i] = game.StateCorrect
		} else {
			counts[answerRunes[i]]++
		}
	}

	// second pass (for yellow)
	for i := range guess {
		char, _ := guess[i].GetInfo()
		if guessResult[i] == game.StateCorrect {
			continue
		}
		if counts[char] > 0 {
			guessResult[i] = game.StatePresent
			counts[char]--
		} else {
			guessResult[i] = game.StateAbsent
		}
	}

	for i := range guessResult {
		_, state := guess[i].GetInfo()
		if guessResult[i] != state {
			return false
		}
	}

	return true
}
