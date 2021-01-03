package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(os.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
		data = []byte(`{"key":"value"}`)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(file, data, 0666); err != nil {
		t.Error(err)
	}
	testSource(t, file, data)
	testSource(t, path, data)
}

func testSource(t *testing.T, path string, data []byte) {
	t.Log(path)

	s := NewSource(path)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	if string(kvs[0].Value) != string(data) {
		t.Errorf("no expected: %s, but got: %s", kvs[0].Value, data)
	}
}
