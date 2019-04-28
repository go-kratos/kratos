package paladin_test

import (
	"io/ioutil"
	"os"
	"testing"

	"go-common/library/conf/paladin"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	// test data
	path := "/tmp/test_conf/"
	assert.Nil(t, os.MkdirAll(path, 0700))
	assert.Nil(t, ioutil.WriteFile(path+"test.toml", []byte(`
		text = "hello"	
		number = 100
		slice = [1, 2, 3]
		sliceStr = ["1", "2", "3"]
	`), 0644))
	// test client
	cli, err := paladin.NewFile(path + "test.toml")
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	// test map
	m := paladin.Map{}
	text, err := cli.Get("test.toml").String()
	assert.Nil(t, err)
	assert.Nil(t, m.Set(text), "text")
	s, err := m.Get("text").String()
	assert.Nil(t, err)
	assert.Equal(t, s, "hello", "text")
	n, err := m.Get("number").Int64()
	assert.Nil(t, err)
	assert.Equal(t, n, int64(100), "number")
}

func TestNewFilePath(t *testing.T) {
	// test data
	path := "/tmp/test_conf/"
	assert.Nil(t, os.MkdirAll(path, 0700))
	assert.Nil(t, ioutil.WriteFile(path+"test.toml", []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	assert.Nil(t, ioutil.WriteFile(path+"abc.toml", []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	// test client
	cli, err := paladin.NewFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	// test map
	m := paladin.Map{}
	text, err := cli.Get("test.toml").String()
	assert.Nil(t, err)
	assert.Nil(t, m.Set(text), "text")
	s, err := m.Get("text").String()
	assert.Nil(t, err, s)
	assert.Equal(t, s, "hello", "text")
	n, err := m.Get("number").Int64()
	assert.Nil(t, err, s)
	assert.Equal(t, n, int64(100), "number")
}
