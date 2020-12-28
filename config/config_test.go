package config

import (
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/config/provider"
	"github.com/go-kratos/kratos/v2/config/provider/memory"
)

func TestConfigYAML(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "yaml",
			Key:    "test",
			Value: []byte(strings.TrimSpace(`
test:
  settings:
    int_key: 100
    float_key: 1000.1
    string_key: string_value
  server:
    addr: 127.0.0.1
    port: 8000`)),
		})),
	)
	testConfig(t, c)
}

func TestConfigJSON(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "json",
			Key:    "test",
			Value: []byte(`
		{
			"test": {
				"settings" : {
					"int_key": 100,
					"float_key": 1000.1, 
					"string_key": "string_value"
				},
				"server": {
					"addr": "127.0.0.1",
					"port": 8000
				}
			}
		}
			`),
		})),
	)
	testConfig(t, c)
}

func TestConfigTOML(t *testing.T) {
	c := New(WithProvider(
		memory.New(nil, provider.KeyValue{
			Format: "toml",
			Key:    "test",
			Value: []byte(strings.TrimSpace(`
[test]
[test.settings]
int_key = 100
float_key = 1000.1
string_key = "string_value"
[test.server]
addr = '127.0.0.1'
port = 8000
			`)),
		})),
	)
	testConfig(t, c)
}

func testConfig(t *testing.T, c Config) {
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	if v, err := c.Value("test.settings.int_key").Int(); err != nil {
		t.Error(err)
	} else {
		t.Logf("int_key: %d", v)
	}
	if v, err := c.Value("test.settings.float_key").Float64(); err != nil {
		t.Error(err)
	} else {
		t.Logf("float_key: %f", v)
	}
	if v, err := c.Value("test.settings.string_key").String(); err != nil {
		t.Error(err)
	} else {
		t.Logf("string_key: %s", v)
	}
	if v, err := c.Value("test.server.addr").String(); err != nil {
		t.Error(err)
	} else {
		t.Logf("server.addr: %s", v)
	}
	if v, err := c.Value("test.server.port").Int(); err != nil {
		t.Error(err)
	} else {
		t.Logf("server.port: %d", v)
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
