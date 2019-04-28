package paladin_test

import (
	"testing"

	"go-common/library/conf/paladin"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	cs := map[string]string{
		"key_toml": `
			key_bool = true
			key_int = 100
			key_float = 100.1
			key_string = "text"	
		`,
	}
	cli := paladin.NewMock(cs)
	// test vlaue
	var m paladin.TOML
	err := cli.Get("key_toml").Unmarshal(&m)
	assert.Nil(t, err)
	b, err := m.Get("key_bool").Bool()
	assert.Nil(t, err)
	assert.Equal(t, b, true)
	i, err := m.Get("key_int").Int64()
	assert.Nil(t, err)
	assert.Equal(t, i, int64(100))
	f, err := m.Get("key_float").Float64()
	assert.Nil(t, err)
	assert.Equal(t, f, float64(100.1))
	s, err := m.Get("key_string").String()
	assert.Nil(t, err)
	assert.Equal(t, s, "text")
}
