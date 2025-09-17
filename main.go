package main

import (
	"bufio"
	"fmt"
	"koutaroyumiba/wordle/game"
	"os"
	"slices"
)

func main() {
	words := game.Preprocess()
	wordle := game.InitGame(words)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n| %s |\n\n", wordle.GetLetters())
		scanner.Scan()
		var e string
		e = scanner.Text()
		if scanner.Err() != nil {
			// handle err
			fmt.Println("error")
		}

		if !slices.Contains(words, e) {
			fmt.Println("not a real word loser")
			continue
		}
		if len(e) != 5 {
			fmt.Println("not even 5 letters, do you even know how to play")
			continue
		}
		res, finished := wordle.Guess(e)
		fmt.Printf("%s (attempt = %d)\n", res, wordle.GetAttempts())
		if finished {
			fmt.Printf("guessed correctly\n")
			break
		}
	}
}
