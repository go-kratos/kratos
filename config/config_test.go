package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/stretchr/testify/assert"
)

func TestFillTemplate(t *testing.T) {
	var (
		envString = "8080"
		envInt    = 10
		envBool   = true
		envFloat  = 0.9
	)

	if err := os.Setenv("PORT", envString); err != nil {
		t.Fatal("set env failed")
	}
	if err := os.Setenv("COUNT", strconv.Itoa(envInt)); err != nil {
		t.Fatal("set env failed")
	}
	if err := os.Setenv("ENABLE", strconv.FormatBool(envBool)); err != nil {
		t.Fatal("set env failed")
	}
	if err := os.Setenv("RATE", fmt.Sprintf("%v", envFloat)); err != nil {
		t.Fatal("set env failed")
	}

	tests := []struct {
		name   string
		format string
		data   string
		path   string
		expect interface{}
	}{
		{
			name:   "json string",
			format: "json",
			data:   `{"http": {"server": {"port": {{ env "PORT" }}}}}`,
			path:   "http.server.port",
			expect: envString,
		},
		{
			name:   "yaml string",
			format: "yaml",
			data: `
http:
  server:
    port: {{ env "PORT" }}`,
			path:   "http.server.port",
			expect: envString,
		},
		{
			name:   "json int",
			format: "json",
			data:   `{"a": {"b": {"c": {{ env "COUNT" }}}}}`,
			path:   "a.b.c",
			expect: envInt,
		},
		{
			name:   "yaml int",
			format: "yaml",
			data: `
a:
  b:
    c: {{ env "COUNT" }}`,
			path:   "a.b.c",
			expect: envInt,
		},
		{
			name:   "json bool",
			format: "json",
			data:   `{"a": {"b": {"c": {{ env "ENABLE" }}}}}`,
			path:   "a.b.c",
			expect: envBool,
		},
		{
			name:   "yaml bool",
			format: "yaml",
			data: `
a:
  b:
    c: {{ env "ENABLE" }}`,
			path:   "a.b.c",
			expect: envBool,
		},
		{
			name:   "json float",
			format: "json",
			data:   `{"a": {"b": {"c": {{ env "RATE" }}}}}`,
			path:   "a.b.c",
			expect: envFloat,
		},
		{
			name:   "yaml float",
			format: "yaml",
			data: `
a:
  b:
    c: {{ env "RATE" }}`,
			path:   "a.b.c",
			expect: envFloat,
		},
		// can not support xml config template because we can not unmarshal xml to map[string]interface{}
		//{
		//	name: "xml",
		//	format: "xml",
		//	data: `<http><server><port>{{ env "PORT" }}</port></server></http>`,
		//	path: "http.server.port",
		//	expect: envString,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := fillTemplate([]byte(test.data))
			if err != nil {
				t.Fatal(err)
			}
			v := make(map[string]interface{})
			codec := encoding.GetCodec(test.format)
			if codec == nil {
				t.Fatalf("unsupported format: %s", test.format)
			}
			t.Log(string(d))
			if err := codec.Unmarshal(d, &v); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}

			rd := reader{
				values: v,
			}
			if v, ok := rd.Value(test.path); ok {
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
					assert.Equal(t, test.expect, actual, "interface{} value should be equal")
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
