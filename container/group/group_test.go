package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupGet(t *testing.T) {
	count := 0
	g := NewGroup(func() interface{} {
		count++
		return count
	})
	v := g.Get("key_0")
	assert.Equal(t, 1, v.(int))

	v = g.Get("key_1")
	assert.Equal(t, 2, v.(int))

	v = g.Get("key_0")
	assert.Equal(t, 1, v.(int))
	assert.Equal(t, 2, count)
}

func TestGroupReset(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("key")
	call := false
	g.Reset(func() interface{} {
		call = true
		return 1
	})

	length := 0
	for range g.vals {
		length++
	}

	assert.Equal(t, 0, length)

	g.Get("key")
	assert.Equal(t, true, call)
}

func TestGroupClear(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("key")
	length := 0
	for range g.vals {
		length++
	}
	assert.Equal(t, 1, length)

	g.Clear()
	length = 0
	for range g.vals {
		length++
	}
	assert.Equal(t, 0, length)
}
