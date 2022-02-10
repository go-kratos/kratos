package env

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

const _testJSON = `
{
    "test":{
        "server":{
			"name":"${SERVICE_NAME}",
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
		path     = filepath.Join(t.TempDir(), "test_config")
		filename = filepath.Join(path, "test.json")
		data     = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(filename, data, 0o666); err != nil {
		t.Error(err)
	}

	// set env
	prefix1, prefix2 := "KRATOS_", "FOO"
	envs := map[string]string{
		prefix1 + "SERVICE_NAME": "kratos_app",
		prefix2 + "ADDR":         "192.168.0.1",
		prefix1 + "AGE":          "20",
		// only prefix
		prefix2:       "foo",
		prefix2 + "_": "foo_",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}

	c := config.New(config.WithSource(
		file.NewSource(path),
		NewSource(prefix1, prefix2),
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
						if !reflect.DeepEqual(test.expect.(int), int(actual.(int64))) {
							t.Errorf("expect %v, actual %v", test.expect, actual)
						}
					}
				case string:
					if actual, err = v.String(); err == nil {
						if !reflect.DeepEqual(test.expect.(string), actual.(string)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						if !reflect.DeepEqual(test.expect.(bool), actual.(bool)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						if !reflect.DeepEqual(test.expect.(float64), actual.(float64)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
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
		path     = filepath.Join(t.TempDir(), "test_config")
		filename = filepath.Join(path, "test.json")
		data     = []byte(_testJSON)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(filename, data, 0o666); err != nil {
		t.Error(err)
	}

	// set env
	envs := map[string]string{
		"SERVICE_NAME": "kratos_app",
		"ADDR":         "192.168.0.1",
		"AGE":          "20",
	}

	for k, v := range envs {
		os.Setenv(k, v)
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
						if !reflect.DeepEqual(test.expect.(int), int(actual.(int64))) {
							t.Errorf("expect %v, actual %v", test.expect, actual)
						}
					}
				case string:
					if actual, err = v.String(); err == nil {
						if !reflect.DeepEqual(test.expect.(string), actual.(string)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						if !reflect.DeepEqual(test.expect.(bool), actual.(bool)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						if !reflect.DeepEqual(test.expect.(float64), actual.(float64)) {
							t.Errorf(`expect %v, actual %v`, test.expect, actual)
						}
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

func Test_env_load(t *testing.T) {
	type fields struct {
		prefixs []string
	}
	type args struct {
		envStrings []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*config.KeyValue
	}{
		{
			name: "without prefixes",
			fields: fields{
				prefixs: nil,
			},
			args: args{
				envStrings: []string{
					"SERVICE_NAME=kratos_app",
					"ADDR=192.168.0.1",
					"AGE=20",
				},
			},
			want: []*config.KeyValue{
				{Key: "SERVICE_NAME", Value: []byte("kratos_app"), Format: ""},
				{Key: "ADDR", Value: []byte("192.168.0.1"), Format: ""},
				{Key: "AGE", Value: []byte("20"), Format: ""},
			},
		},

		{
			name: "empty prefix",
			fields: fields{
				prefixs: []string{""},
			},
			args: args{
				envStrings: []string{
					"__SERVICE_NAME=kratos_app",
					"__ADDR=192.168.0.1",
					"__AGE=20",
				},
			},
			want: []*config.KeyValue{
				{Key: "_SERVICE_NAME", Value: []byte("kratos_app"), Format: ""},
				{Key: "_ADDR", Value: []byte("192.168.0.1"), Format: ""},
				{Key: "_AGE", Value: []byte("20"), Format: ""},
			},
		},

		{
			name: "underscore prefix",
			fields: fields{
				prefixs: []string{"_"},
			},
			args: args{
				envStrings: []string{
					"__SERVICE_NAME=kratos_app",
					"__ADDR=192.168.0.1",
					"__AGE=20",
				},
			},
			want: []*config.KeyValue{
				{Key: "SERVICE_NAME", Value: []byte("kratos_app"), Format: ""},
				{Key: "ADDR", Value: []byte("192.168.0.1"), Format: ""},
				{Key: "AGE", Value: []byte("20"), Format: ""},
			},
		},

		{
			name: "with prefixes",
			fields: fields{
				prefixs: []string{"KRATOS_", "FOO"},
			},
			args: args{
				envStrings: []string{
					"KRATOS_SERVICE_NAME=kratos_app",
					"KRATOS_ADDR=192.168.0.1",
					"FOO_AGE=20",
				},
			},
			want: []*config.KeyValue{
				{Key: "SERVICE_NAME", Value: []byte("kratos_app"), Format: ""},
				{Key: "ADDR", Value: []byte("192.168.0.1"), Format: ""},
				{Key: "AGE", Value: []byte("20"), Format: ""},
			},
		},

		{
			name: "should not panic #1",
			fields: fields{
				prefixs: []string{"FOO"},
			},
			args: args{
				envStrings: []string{
					"FOO=123",
				},
			},
			want: nil,
		},

		{
			name: "should not panic #2",
			fields: fields{
				prefixs: []string{"FOO=1"},
			},
			args: args{
				envStrings: []string{
					"FOO=123",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &env{
				prefixs: tt.fields.prefixs,
			}
			got := e.load(tt.args.envStrings)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("env.load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchPrefix(t *testing.T) {
	type args struct {
		prefixes []string
		s        string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		wantOk bool
	}{
		{args: args{prefixes: nil, s: "foo=123"}, want: "", wantOk: false},
		{args: args{prefixes: []string{""}, s: "foo=123"}, want: "", wantOk: true},
		{args: args{prefixes: []string{"foo"}, s: "foo=123"}, want: "foo", wantOk: true},
		{args: args{prefixes: []string{"foo=1"}, s: "foo=123"}, want: "foo=1", wantOk: true},
		{args: args{prefixes: []string{"foo=1234"}, s: "foo=123"}, want: "", wantOk: false},
		{args: args{prefixes: []string{"bar"}, s: "foo=123"}, want: "", wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := matchPrefix(tt.args.prefixes, tt.args.s)
			if got != tt.want {
				t.Errorf("matchPrefix() got = %v, want %v", got, tt.want)
			}
			if gotOk != tt.wantOk {
				t.Errorf("matchPrefix() gotOk = %v, wantOk %v", gotOk, tt.wantOk)
			}
		})
	}
}
