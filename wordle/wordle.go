package wordle

import (
	"fmt"
	"sort"
)

const (
	defaultWordLen = 5
	freqThreshold  = 500
)

type Color int

const (
	Grey Color = iota
	Yellow
	Green
)

type Solver struct {
	dictionary         map[string]int
	lettersInWord      map[rune]struct{}
	lettersNotInWord   map[rune]struct{}
	lettersAtWrongSpot map[int]rune
	answer             []rune
}

// New returns a new instance of wordle solver.
func New() (*Solver, error) {
	s := &Solver{
		answer: make([]rune, defaultWordLen),
	}
	s.Reset()

	d, err := loadDict()
	if err != nil {
		return nil, err
	}
	s.dictionary = d

	return s, nil
}

// Reset resets internal state of the wordle solver enabling it to be
// reused for the next wordle puzzle.
func (s *Solver) Reset() {
	s.lettersInWord = make(map[rune]struct{})
	s.lettersNotInWord = make(map[rune]struct{})
	s.lettersAtWrongSpot = make(map[int]rune)

	for i := range s.answer {
		s.answer[i] = '_'
	}
}

// Seed returns a list of suggested seed words.
func (s *Solver) Seed() []string {
	seeds := []string{
		"AISLE", "RAISE", "RATIO", "TRAIN", "AROSE",
		"LASER", "REALS", "EARLS", "POISE", "SLATE",
	}

	return seeds
}

// Answer returns the answer word in its current state containing
// letters that known to be present in its final position (Green).
func (s *Solver) Answer() string {
	return string(s.answer[:])
}

// WordLength returns length of the words allowed in wordle. Defaults
// to 5.
func (s *Solver) WordLength() int {
	return len(s.answer)
}

// Next returns suggestions for entering next into wordle.
func (s *Solver) Next(word []rune, colors []Color) ([]string, error) {

	// validate input
	if len(word) != len(s.answer) || len(colors) != len(s.answer) {
		return nil, fmt.Errorf("word/color length != %d", len(s.answer))
	}
	for i := range word {
		color := colors[i]
		if color != Green && color != Yellow && color != Grey {
			return nil, fmt.Errorf("Invalid color: %d", color)
		}
	}

	for i, letter := range word {
		color := colors[i]
		switch color {
		case Green:
			s.lettersInWord[letter] = struct{}{}
			s.answer[i] = letter
		case Yellow:
			s.lettersInWord[letter] = struct{}{}
			s.lettersAtWrongSpot[i] = letter
		case Grey:
			_, ok := s.lettersInWord[letter]
			if !ok {
				// edge case
				s.lettersNotInWord[letter] = struct{}{}
			}
		}
	}

	return s.nextWords(), nil
}

// nextWords goes through the dictionary and returns words
// that are candidates to be the next word sorted by their
// popularity.
func (s *Solver) nextWords() []string {
	var words []string
	for word := range s.dictionary {
		if s.canBeNext(word) {
			words = append(words, word)
		}
	}

	return s.sortByFrequency(words)
}

// canBeNext evaluates if the input word satisfies all
// conditions necessary to be a candidate for next word.
func (s *Solver) canBeNext(next string) bool {
	word := []rune(next)

	if len(word) != len(s.answer) {
		return false
	}

	for i := range word {
		// skip if known letters at position don't match
		if word[i] != s.answer[i] && s.answer[i] != '_' {
			return false
		}

		// skip if the word contains letters
		// that are known to be absent in the final answer.
		_, contains := s.lettersNotInWord[word[i]]
		if contains {
			return false
		}

		// skip if the word contains letters at places
		// that are known to be incorrect.
		wrongLetter, ok := s.lettersAtWrongSpot[i]
		if ok && word[i] == wrongLetter {
			return false
		}
	}

	// skip if the word doesn't have letters
	// that are known be present in the final answer.
	for letter := range s.lettersInWord {
		if !containsRune(word, letter) {
			return false
		}
	}

	return true
}

// sortByFrequency sorts the words by how frequent they occur in literature.
func (s *Solver) sortByFrequency(words []string) []string {
	var result []string

	for i := range words {
		f := s.dictionary[words[i]]
		if f > freqThreshold {
			result = append(result, words[i])
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return s.dictionary[result[i]] > s.dictionary[result[j]]
	})

	return result
}
