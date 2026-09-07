package config

import (
	"errors"
	"testing"
)

const (
	_testJSON = `
{
    "server":{
        "http":{
            "addr":"0.0.0.0",
			"port":80,
            "timeout":0.5,
			"enable_ssl":true
        },
        "grpc":{
            "addr":"0.0.0.0",
			"port":10080,
            "timeout":0.2
        }
    },
    "data":{
        "database":{
            "driver":"mysql",
            "source":"root:root@tcp(127.0.0.1:3306)/karta_id?parseTime=true"
        }
    },
	"endpoints":[
		"www.aaa.com",
		"www.bbb.org"
	]
}`
)

type testConfigStruct struct {
	Server struct {
		HTTP struct {
			Addr      string  `json:"addr"`
			Port      int     `json:"port"`
			Timeout   float64 `json:"timeout"`
			EnableSSL bool    `json:"enable_ssl"`
		} `json:"http"`
		GRPC struct {
			Addr    string  `json:"addr"`
			Port    int     `json:"port"`
			Timeout float64 `json:"timeout"`
		} `json:"grpc"`
	} `json:"server"`
	Data struct {
		Database struct {
			Driver string `json:"driver"`
			Source string `json:"source"`
		} `json:"database"`
	} `json:"data"`
	Endpoints []string `json:"endpoints"`
}

type testJSONSource struct {
	data string
	sig  chan struct{}
	err  chan struct{}
}

func newTestJSONSource(data string) *testJSONSource {
	return &testJSONSource{data: data, sig: make(chan struct{}), err: make(chan struct{})}
}

func (p *testJSONSource) Load() ([]*KeyValue, error) {
	kv := &KeyValue{
		Key:    "json",
		Value:  []byte(p.data),
		Format: "json",
	}
	return []*KeyValue{kv}, nil
}

func (p *testJSONSource) Watch() (Watcher, error) {
	return newTestWatcher(p.sig, p.err), nil
}

type testWatcher struct {
	sig  chan struct{}
	err  chan struct{}
	exit chan struct{}
}

func newTestWatcher(sig, err chan struct{}) Watcher {
	return &testWatcher{sig: sig, err: err, exit: make(chan struct{})}
}

func (w *testWatcher) Next() ([]*KeyValue, error) {
	select {
	case <-w.sig:
		return nil, nil
	case <-w.err:
		return nil, errors.New("error")
	case <-w.exit:
		return nil, nil
	}
}

func (w *testWatcher) Stop() error {
	close(w.exit)
	return nil
}

func TestConfig(t *testing.T) {
	var (
		err            error
		httpAddr       = "0.0.0.0"
		httpTimeout    = 0.5
		grpcPort       = 10080
		endpoint1      = "www.aaa.com"
		databaseDriver = "mysql"
	)

	c := New(
		WithSource(newTestJSONSource(_testJSON)),
		WithDecoder(defaultDecoder),
		WithResolver(defaultResolver),
	)
	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}

	jSource := newTestJSONSource(_testJSON)
	opts := options{
		sources:  []Source{jSource},
		decoder:  defaultDecoder,
		resolver: defaultResolver,
		merge:    defaultMerge,
	}
	cf := &config{}
	cf.opts = opts
	cf.reader = newReader(opts)

	err = cf.Load()
	if err != nil {
		t.Fatal(err)
	}

	driver, err := cf.Value("data.database.driver").String()
	if err != nil {
		t.Fatal(err)
	}
	if databaseDriver != driver {
		t.Fatal("databaseDriver is not equal to val")
	}

	driverGet, err := Get[string](cf, "data.database.driver")
	if err != nil {
		t.Fatal(err)
	}
	if databaseDriver != driverGet {
		t.Errorf("Get[string] want: %s, got: %s", databaseDriver, driverGet)
	}

	type HTTPConfig struct {
		Addr string `json:"addr"`
		Port int    `json:"port"`
	}
	v, err := Get[HTTPConfig](cf, "server.http")
	if err != nil {
		t.Fatal(err)
	} else if v.Addr != httpAddr {
		t.Errorf("Get[HttpConfig] Addr want: %s, got: %s", httpAddr, v.Addr)
	} else if v.Port != 80 {
		t.Errorf("Get[HttpConfig] Port want: 80, got: %d", v.Port)
	}

	err = cf.Watch("endpoints", func(string, Value) {})
	if err != nil {
		t.Fatal(err)
	}

	jSource.sig <- struct{}{}
	jSource.err <- struct{}{}

	var testConf testConfigStruct
	err = cf.Scan(&testConf)
	if err != nil {
		t.Fatal(err)
	}
	if httpAddr != testConf.Server.HTTP.Addr {
		t.Errorf("testConf.Server.HTTP.Addr want: %s, got: %s", httpAddr, testConf.Server.HTTP.Addr)
	}
	if httpTimeout != testConf.Server.HTTP.Timeout {
		t.Errorf("testConf.Server.HTTP.Timeout want: %.1f, got: %.1f", httpTimeout, testConf.Server.HTTP.Timeout)
	}
	if !testConf.Server.HTTP.EnableSSL {
		t.Error("testConf.Server.HTTP.EnableSSL is not equal to true")
	}
	if grpcPort != testConf.Server.GRPC.Port {
		t.Errorf("testConf.Server.GRPC.Port want: %d, got: %d", grpcPort, testConf.Server.GRPC.Port)
	}
	if endpoint1 != testConf.Endpoints[0] {
		t.Errorf("testConf.Endpoints[0] want: %s, got: %s", endpoint1, testConf.Endpoints[0])
	}
	if len(testConf.Endpoints) != 2 {
		t.Error("len(testConf.Endpoints) is not equal to 2")
	}
}

// stringKVSource is a config source that returns key-value pairs with
// string values (like env source), for testing type coercion during Scan.
type stringKVSource struct {
	kvs map[string]string
}

func newStringKVSource(kvs map[string]string) *stringKVSource {
	return &stringKVSource{kvs: kvs}
}

func (s *stringKVSource) Load() ([]*KeyValue, error) {
	var out []*KeyValue
	for k, v := range s.kvs {
		out = append(out, &KeyValue{
			Key:   k,
			Value: []byte(v),
		})
	}
	return out, nil
}

func (s *stringKVSource) Watch() (Watcher, error) {
	return newTestWatcher(make(chan struct{}), make(chan struct{})), nil
}

func TestScanStringToTyped(t *testing.T) {
	type nested struct {
		Debug bool    `json:"debug"`
		Port  int     `json:"port"`
		Rate  float64 `json:"rate"`
		Level string  `json:"level"`
	}
	type config struct {
		Debug   bool    `json:"debug"`
		Port    int     `json:"port"`
		Rate    float64 `json:"rate"`
		Level   string  `json:"level"`
		Version string  `json:"version"`
		Nested  nested  `json:"nested"`
	}

	src := newStringKVSource(map[string]string{
		"debug":        "true",
		"port":         "8080",
		"rate":         "0.9",
		"level":        "debug",
		"version":      "1234",
		"nested.debug": "false",
		"nested.port":  "9090",
		"nested.rate":  "1.5",
		"nested.level": "info",
	})

	c := New(WithSource(src))
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	var cfg config
	if err := c.Scan(&cfg); err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if cfg.Debug != true {
		t.Errorf("Debug want true, got %v", cfg.Debug)
	}
	if cfg.Port != 8080 {
		t.Errorf("Port want 8080, got %d", cfg.Port)
	}
	if cfg.Rate != 0.9 {
		t.Errorf("Rate want 0.9, got %f", cfg.Rate)
	}
	if cfg.Level != "debug" {
		t.Errorf("Level want 'debug', got %q", cfg.Level)
	}
	if cfg.Version != "1234" {
		t.Errorf("Version want '1234' (string), got %q", cfg.Version)
	}
	if cfg.Nested.Debug != false {
		t.Errorf("Nested.Debug want false, got %v", cfg.Nested.Debug)
	}
	if cfg.Nested.Port != 9090 {
		t.Errorf("Nested.Port want 9090, got %d", cfg.Nested.Port)
	}
	if cfg.Nested.Rate != 1.5 {
		t.Errorf("Nested.Rate want 1.5, got %f", cfg.Nested.Rate)
	}
	if cfg.Nested.Level != "info" {
		t.Errorf("Nested.Level want 'info', got %q", cfg.Nested.Level)
	}
}

func TestScanStringToBoolEdgeCases(t *testing.T) {
	type cfg struct {
		Debug bool `json:"debug"`
	}

	tests := []struct {
		name    string
		value   string
		want    bool
		wantErr bool
	}{
		{"true", "true", true, false},
		{"false", "false", false, false},
		{"1", "1", true, false},
		{"0", "0", false, false},
		{"TRUE", "TRUE", true, false},
		{"FALSE", "FALSE", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := newStringKVSource(map[string]string{"debug": tt.value})
			c := New(WithSource(src))
			if err := c.Load(); err != nil {
				t.Fatal(err)
			}
			var c cfg
			err := c.Scan(&c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && c.Debug != tt.want {
				t.Errorf("Debug want %v, got %v", tt.want, c.Debug)
			}
		})
	}
}
