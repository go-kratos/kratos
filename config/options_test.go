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
	target := make(map[string]interface{})
	err := defaultDecoder(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(target, map[string]interface{}{"service": []byte("config")}) {
		t.Fatal(`target is not equal to map[string]interface{}{"service": "config"}`)
	}

	src = &KeyValue{
		Key:    "service.name.alias",
		Value:  []byte("2233"),
		Format: "",
	}
	target = make(map[string]interface{})
	err = defaultDecoder(src, target)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(map[string]interface{}{
		"service": map[string]interface{}{
			"name": map[string]interface{}{
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

	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:8081}",
				"count":    "${COUNT:0}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"url":      "${URL:http://example.com}",
				"array": []interface{}{
					"${PORT}",
					map[string]interface{}{"foobar": "${NOTEXIST:8081}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]interface{}{
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
			expect: []interface{}{portString, map[string]interface{}{"foobar": "8081"}},
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
				var actual interface{}
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

	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"notexist": "${NOTEXIST:100}",
				"port":     "${PORT:\"8081\"}",
				"count":    "${COUNT:\"0\"}",
				"enable":   "${ENABLE:false}",
				"rate":     "${RATE}",
				"empty":    "${EMPTY:foobar}",
				"url":      "${URL:\"http://example.com\"}",
				"array": []interface{}{
					"${PORT}",
					map[string]interface{}{"foobar": "${NOTEXIST:\"8081\"}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]interface{}{
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
			expect: []interface{}{portString, map[string]interface{}{"foobar": "8081"}},
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
				var actual interface{}
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
