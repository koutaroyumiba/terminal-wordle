package main

import (
	"bufio"
	"fmt"
	"koutaroyumiba/wordle/game"
	"os"
	"slices"
)

func main() {
	choiceWords := game.ProcessFile("data/wordle-answers-alphabetical.txt")
	validWords := game.ProcessFile("data/valid-wordle-words.txt")
	wordle := game.InitGame(choiceWords)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n| %s |\n", wordle.GetLetters())
		wordle.PrintBoard()
		fmt.Println()
		scanner.Scan()
		var currGuess string
		currGuess = scanner.Text()
		if scanner.Err() != nil {
			fmt.Println("error scanning input")
		}

		if len(currGuess) != 5 {
			fmt.Println("not even 5 letters, do you even know how to play")
			continue
		}
		if !slices.Contains(validWords, currGuess) {
			fmt.Println("not a real word loser")
			continue
		}
		_, finished := wordle.Guess(currGuess)
		if finished {
			fmt.Printf("xxxxx\nguessed correctly\n")
			break
		}
	}
}
