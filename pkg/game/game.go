package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jayjyli/openwordle-go/data"
)

type Correctness string

const (
	Correct   Correctness = "correct"
	Partial               = "partial"
	Incorrect             = "incorrect"
)

// An instance of a Wordle game. This stores state,
// e.g. the guess word and previous guess states.
// Assumption: we're playing with a 5-letter word.
// TODO: make this work for games using words that
// aren't 5 letters.
type Game struct {
	runes  []rune
	counts map[rune]uint
}

func NewGame() *Game {
	word := randomWord()
	runes := []rune(word)
	// A map of the number of times a letter shows up
	// in the word. This is important because a guess
	// containing 2 E's are both valid if the answer
	// also contains >=2 E's.
	counts := map[rune]uint{}

	for i := 0; i < len(word); i++ {
		// This works because of Go's automatic 0 value.
		// Important to note that no char would make it in
		// this map with an initial count of 0.
		counts[runes[i]]++
	}

	return &Game{
		runes:  runes,
		counts: counts,
	}
}

func (g *Game) Guess(guess string) ([]Correctness, error) {
	if err := g.validate(guess); err != nil {
		return nil, fmt.Errorf("invalid guess: %v", err)
	}

	// Answer correctness array
	c := make([]Correctness, len(guess))

	if string(g.runes) == guess {
		for i := 0; i < len(c); i++ {
			c[i] = Correct
		}

		return c, nil
	}

	// Copy the counts, don't modify the game word's state
	counts := copyMap(g.counts)
	gr := []rune(guess)
	// Process each letter in the guess
	for i := 0; i < len(gr); i++ {
		if _, ok := counts[gr[i]]; !ok {
			// This letter isn't in the answer word
			c[i] = Incorrect
		} else {
			// This letter is in the answer word
			if g.runes[i] == gr[i] {
				// This letter is in the correct position
				c[i] = Correct
			} else {
				// This letter is in the incorrect position
				c[i] = Partial
			}

			// Reduce the number of times this letter can be guessed
			counts[gr[i]]--
			if counts[gr[i]] == 0 {
				// Remove the letter from the guessable counts if it's the last time.
				// Important: there should be no letters with counts of 0 at the
				// beginning of a guess.
				delete(counts, gr[i])
			}
		}
	}

	return c, nil
}

func (g *Game) validate(guess string) error {
	if len(g.runes) != len(guess) {
		return errors.New("guess length does not match word length")
	}

	if _, ok := data.ALLOWED_GUESSES[guess]; !ok {
		return errors.New("not a valid word")
	}

	return nil
}

func randomWord() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(data.ANSWERS))
	return data.ANSWERS[i]
}

func copyMap(i map[rune]uint) map[rune]uint {
	o := map[rune]uint{}
	for k, v := range i {
		o[k] = v
	}

	return o
}
