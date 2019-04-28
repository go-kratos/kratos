package random

import (
	"math"
	"math/rand"
	"time"
)

var (
	rnd *rand.Rand
	ch  chan int64
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	ch = make(chan int64, 1000)
	go randomBase(ch)
}

func randomBase(c chan int64) {
	for {
		c <- rnd.Int63()
	}
}

//Uniqid 随机数，length是需要返回的长度，只支持10~19位
func Uniqid(length int) int64 {
	if length < 10 || length > 19 {
		return 0
	}
	prefix := (time.Now().UnixNano() / 100000000) & 0x3fffffff
	cut := int64(math.Pow10(length - 9))
	suffix := <-ch % cut
	return prefix*cut + suffix
}
