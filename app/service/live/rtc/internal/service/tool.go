package service

import (
	"bytes"
	"math/rand"
	"time"
)

type Tool struct {
	rnd  *rand.Rand
	dict []byte
}

func NewTool() *Tool {
	t := &Tool{
		rnd:  rand.New(rand.NewSource(time.Now().UnixNano())),
		dict: []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}
	return t
}

func (t *Tool) RandomString(length int) string {
	var result bytes.Buffer
	result.Grow(length)
	for i := 0; i < length; i++ {
		result.WriteByte(t.dict[t.rnd.Intn(len(t.dict))])
	}
	return result.String()
}
