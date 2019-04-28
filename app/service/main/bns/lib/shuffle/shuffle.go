package shuffle

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Shuffler A type, typically a collection, that satisfies Shuffler
// can be shuffle by Shuffle func
type Shuffler interface {
	Len() int
	Swap(i, j int)
}

// Shuffle s
func Shuffle(s Shuffler) {
	l := s.Len()
	for i := l; i > 0; i-- {
		j := rand.Intn(i)
		s.Swap(l-i, j)
	}
}
