package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-kratos/kratos/v2/config"
)

func TestConfig_Load(t *testing.T) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	source := NewConfigSource(client, WithGroup("test"), WithDataID("test.yaml"))

	type fields struct {
		source config.Source
	}
	tests := []struct {
		name      string
		fields    fields
		want      []*config.KeyValue
		wantErr   bool
		preFunc   func(t *testing.T)
		deferFunc func(t *testing.T)
	}{
		{
			name: "normal",
			fields: fields{
				source: source,
			},
			wantErr: false,
			preFunc: func(t *testing.T) {
				_, err = client.PublishConfig(vo.ConfigParam{DataId: "test.yaml", Group: "test", Content: "test: test"})
				if err != nil {
					t.Error(err)
				}
				time.Sleep(time.Second * 1)
			},
			deferFunc: func(t *testing.T) {
				_, dErr := client.DeleteConfig(vo.ConfigParam{DataId: "test.yaml", Group: "test"})
				if dErr != nil {
					t.Error(dErr)
				}
			},
			want: []*config.KeyValue{{
				Key:    "test.yaml",
				Value:  []byte("test: test"),
				Format: "yaml",
			}},
		},
		{
			name: "error",
			fields: fields{
				source: source,
			},
			wantErr: false,
			preFunc: func(t *testing.T) {
				_, err = client.PublishConfig(vo.ConfigParam{DataId: "111.yaml", Group: "notExist", Content: "test: test"})
				if err != nil {
					t.Error(err)
				}
				time.Sleep(time.Second * 1)
			},
			deferFunc: func(t *testing.T) {
				_, dErr := client.DeleteConfig(vo.ConfigParam{DataId: "111.yaml", Group: "notExist"})
				if dErr != nil {
					t.Error(dErr)
				}
			},
			want: []*config.KeyValue{{
				Key:    "test.yaml",
				Value:  []byte{},
				Format: "yaml",
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.preFunc != nil {
				test.preFunc(t)
			}
			if test.deferFunc != nil {
				defer test.deferFunc(t)
			}
			s := test.fields.source
			configs, lErr := s.Load()
			if (lErr != nil) != test.wantErr {
				t.Errorf("Load error = %v, wantErr %v", lErr, test.wantErr)
				t.Errorf("Load configs = %v", configs)
				return
			}
			if !reflect.DeepEqual(configs, test.want) {
				t.Errorf("Load configs = %v, want %v", configs, test.want)
			}
		})
	}
}

func TestConfig_Watch(t *testing.T) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	source := NewConfigSource(client, WithGroup("test"), WithDataID("test.yaml"))

	type fields struct {
		source config.Source
	}
	tests := []struct {
		name        string
		fields      fields
		want        []*config.KeyValue
		wantErr     bool
		processFunc func(t *testing.T, w config.Watcher)
		deferFunc   func(t *testing.T, w config.Watcher)
	}{
		{
			name: "normal",
			fields: fields{
				source: source,
			},
			wantErr: false,
			processFunc: func(t *testing.T, w config.Watcher) {
				_, pErr := client.PublishConfig(vo.ConfigParam{DataId: "test.yaml", Group: "test", Content: "test: test"})
				if pErr != nil {
					t.Error(pErr)
				}
			},
			deferFunc: func(t *testing.T, w config.Watcher) {
				_, dErr := client.DeleteConfig(vo.ConfigParam{DataId: "test.yaml", Group: "test"})
				if dErr != nil {
					t.Error(dErr)
				}
			},
			want: []*config.KeyValue{{
				Key:    "test.yaml",
				Value:  []byte("test: test"),
				Format: "yaml",
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.fields.source
			watch, wErr := s.Watch()
			if wErr != nil {
				t.Error(wErr)
				return
			}
			if test.processFunc != nil {
				test.processFunc(t, watch)
			}
			if test.deferFunc != nil {
				defer test.deferFunc(t, watch)
			}
			want, nErr := watch.Next()
			if (nErr != nil) != test.wantErr {
				t.Errorf("Watch error = %v, wantErr %v", nErr, test.wantErr)
				return
			}
			if !reflect.DeepEqual(want, test.want) {
				t.Errorf("Watch watcher = %v, want %v", watch, test.want)
			}
		})
	}
}
