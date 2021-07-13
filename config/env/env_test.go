package env

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/stretchr/testify/assert"
)

const _testJSON = `
{
    "test":{
        "server":{
			"name":"$SERVICE_NAME",
            "addr":"${ADDR:127.0.0.1}",
            "port":"${PORT:8080}"
        }
    },
    "foo":[
        {
            "name":"Tom",
            "age":"${AGE}"
        }
    ]
}`

func TestEnvWithPrefix(t *testing.T) {
	var (
		path     = filepath.Join(os.TempDir(), "test_config")
		filename = filepath.Join(path, "test.json")
		data     = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		t.Error(err)
	}

	// set env
	prefix1, prefix2 := "KRATOS_", "FOO"
	envs := map[string]string{
		prefix1 + "SERVICE_NAME": "kratos_app",
		prefix2 + "ADDR":         "192.168.0.1",
		prefix1 + "AGE":          "20",
	}

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(err)
		}
	}

	c := config.New(config.WithSource(
		NewSource(prefix1, prefix2),
		file.NewSource(path),
	))

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "test $KEY",
			path:   "test.server.name",
			expect: "kratos_app",
		},
		{
			name:   "test ${KEY:DEFAULT} without default",
			path:   "test.server.addr",
			expect: "192.168.0.1",
		},
		{
			name:   "test ${KEY:DEFAULT} with default",
			path:   "test.server.port",
			expect: "8080",
		},
		{
			name: "test ${KEY} in array",
			path: "foo",
			expect: []interface{}{
				map[string]interface{}{
					"name": "Tom",
					"age":  "20",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			v := c.Value(test.path)
			if v.Load() != nil {
				var actual interface{}
				switch test.expect.(type) {
				case int:
					if actual, err = v.Int(); err == nil {
						assert.Equal(t, test.expect, int(actual.(int64)), "int value should be equal")
					}
				case string:
					if actual, err = v.String(); err == nil {
						assert.Equal(t, test.expect, actual, "string value should be equal")
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						assert.Equal(t, test.expect, actual, "bool value should be equal")
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						assert.Equal(t, test.expect, actual, "float64 value should be equal")
					}
				default:
					actual = v.Load()
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("\nexpect: %#v\nactural: %#v", test.expect, actual)
						t.Fail()
					}
				}
				if err != nil {
					t.Error(err)
				}
			} else {
				t.Error("value path not found")
			}
		})
	}
}

func TestEnvWithoutPrefix(t *testing.T) {
	var (
		path     = filepath.Join(os.TempDir(), "test_config")
		filename = filepath.Join(path, "test.json")
		data     = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		t.Error(err)
	}

	// set env
	envs := map[string]string{
		"SERVICE_NAME": "kratos_app",
		"ADDR":         "192.168.0.1",
		"AGE":          "20",
	}

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(err)
		}
	}

	c := config.New(config.WithSource(
		NewSource(),
		file.NewSource(path),
	))

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "test $KEY",
			path:   "test.server.name",
			expect: "kratos_app",
		},
		{
			name:   "test ${KEY:DEFAULT} without default",
			path:   "test.server.addr",
			expect: "192.168.0.1",
		},
		{
			name:   "test ${KEY:DEFAULT} with default",
			path:   "test.server.port",
			expect: "8080",
		},
		{
			name: "test ${KEY} in array",
			path: "foo",
			expect: []interface{}{
				map[string]interface{}{
					"name": "Tom",
					"age":  "20",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			v := c.Value(test.path)
			if v.Load() != nil {
				var actual interface{}
				switch test.expect.(type) {
				case int:
					if actual, err = v.Int(); err == nil {
						assert.Equal(t, test.expect, int(actual.(int64)), "int value should be equal")
					}
				case string:
					if actual, err = v.String(); err == nil {
						assert.Equal(t, test.expect, actual, "string value should be equal")
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						assert.Equal(t, test.expect, actual, "bool value should be equal")
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						assert.Equal(t, test.expect, actual, "float64 value should be equal")
					}
				default:
					actual = v.Load()
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("\nexpect: %#v\nactural: %#v", test.expect, actual)
						t.Fail()
					}
				}
				if err != nil {
					t.Error(err)
				}
			} else {
				t.Error("value path not found")
			}
		})
	}
}
