package paladin

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	// test data
	path := "/tmp/test_conf/"
	assert.Nil(t, os.MkdirAll(path, 0700))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`
		text = "hello"	
		number = 100
		slice = [1, 2, 3]
		sliceStr = ["1", "2", "3"]
	`), 0644))
	// test client
	cli, err := NewFile(filepath.Join(path, "test.toml"))
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	// test map
	m := Map{}
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
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "abc.toml"), []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	// test client
	cli, err := NewFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	// test map
	m := Map{}
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

func TestFileEvent(t *testing.T) {
	// test data
	path := "/tmp/test_conf_event/"
	assert.Nil(t, os.MkdirAll(path, 0700))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "abc.toml"), []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	// test client
	cli, err := NewFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	ch := cli.WatchEvent(context.Background(), "test.toml", "abc.toml")
	time.Sleep(time.Millisecond)
	timeout := time.NewTimer(time.Second)

	// for file test.toml
	ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`hello`), 0644)
	select {
	case <-timeout.C:
		t.Fatalf("run test timeout")
	case ev := <-ch:
		if ev.Key == "test.toml" {
			assert.Equal(t, EventUpdate, ev.Event)
			assert.Equal(t, "hello", ev.Value)
		}
	}
	content1, _ := cli.Get("test.toml").String()
	assert.Equal(t, "hello", content1)

	// for file abc.toml
	ioutil.WriteFile(filepath.Join(path, "abc.toml"), []byte(`test`), 0644)
	select {
	case <-timeout.C:
		t.Fatalf("run test timeout")
	case ev := <-ch:
		if ev.Key == "abc.toml" {
			assert.Equal(t, EventUpdate, ev.Event)
			assert.Equal(t, "test", ev.Value)
		}
	}
	content2, _ := cli.Get("abc.toml").String()
	assert.Equal(t, "test", content2)
}

func TestHiddenFile(t *testing.T) {
	path := "/tmp/test_hidden_event/"
	assert.Nil(t, os.MkdirAll(path, 0700))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`hello`), 0644))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "abc.toml"), []byte(`
		text = "hello"	
		number = 100
	`), 0644))
	// test client
	cli, err := NewFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	cli.WatchEvent(context.Background(), "test.toml")
	time.Sleep(time.Millisecond)
	ioutil.WriteFile(filepath.Join(path, "abc.toml"), []byte(`hello`), 0644)
	time.Sleep(time.Second)
	content1, _ := cli.Get("test.toml").String()
	assert.Equal(t, "hello", content1)
	_, err = cli.Get(".abc.toml").String()
	assert.NotNil(t, err)
}

func TestOneLevelSymbolicFile(t *testing.T) {
	path := "/tmp/test_symbolic_link/"
	path2 := "/tmp/test_symbolic_link/configs/"
	assert.Nil(t, os.MkdirAll(path2, 0700))
	assert.Nil(t, ioutil.WriteFile(filepath.Join(path, "test.toml"), []byte(`hello`), 0644))
	assert.Nil(t, os.Symlink(filepath.Join(path, "test.toml"), filepath.Join(path2, "test.toml.ln")))
	// test client
	cli, err := NewFile(path2)
	assert.Nil(t, err)
	assert.NotNil(t, cli)
	content, _ := cli.Get("test.toml.ln").String()
	assert.Equal(t, "hello", content)
	os.Remove(filepath.Join(path, "test.toml"))
	os.Remove(filepath.Join(path2, "test.toml.ln"))
}
