package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jayjyli/openwordle-go/data"
)

type Correctness string

const (
	Correct   Correctness = "correct"
	Partial   Correctness = "partial"
	Incorrect Correctness = "incorrect"
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

func (g *Game) Guess(guess string) (string, error) {
	if g.limit <= g.guesses {
		return "", errors.New("guess limit reached, please start a new game")
	}

	if err := g.validate(guess); err != nil {
		return "", fmt.Errorf("invalid guess: %v", err)
	}

	// If the answer was solved
	if string(g.runes) == guess {
		return strings.ToUpper(guess), nil
	}

	r := []rune(guess)
	ca := g.check(r)
	var a string
	for i, c := range ca {
		char := string(r[i])

		switch c {
		case Correct:
			a += strings.ToUpper(char)
		case Partial:
			a += char
		case Incorrect:
			a += "_"
		default:
			return "", fmt.Errorf("unrecognized correctness value: %s", c)
		}
	}

	g.guesses++
	return a, nil
}

func (g *Game) check(gr []rune) []Correctness {
	c := make([]Correctness, len(gr))

	// Copy the counts, don't modify the game word's state
	counts := copyMap(g.counts)
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

	return c
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
