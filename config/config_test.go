package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultResolver(t *testing.T) {
	var (
		envString = "8080"
		envInt    = 10
		envBool   = true
		envFloat  = 0.9
	)
	var err error
	err = os.Setenv("PORT", envString)
	assert.NoError(t, err)
	err = os.Setenv("COUNT", strconv.Itoa(envInt))
	assert.NoError(t, err)
	err = os.Setenv("ENABLE", strconv.FormatBool(envBool))
	assert.NoError(t, err)
	err = os.Setenv("RATE", fmt.Sprintf("%v", envFloat))
	assert.NoError(t, err)
	err = os.Setenv("EMPTY", "")
	assert.NoError(t, err)

	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:8081}",
				"count":    "${COUNT:0}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"array":    []interface{}{"${PORT}", "${NOTEXIST:8081}"},
			},
		},
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "test not exist int env with default",
			path:   "foo.bar.notexist",
			expect: 100,
		},
		{
			name:   "test string with default",
			path:   "foo.bar.port",
			expect: envString,
		},
		{
			name:   "test int with default",
			path:   "foo.bar.count",
			expect: envInt,
		},
		{
			name:   "test bool with default",
			path:   "foo.bar.enable",
			expect: envBool,
		},
		{
			name:   "test float with default",
			path:   "foo.bar.rate",
			expect: envFloat,
		},
		{
			name:   "test empty value with default",
			path:   "foo.bar.empty",
			expect: "",
		},
		//TODO: add array test case

		// // can not support xml config template because
		// // we can not decode xml to map[string]interface{}
		//{
		//	name: "xml",
		//	path: "http.server.port",
		//	expect: envString,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := defaultResolver(data)
			assert.NoError(t, err)
			rd := reader{
				values: data,
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
