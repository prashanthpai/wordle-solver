package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/prashanthpai/wordle-solver/wordle"
)

const usageText = `Input line format:
<WORD> <COLOR>

where COLOR is a number representing color of each letter:
0: GREY     The letter is not in the word in any spot.
1: YELLOW   The letter is in the word but in the wrong spot.
2: GREEN    The letter is in the word and in the correct spot.

Example:
TROVE 12001`

func main() {
	w, err := wordle.New()
	if err != nil {
		log.Fatalf("Unable to load wordle solver: %v", err)
	}

	fmt.Println(usageText)
	fmt.Printf("\nSeed words: %s\n> ", w.Seed())

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		word, colors, err := parseInputLine(line, w.WordLength())
		if err != nil {
			fmt.Printf("parsing error: %s\n> ", err.Error())
			continue
		}

		next, err := w.Next(word, colors)
		if err != nil {
			fmt.Printf("parsing error: %s\n> ", err.Error())
			continue
		}

		fmt.Printf("Word: %s\nNext: %v\n> ", w.Answer(), next)
		if len(next) == 0 || len(next) == 1 {
			fmt.Println("the end")
			break
		}
	}

	if scanner.Err() != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}

func parseInputLine(line string, wordLen int) ([]rune, []wordle.Color, error) {
	word := make([]rune, wordLen)
	colors := make([]wordle.Color, wordLen)

	s := strings.Split(line, " ")
	if len(s) != 2 || len(s[0]) != wordLen || len(s[1]) != wordLen {
		return word, colors, fmt.Errorf("Invalid input. Example: TROVE 00210")
	}

	for i, r := range s[0] {
		word[i] = r
	}

	for i, r := range s[1] {
		num := r - '0'
		colors[i] = wordle.Color(num)
	}

	return word, colors, nil
}
