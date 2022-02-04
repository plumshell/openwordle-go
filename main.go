package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jayjyli/openwordle-go/pkg/game"
)

func main() {
	g := game.NewGame()
	fmt.Println(`Welcome to OpenWordle!

How to play:
- You have 6 turns to guess a 5-letter word.
- Each turn, you make a guess for what the word is. If you guess the word, you win.
- If you guess a nonexistent word, e.g. 'aaaaa', it will not count as a valid guess.
- When you guess a word, you'll receive feedback on how correct each letter was:
	- _____ means none of the letters in your guess were correct
	- _e___ means the answer contains the letter 'E', but in a different spot
	- _E___ means the answer contains the letter 'E' in the spot you guessed

Have fun!`)

	s := bufio.NewScanner(os.Stdin)
	fmt.Printf("Please make a guess (you have %d left): ", g.Remaining())
	for s.Scan() {
		in := s.Text()

		c, err := g.Guess(in)
		if err != nil {
			fmt.Printf("failed to guess: %v\n", err)
		} else {
			fmt.Println(c)
		}

		fmt.Printf("Please make a guess (you have %d left): ", g.Remaining())
	}
}
