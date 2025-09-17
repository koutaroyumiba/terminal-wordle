package game

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Preprocess() []string {
	file, err := os.Open("valid-wordle-words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) != 5 {
			fmt.Println("what happened")
		}
		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return words
}
