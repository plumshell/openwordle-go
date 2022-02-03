package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jayjyli/openwordle-go/pkg/game"
)

func main() {
	g := game.NewGame()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()

		c, err := g.Guess(line)
		if err != nil {
			fmt.Printf("failed to guess: %v\n", err)
		} else {
			fmt.Println(c)
		}
	}
}
