package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"
)

const (
	_testJSON = `
{
    "test":{
        "settings":{
            "int_key":1000,
            "float_key":1000.1,
            "duration_key":10000,
            "string_key":"string_value"
        },
        "server":{
            "addr":"127.0.0.1",
            "port":8000
        }
    },
    "foo":[
        {
            "name":"nihao",
            "age":18
        },
        {
            "name":"nihao",
            "age":18
        }
    ]
}`

//	_testYaml = `
//Foo:
//    bar :
//        - {name: nihao,age: 1}
//        - {name: nihao,age: 1}
//
//
//`
)

//func TestScan(t *testing.T) {
//
//}

func TestFile(t *testing.T) {
	var (
		path = filepath.Join(os.TempDir(), "test_config")
		file = filepath.Join(path, "test.json")
		data = []byte(_testJSON)
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

func TestConfig(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test_config.json")
	defer os.Remove(path)
	if err := ioutil.WriteFile(path, []byte(_testJSON), 0666); err != nil {
		t.Error(err)
	}
	c := config.New(config.WithSource(
		NewSource(path),
	))
	testScan(t, c)

	testConfig(t, c)
}

func testConfig(t *testing.T, c config.Config) {
	var expected = map[string]interface{}{
		"test.settings.int_key":      int64(1000),
		"test.settings.float_key":    float64(1000.1),
		"test.settings.string_key":   "string_value",
		"test.settings.duration_key": time.Duration(10000),
		"test.server.addr":           "127.0.0.1",
		"test.server.port":           int64(8000),
	}
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	for key, value := range expected {
		switch value.(type) {
		case int64:
			if v, err := c.Value(key).Int(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case float64:
			if v, err := c.Value(key).Float(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case string:
			if v, err := c.Value(key).String(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case time.Duration:
			if v, err := c.Value(key).Duration(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		}
	}
	// scan
	var settings struct {
		IntKey      int64         `json:"int_key"`
		FloatKey    float64       `json:"float_key"`
		StringKey   string        `json:"string_key"`
		DurationKey time.Duration `json:"duration_key"`
	}
	if err := c.Value("test.settings").Scan(&settings); err != nil {
		t.Error(err)
	}
	if v := expected["test.settings.int_key"]; settings.IntKey != v {
		t.Errorf("no expect int_key value: %v, but got: %v", settings.IntKey, v)
	}
	if v := expected["test.settings.float_key"]; settings.FloatKey != v {
		t.Errorf("no expect float_key value: %v, but got: %v", settings.FloatKey, v)
	}
	if v := expected["test.settings.string_key"]; settings.StringKey != v {
		t.Errorf("no expect string_key value: %v, but got: %v", settings.StringKey, v)
	}
	if v := expected["test.settings.duration_key"]; settings.DurationKey != v {
		t.Errorf("no expect duration_key value: %v, but got: %v", settings.DurationKey, v)
	}

	// not found
	if _, err := c.Value("not_found_key").Bool(); errors.Is(err, config.ErrNotFound) {
		t.Logf("not_found_key not match: %v", err)
	}

}

func testScan(t *testing.T, c config.Config) {
	type TestJSON struct {
		Test struct {
			Settings struct {
				IntKey      int     `json:"int_key"`
				FloatKey    float64 `json:"float_key"`
				DurationKey int     `json:"duration_key"`
				StringKey   string  `json:"string_key"`
			} `json:"settings"`
			Server struct {
				Addr string `json:"addr"`
				Port int    `json:"port"`
			} `json:"server"`
		} `json:"test"`
		Foo []struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"foo"`
	}
	var conf TestJSON
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	if err := c.Scan(&conf); err != nil {
		t.Error(err)
	}
	t.Log(conf)
}
