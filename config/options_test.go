package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestDefaultDecoder(t *testing.T) {
	src := &KeyValue{
		Key:    "service",
		Value:  []byte("config"),
		Format: "",
	}
	target := make(map[string]any)
	err := defaultDecoder(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(target, map[string]any{"service": []byte("config")}) {
		t.Fatal(`target is not equal to map[string]interface{}{"service": "config"}`)
	}

	src = &KeyValue{
		Key:    "service.name.alias",
		Value:  []byte("2233"),
		Format: "",
	}
	target = make(map[string]any)
	err = defaultDecoder(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(map[string]any{
		"service": map[string]any{
			"name": map[string]any{
				"alias": []byte("2233"),
			},
		},
	}, target) {
		t.Fatal(`target is not equal to map[string]interface{}{"service": map[string]interface{}{"name": map[string]interface{}{"alias": []byte("2233")}}}`)
	}
}

func TestDefaultResolver(t *testing.T) {
	var (
		portString = "8080"
		countInt   = 10
		rateFloat  = 0.9
	)

	data := map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:8081}",
				"count":    "${COUNT:0}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"url":      "${URL:http://example.com}",
				"array": []any{
					"${PORT}",
					map[string]any{"foobar": "${NOTEXIST:8081}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]any{
			"value": "foobar",
		},
		"PORT":   "8080",
		"COUNT":  "10",
		"ENABLE": "true",
		"RATE":   "0.9",
		"EMPTY":  "",
	}

	tests := []struct {
		name   string
		path   string
		expect any
	}{
		{
			name:   "test not exist int env with default",
			path:   "foo.bar.notexist",
			expect: 100,
		},
		{
			name:   "test string with default",
			path:   "foo.bar.port",
			expect: portString,
		},
		{
			name:   "test int with default",
			path:   "foo.bar.count",
			expect: countInt,
		},
		{
			name:   "test bool with default",
			path:   "foo.bar.enable",
			expect: true,
		},
		{
			name:   "test float without default",
			path:   "foo.bar.rate",
			expect: rateFloat,
		},
		{
			name:   "test empty value with default",
			path:   "foo.bar.empty",
			expect: "",
		},
		{
			name:   "test url with default",
			path:   "foo.bar.url",
			expect: "http://example.com",
		},
		{
			name:   "test array",
			path:   "foo.bar.array",
			expect: []any{portString, map[string]any{"foobar": "8081"}},
		},
		{
			name:   "test ${test.value}",
			path:   "foo.bar.value1",
			expect: "foobar",
		},
		{
			name:   "test $PORT",
			path:   "foo.bar.value2",
			expect: "$PORT",
		},
		{
			name:   "test abc${PORT}foo${COUNT}bar",
			path:   "foo.bar.value3",
			expect: "abc8080foo10bar",
		},
		{
			name:   "test ${foo${bar}}",
			path:   "foo.bar.value4",
			expect: "}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := defaultResolver(data)
			if err != nil {
				t.Fatal(err)
			}
			rd := reader{
				values: data,
			}
			if v, ok := rd.Value(test.path); ok {
				var actual any
				switch test.expect.(type) {
				case int:
					if actual, err = v.Int(); err == nil {
						if !reflect.DeepEqual(test.expect.(int), int(actual.(int64))) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case string:
					if actual, err = v.String(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				default:
					actual = v.Load()
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("expect: %#v, actual: %#v", test.expect, actual)
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

func TestNewDefaultResolver(t *testing.T) {
	var (
		portString = "8080"
		countInt   = 10
		rateFloat  = 0.9
	)

	data := map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:\"8081\"}",
				"count":    "${COUNT:\"0\"}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"url":      "${URL:\"http://example.com\"}",
				"array": []any{
					"${PORT}",
					map[string]any{"foobar": "${NOTEXIST:\"8081\"}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]any{
			"value": "foobar",
		},
		"PORT":   "\"8080\"",
		"COUNT":  "\"10\"",
		"ENABLE": "true",
		"RATE":   "0.9",
		"EMPTY":  "",
	}

	tests := []struct {
		name   string
		path   string
		expect any
	}{
		{
			name:   "test not exist int env with default",
			path:   "foo.bar.notexist",
			expect: 100,
		},
		{
			name:   "test string with default",
			path:   "foo.bar.port",
			expect: portString,
		},
		{
			name:   "test int with default",
			path:   "foo.bar.count",
			expect: countInt,
		},
		{
			name:   "test bool with default",
			path:   "foo.bar.enable",
			expect: true,
		},
		{
			name:   "test float without default",
			path:   "foo.bar.rate",
			expect: rateFloat,
		},
		{
			name:   "test empty value with default",
			path:   "foo.bar.empty",
			expect: "",
		},
		{
			name:   "test url with default",
			path:   "foo.bar.url",
			expect: "http://example.com",
		},
		{
			name:   "test array",
			path:   "foo.bar.array",
			expect: []any{portString, map[string]any{"foobar": "8081"}},
		},
		{
			name:   "test ${test.value}",
			path:   "foo.bar.value1",
			expect: "foobar",
		},
		{
			name:   "test $PORT",
			path:   "foo.bar.value2",
			expect: "$PORT",
		},
		//{
		//	name:   "test abc${PORT}foo${COUNT}bar",
		//	path:   "foo.bar.value3",
		//	expect: "abc8080foo10bar",
		//},
		{
			name:   "test ${foo${bar}}",
			path:   "foo.bar.value4",
			expect: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn := newActualTypesResolver(true)
			err := fn(data)
			if err != nil {
				t.Fatal(err)
			}
			rd := reader{
				values: data,
			}
			if v, ok := rd.Value(test.path); ok {
				var actual any
				switch test.expect.(type) {
				case int:
					if actual, err = v.Int(); err == nil {
						if !reflect.DeepEqual(test.expect.(int), int(actual.(int64))) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case string:
					if actual, err = v.String(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal("expect is not equal to actual")
						}
					}
				default:
					actual = v.Load()
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("expect: %#v, actual: %#v", test.expect, actual)
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

func TestExpand(t *testing.T) {
	tests := []struct {
		input   string
		mapping func(string) string
		want    string
	}{
		{
			input: "${a}",
			mapping: func(s string) string {
				return strings.ToUpper(s)
			},
			want: "A",
		},
		{
			input: "a",
			mapping: func(s string) string {
				return strings.ToUpper(s)
			},
			want: "a",
		},
	}
	for _, tt := range tests {
		if got := expand(tt.input, tt.mapping, false); got != tt.want {
			t.Errorf("expand() want: %s, got: %s", tt.want, got)
		}
	}
}

func TestExpandBoolPlaceholder(t *testing.T) {
	mapping := func(s string) string {
		switch s {
		case "port":
			return "8080"
		case "debug":
			return "true"
		case "disabled":
			return "false"
		}
		return ""
	}
	tests := []struct {
		name  string
		input string
		want  any
	}{
		{
			name:  "whole value is a placeholder resolving to true",
			input: "${debug}",
			want:  true,
		},
		{
			name:  "whole value is a placeholder resolving to false",
			input: "${disabled}",
			want:  false,
		},
		{
			name:  "whole value is a placeholder resolving to a non-bool string",
			input: "${port}",
			want:  "8080",
		},
		{
			name:  "placeholder embedded in a longer value stays a string",
			input: "port=${port}",
			want:  "port=8080",
		},
		{
			name:  "value without placeholder stays a string",
			input: "true",
			want:  "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expand(tt.input, mapping, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expand() want: %#v, got: %#v", tt.want, got)
			}
		})
	}
}

// TestDefaultResolverBoolPlaceholder guards the bool field case end to end: the
// default resolver leaves values as strings, which json and protojson then
// reject when the destination field is a bool.
func TestDefaultResolverBoolPlaceholder(t *testing.T) {
	data := map[string]any{
		"debug": "${DEBUG:false}",
		"port":  "${PORT:8080}",
		"name":  "${NAME:1234}",
	}
	if err := defaultResolver(data); err != nil {
		t.Fatal(err)
	}
	if got, ok := data["debug"].(bool); !ok || got {
		t.Errorf("debug want: false (bool), got: %#v", data["debug"])
	}
	// Values that are not booleans keep their string shape, so a string field
	// holding a numeric-looking value is not mistyped.
	if got, ok := data["port"].(string); !ok || got != "8080" {
		t.Errorf("port want: %q (string), got: %#v", "8080", data["port"])
	}
	if got, ok := data["name"].(string); !ok || got != "1234" {
		t.Errorf("name want: %q (string), got: %#v", "1234", data["name"])
	}
}

func TestWithMergeFunc(t *testing.T) {
	c := &options{}
	a := func(any, any) error {
		return nil
	}
	WithMergeFunc(a)(c)
	if c.merge == nil {
		t.Fatal("c.merge is nil")
	}
}
