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
	v := g.Get("/x/internal/dummy/user")
	assert.Equal(t, 1, v.(int))

	v = g.Get("/x/internal/dummy/avatar")
	assert.Equal(t, 2, v.(int))

	v = g.Get("/x/internal/dummy/user")
	assert.Equal(t, 1, v.(int))
	assert.Equal(t, 2, count)

}

func TestGroupReset(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("/x/internal/dummy/user")
	call := false
	g.Reset(func() interface{} {
		call = true
		return 1
	})

	length := 0
	g.objs.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)

	g.Get("/x/internal/dummy/user")
	assert.Equal(t, true, call)
}

func TestGroupClear(t *testing.T) {
	g := NewGroup(func() interface{} {
		return 1
	})
	g.Get("/x/internal/dummy/user")
	length := 0
	g.objs.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 1, length)

	g.Clear()
	length = 0
	g.objs.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	assert.Equal(t, 0, length)

}
