package config

import (
	"errors"
	"testing"

	"dario.cat/mergo"
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
		merge: func(dst, src interface{}) error {
			return mergo.Map(dst, src, mergo.WithOverride)
		},
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
