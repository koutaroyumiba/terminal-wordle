package game_tests

import (
	"fmt"
	"koutaroyumiba/wordle/bot"
	"koutaroyumiba/wordle/game"
	"testing"
)

func TestWordleGame1(t *testing.T) {
	wordLength := 5
	maxGuesses := 6
	wordle := game.InitGameWithWord(wordLength, maxGuesses, "spice")
	wordleBot := bot.InitBot(wordLength, maxGuesses)

	t.Run("guess #1: salet", func(t *testing.T) {
		finished, won := wordle.EvaluateGuess("salet")
		if finished {
			t.Error("finished too soon: 1.salet")
		}
		if won {
			t.Error("not won: 1.salet")
		}
	})
	t.Run("guess #2: spine", func(t *testing.T) {
		finished, won := wordle.EvaluateGuess("spine")
		if finished {
			t.Error("finished too soon: 2.salet")
		}
		if won {
			t.Error("not won: 2.spine")
		}
	})
	t.Run("guess #3: spire", func(t *testing.T) {
		finished, won := wordle.EvaluateGuess("spire")
		if finished {
			t.Error("finished too soon: 3.spire")
		}
		if won {
			t.Error("not won: 3.spire")
		}
	})
	t.Run("guess #4: spike", func(t *testing.T) {
		finished, won := wordle.EvaluateGuess("spike")
		if finished {
			t.Error("finished too soon: 4.spike")
		}
		if won {
			t.Error("not won: 4.spike")
		}
	})
	t.Run("guess #5: spice", func(t *testing.T) {
		finished, won := wordle.EvaluateGuess("spice")
		if !finished {
			t.Error("why is it not finished 5.spice")
		}
		if !won {
			t.Error("why u no win 5.spice")
		}
	})

	stuff := wordleBot.Analysis(wordle.GetGuesses())
	for _, c := range stuff {
		fmt.Printf("c: %d\n", c)
	}
}
