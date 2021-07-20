package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
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
		Http struct {
			Addr      string  `json:"addr"`
			Port      int     `json:"port"`
			Timeout   float64 `json:"timeout"`
			EnableSSL bool    `json:"enable_ssl"`
		} `json:"http"`
		GRpc struct {
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

type testJsonSource struct {
	data string
	sig  chan struct{}
	err  chan struct{}
}

func newTestJsonSource(data string) *testJsonSource {
	return &testJsonSource{data: data, sig: make(chan struct{}), err: make(chan struct{})}
}

func (p *testJsonSource) Load() ([]*KeyValue, error) {
	kv := &KeyValue{
		Key:    "json",
		Value:  []byte(p.data),
		Format: "json",
	}
	return []*KeyValue{kv}, nil
}

func (p *testJsonSource) Watch() (Watcher, error) {
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
		enableSSL      = true
		endpoint1      = "www.aaa.com"
		databaseDriver = "mysql"
	)

	c := New(
		WithSource(newTestJsonSource(_testJSON)),
		WithDecoder(defaultDecoder),
		WithResolver(defaultResolver),
		WithLogger(log.DefaultLogger),
	)
	err = c.Close()
	assert.Nil(t, err)

	jSource := newTestJsonSource(_testJSON)
	opts := options{
		sources:  []Source{jSource},
		decoder:  defaultDecoder,
		resolver: defaultResolver,
		logger:   log.DefaultLogger,
	}
	cf := &config{}
	cf.opts = opts
	cf.reader = newReader(opts)

	err = cf.Load()
	assert.Nil(t, err)

	val, err := cf.Value("data.database.driver").String()
	assert.Nil(t, err)
	assert.Equal(t, databaseDriver, val)

	err = cf.Watch("endpoints", func(key string, value Value) {
	})
	assert.Nil(t, err)

	jSource.sig <- struct{}{}
	jSource.err <- struct{}{}

	var testConf testConfigStruct
	err = cf.Scan(&testConf)
	assert.Nil(t, err)
	assert.Equal(t, httpAddr, testConf.Server.Http.Addr)
	assert.Equal(t, httpTimeout, testConf.Server.Http.Timeout)
	assert.Equal(t, enableSSL, testConf.Server.Http.EnableSSL)
	assert.Equal(t, grpcPort, testConf.Server.GRpc.Port)
	assert.Equal(t, endpoint1, testConf.Endpoints[0])
	assert.Equal(t, 2, len(testConf.Endpoints))
}

func TestDefaultResolver(t *testing.T) {
	var (
		portString = "8080"
		countInt   = 10
		enableBool = true
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
				"array": []interface{}{
					"${PORT}",
					map[string]interface{}{"foobar": "${NOTEXIST:8081}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "$PORT:default",
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
			expect: enableBool,
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
			name:   "test $value",
			path:   "foo.bar.value2",
			expect: portString,
		},
		{
			name:   "test $value:default",
			path:   "foo.bar.value3",
			expect: portString + ":default",
		},
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
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("expect: %#v, actural: %#v", test.expect, actual)
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
