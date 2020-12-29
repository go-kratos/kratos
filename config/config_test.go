package config

import (
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config/provider"
	"github.com/go-kratos/kratos/v2/config/provider/memory"
)

const (
	_yaml = `
test:
  settings:
    int_key: 100
    int64_key: 1000
    float_key: 1000.1
    duration: 10000
    string_key: string_value
  server:
    addr: 127.0.0.1
    port: 8000
`
	_toml = `
[test]
[test.settings]
  int_key = 100
  int64_key = 1000
  float_key = 1000.1
  duration = 10000
  string_key = "string_value"
[test.server]
  addr = '127.0.0.1'
  port = 8000
`
	_json = `
{
  "test": {
    "settings" : {
      "int_key": 100,
      "int64_key": 1000,
      "float_key": 1000.1, 
      "duration": 10000, 
      "string_key": "string_value"
    },
    "server": {
      "addr": "127.0.0.1",
      "port": 8000
    }
  }
}`
)

func TestConfigYAML(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "yaml",
			Key:    "test",
			Value:  []byte(strings.TrimSpace(_yaml)),
		})),
	)
	testConfig(t, c)
}

func TestConfigTOML(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "toml",
			Key:    "test",
			Value:  []byte(strings.TrimSpace(_toml)),
		})),
	)
	testConfig(t, c)
}

func TestConfigJSON(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "json",
			Key:    "test",
			Value:  []byte(strings.TrimSpace(_json)),
		})),
	)
	testConfig(t, c)
}

func testConfig(t *testing.T, c Config) {
	var expected = map[string]interface{}{
		"test.settings.int_key":    int(100),
		"test.settings.int64_key":  int64(1000),
		"test.settings.float_key":  float64(1000.1),
		"test.settings.duration":   time.Duration(10000),
		"test.settings.string_key": "string_value",
		"test.server.addr":         "127.0.0.1",
		"test.server.port":         int(8000),
	}
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	for key, value := range expected {
		switch value.(type) {
		case int:
			if v, err := c.Value(key).Int(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case int64:
			if v, err := c.Value(key).Int64(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case float64:
			if v, err := c.Value(key).Float64(); err != nil {
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
		IntKey    int     `json:"int_key"`
		FloatKey  float32 `json:"float_key"`
		StringKey string  `json:"string_key"`
	}
	if err := c.Value("test.settings").Scan(&settings); err != nil {
		t.Error(err)
	}
	t.Log(settings)

	// not found
	if _, err := c.Value("not_found_key").Bool(); err == nil {
		t.Logf("not_found_key not match: %v", err)
	}

}
