package game

import (
	"math/rand"
	"time"

	"github.com/jayjyli/openwordle-go/data"
)

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
