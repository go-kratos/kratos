package config

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	testcomplex "github.com/go-kratos/kratos/v3/internal/testdata/complex"
)

// bootstrap mirrors the struct used in the issue repro (go-kratos/kratos#3881).
type bootstrapStruct struct {
	Debug   bool   `json:"debug"`
	Port    int    `json:"port"`
	Level   string `json:"level"`
	Version string `json:"version"`
}

// envLikeSource emits format-less KeyValues exactly like the env source does.
type envLikeSource struct {
	keys map[string]string
}

func (e *envLikeSource) Load() ([]*KeyValue, error) {
	var kvs []*KeyValue
	for k, v := range e.keys {
		kvs = append(kvs, &KeyValue{Key: k, Value: []byte(v)})
	}
	return kvs, nil
}

func (e *envLikeSource) Watch() (Watcher, error) {
	return newTestWatcher(make(chan struct{}), make(chan struct{})), nil
}

// TestScan_EnvBoolOverride reproduces the issue: a bool field overridden from
// the environment must be accepted by Scan instead of failing with
// "json: cannot unmarshal string into Go struct field ... of type bool".
func TestScan_EnvBoolOverride(t *testing.T) {
	fileSrc := newTestJSONSource(`{
		"debug": false,
		"port": 8000,
		"level": "info",
		"version": "1234"
	}`)
	envSrc := &envLikeSource{keys: map[string]string{
		"debug":   "true",
		"level":   "debug",
		"version": "1234",
	}}

	c := New(WithSource(fileSrc, envSrc))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	var bc bootstrapStruct
	if err := c.Scan(&bc); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if !bc.Debug {
		t.Errorf("expect Debug=true, got %v", bc.Debug)
	}
	if bc.Port != 8000 {
		t.Errorf("expect Port=8000, got %v", bc.Port)
	}
	if bc.Level != "debug" {
		t.Errorf("expect Level=debug, got %q", bc.Level)
	}
	// numeric-looking string must stay a string, not become a number.
	if bc.Version != "1234" {
		t.Errorf("expect Version=\"1234\" (string), got %T(%v)", bc.Version, bc.Version)
	}
}

// TestScan_EnvIntOverride checks int fields are coerced from env strings too.
func TestScan_EnvIntOverride(t *testing.T) {
	fileSrc := newTestJSONSource(`{"port": 8000}`)
	envSrc := &envLikeSource{keys: map[string]string{"port": "9000"}}

	c := New(WithSource(fileSrc, envSrc))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	var bc struct {
		Port int `json:"port"`
	}
	if err := c.Scan(&bc); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if bc.Port != 9000 {
		t.Errorf("expect Port=9000, got %v", bc.Port)
	}
}

// TestScan_EnvNestedStruct checks coercion inside nested structs.
func TestScan_EnvNestedStruct(t *testing.T) {
	fileSrc := newTestJSONSource(`{
		"server": {
			"http": {"addr": "0.0.0.0", "port": 8080, "enable_ssl": false, "timeout": 0.5}
		}
	}`)
	envSrc := &envLikeSource{keys: map[string]string{
		"server.http.enable_ssl": "true",
		"server.http.timeout":    "1.5",
	}}

	c := New(WithSource(fileSrc, envSrc))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	var sc struct {
		Server struct {
			HTTP struct {
				Addr      string  `json:"addr"`
				Port      int     `json:"port"`
				EnableSSL bool    `json:"enable_ssl"`
				Timeout   float64 `json:"timeout"`
			} `json:"http"`
		} `json:"server"`
	}
	if err := c.Scan(&sc); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if !sc.Server.HTTP.EnableSSL {
		t.Error("expect EnableSSL=true")
	}
	if sc.Server.HTTP.Timeout != 1.5 {
		t.Errorf("expect Timeout=1.5, got %v", sc.Server.HTTP.Timeout)
	}
	if sc.Server.HTTP.Port != 8080 {
		t.Errorf("expect Port=8080, got %v", sc.Server.HTTP.Port)
	}
}

// TestScan_ProtoBoolOverride reproduces the proto path from the issue:
// protojson accepts a quoted string for every scalar kind except bool.
func TestScan_ProtoBoolOverride(t *testing.T) {
	fileSrc := newTestJSONSource(`{"b": false, "age": 30, "count": 100}`)
	envSrc := &envLikeSource{keys: map[string]string{
		"b":     "true",
		"age":   "42",
		"count": "7",
	}}

	c := New(WithSource(fileSrc, envSrc))
	defer c.Close()
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}
	var msg testcomplex.Complex
	if err := c.Scan(&msg); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}
	if !msg.GetB() {
		t.Error("expect B=true")
	}
	if msg.GetAge() != 42 {
		t.Errorf("expect Age=42, got %d", msg.GetAge())
	}
	if msg.GetCount() != 7 {
		t.Errorf("expect Count=7, got %d", msg.GetCount())
	}
}

// TestConvertStringLeaf covers the no-text-guessing contract: numeric-looking
// strings targeting string fields must never be converted.
func TestConvertStringLeaf(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		typ    reflect.Type
		want   any
		wantOk bool
	}{
		{name: "bool true", value: "true", typ: reflect.TypeOf(false), want: true, wantOk: true},
		{name: "bool false", value: "false", typ: reflect.TypeOf(false), want: false, wantOk: true},
		{name: "bool invalid", value: "yes", typ: reflect.TypeOf(false), want: "yes", wantOk: false},
		{name: "int", value: "42", typ: reflect.TypeOf(int(0)), want: int64(42), wantOk: true},
		{name: "int invalid", value: "abc", typ: reflect.TypeOf(int(0)), want: "abc", wantOk: false},
		{name: "uint32", value: "7", typ: reflect.TypeOf(uint32(0)), want: uint64(7), wantOk: true},
		{name: "float", value: "1.5", typ: reflect.TypeOf(0.0), want: 1.5, wantOk: true},
		{name: "string numeric stays", value: "1234", typ: reflect.TypeOf(""), want: "1234", wantOk: false},
		{name: "string normal stays", value: "info", typ: reflect.TypeOf(""), want: "info", wantOk: false},
		{name: "non-string leaf", value: 42, typ: reflect.TypeOf(int(0)), want: 42, wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := convertStringLeaf(tt.value, tt.typ)
			if ok != tt.wantOk {
				t.Fatalf("convertStringLeaf(%v, %v) ok = %v, want %v", tt.value, tt.typ, ok, tt.wantOk)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertStringLeaf(%v, %v) = %v, want %v", tt.value, tt.typ, got, tt.want)
			}
		})
	}
}

// TestCoerceScalarTypes_NoChange verifies that a config with no string leaves
// needing conversion round-trips byte-identically (no spurious rewrite).
func TestCoerceScalarTypes_NoChange(t *testing.T) {
	data := []byte(`{"debug":false,"port":8000,"level":"info"}`)
	got, err := coerceScalarTypes(data, &bootstrapStruct{})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(data) {
		t.Errorf("no-change config was rewritten: got %s", got)
	}
}

// TestCoerceScalarTypes_ProtoSnakeNames checks protojson accepts the proto field
// name (snake_case or json_name) when the source emits dotted env keys.
func TestCoerceScalarTypes_ProtoSnakeNames(t *testing.T) {
	var msg testcomplex.Complex
	data := []byte(`{"no_one":"x","b":"true"}`)
	got, err := coerceScalarTypes(data, &msg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), `"b":true`) {
		t.Errorf("expect b coerced to true, got %s", got)
	}
	if strings.Contains(string(got), `"no_one"`) == false {
		t.Errorf("string field no_one should be untouched, got %s", got)
	}
	_ = json.Valid(got)
}
