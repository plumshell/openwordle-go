package game

import (
	"errors"
	"fmt"

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
	guesses uint
	limit   uint

	runes  []rune        // The word, as a char array
	counts map[rune]uint // Letter counts
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
		guesses: 0,
		limit:   6, // A standard game is 6 guess attempts
		runes:   runes,
		counts:  counts,
	}
}

func (g *Game) Remaining() uint {
	return g.limit - g.guesses
}

func (g *Game) Guess(guess string) ([]Correctness, error) {
	if g.limit <= g.guesses {
		return nil, errors.New("guess limit reached, please start a new game")
	}

	if err := g.validate(guess); err != nil {
		return nil, fmt.Errorf("invalid guess: %v", err)
	}

	g.guesses++

	// Answer correctness array
	c := make([]Correctness, len(guess))

	// If the answer was solved
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
