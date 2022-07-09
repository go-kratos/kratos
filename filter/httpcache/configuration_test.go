package httpcache

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v2"
)

var dummyConfig = []byte(`
httpcache:
  api:
    basepath: /souin-api
    prometheus:
      basepath: /anything-for-prometheus-metrics
    souin:
      basepath: /anything-for-souin
  cache_keys:
    .+:
      disable_body: true
      disable_host: true
      disable_method: true
  default_cache:
    allowed_http_verbs:
      - GET
      - POST
      - HEAD
    badger:
      url: /badger/url
      path: /badger/path
      configuration:
        SyncEnable: false
    cdn:
      api_key: XXXX
      dynamic: true
      hostname: XXXX
      network: XXXX
      provider: fastly
      strategy: soft
    etcd:
      url: /etcd/url
      path: /etcd/path
      configuration:
        SyncEnable: false
    headers:
      - Authorization
    nuts:
      url: /etcd/url
      path: /etcd/path
      configuration:
        SyncEnable: false
        ValueFloat: 1.123
        ValueDuration: 1s
        ValueInt: 2
        ValueObject:
          ValueFloat: 1.123
    olric:
      url: 'olric:3320'
      path: /olric/path
      configuration:
        SyncEnable: false
    regex:
      exclude: 'ARegexHere'
    ttl: 10s
    stale: 10s
    default_cache_control: no-store
  log_level: debug
  urls:
    'https:\/\/domain.com\/first-.+':
      ttl: 1000s
    'https:\/\/domain.com\/second-route':
      ttl: 10s
      headers:
      - Authorization
    'https?:\/\/mysubdomain\.domain\.com':
      ttl: 50s
      default_cache_control: no-cache
      headers:
      - Authorization
      - 'Content-Type'
  ykeys:
    The_First_Test:
      headers:
        Content-Type: '.+'
    The_Second_Test:
      url: 'the/second/.+'
  surrogate_keys:
    The_First_Test:
      headers:
        Content-Type: '.+'
    The_Second_Test:
      url: 'the/second/.+'
`)

func Test_ParseConfiguration(t *testing.T) {
	filename := "/tmp/httpcache-" + time.Now().String() + ".yml"
	ioutil.WriteFile(filename, dummyConfig, 0777)
	c := config.New(
		config.WithSource(file.NewSource(filename)),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	_ = ParseConfiguration(c)
}
