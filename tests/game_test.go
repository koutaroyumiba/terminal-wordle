package game_tests

import (
	"koutaroyumiba/wordle/game"
	"testing"
)

var input_file string = "../data/wordle-answers-alphabetical.txt"

func TestInit(t *testing.T) {
	words := game.ProcessFile(input_file)
	t.Run("testing initialisation of game", func(t *testing.T) {
		gs := game.InitGame(words)
		if gs.GetAttempts() != 0 {
			t.Errorf("what attempt do we start with?! %d", gs.GetAttempts())
		}
	})
}

func TestWordleWordElate(t *testing.T) {
	t.Run("testing elate", func(t *testing.T) {
		gs := game.InitGameWithWord("elate")
		if gs.GetAttempts() != 0 {
			t.Errorf("what attempt do we start with?! %d", gs.GetAttempts())
		}
		some, _ := gs.Guess("geese")
		if some != "-^--x" {
			t.Errorf("geese: should be -^--x, got %s", some)
		}
		some, _ = gs.Guess("teeth")
		if some != "-^^x-" {
			t.Errorf("teeth: should be -^^x-, got %s", some)
		}
	})
}
