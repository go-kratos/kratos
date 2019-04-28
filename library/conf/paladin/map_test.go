package paladin_test

import (
	"testing"

	"go-common/library/conf/paladin"

	"github.com/naoina/toml"
	"github.com/stretchr/testify/assert"
)

type fruit struct {
	Fruit []struct {
		Name string
	}
}

func (f *fruit) Set(text string) error {
	return toml.Unmarshal([]byte(text), f)
}

func TestMap(t *testing.T) {
	s := `
		# kv
		text = "hello"	
		number = 100
		point = 100.1
		boolean = true
		KeyCase = "test"

		# slice
		numbers = [1, 2, 3]
		strings = ["a", "b", "c"]
		empty = []
		[[fruit]]
		  name = "apple"
		[[fruit]]
		  name = "banana"

		# table 
		[database]
		server = "192.168.1.1"
		connection_max = 5000
		enabled = true

		[pool]
		[pool.breaker]
		xxx = "xxx"
		`
	m := paladin.Map{}
	assert.Nil(t, m.Set(s), s)
	str, err := m.Get("text").String()
	assert.Nil(t, err)
	assert.Equal(t, str, "hello", "text")
	n, err := m.Get("number").Int64()
	assert.Nil(t, err)
	assert.Equal(t, n, int64(100), "number")
	p, err := m.Get("point").Float64()
	assert.Nil(t, err)
	assert.Equal(t, p, 100.1, "point")
	b, err := m.Get("boolean").Bool()
	assert.Nil(t, err)
	assert.Equal(t, b, true, "boolean")
	// key lower case
	lb, err := m.Get("Boolean").Bool()
	assert.Nil(t, err)
	assert.Equal(t, lb, true, "boolean")
	lt, err := m.Get("KeyCase").String()
	assert.Nil(t, err)
	assert.Equal(t, lt, "test", "key case")
	var sliceInt []int64
	err = m.Get("numbers").Slice(&sliceInt)
	assert.Nil(t, err)
	assert.Equal(t, sliceInt, []int64{1, 2, 3})
	var sliceStr []string
	err = m.Get("strings").Slice(&sliceStr)
	assert.Nil(t, err)
	assert.Equal(t, sliceStr, []string{"a", "b", "c"})
	err = m.Get("strings").Slice(&sliceStr)
	assert.Nil(t, err)
	assert.Equal(t, sliceStr, []string{"a", "b", "c"})
	// errors
	err = m.Get("strings").Slice(sliceInt)
	assert.NotNil(t, err)
	err = m.Get("strings").Slice(&sliceInt)
	assert.NotNil(t, err)
	var obj struct {
		Name string
	}
	err = m.Get("strings").Slice(obj)
	assert.NotNil(t, err)
	err = m.Get("strings").Slice(&obj)
	assert.NotNil(t, err)
}
