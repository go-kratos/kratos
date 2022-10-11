package consul

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/go-kratos/kratos/v2/config"
)

const testPath = "kratos/test/config"

const testKey = "kratos/test/config/key"

func TestConfig(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err = client.KV().Put(&api.KVPair{Key: testKey, Value: []byte("test config")}, nil); err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(testPath))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != "key" || string(kvs[0].Value) != "test config" {
		t.Fatal("config error")
	}

	w, err := source.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()

	if _, err = client.KV().Put(&api.KVPair{Key: testKey, Value: []byte("new config")}, nil); err != nil {
		t.Error(err)
	}

	if kvs, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != "key" || string(kvs[0].Value) != "new config" {
		t.Fatal("config error")
	}

	if _, err := client.KV().Delete(testKey, nil); err != nil {
		t.Error(err)
	}
}

func TestExtToFormat(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		t.Fatal(err)
	}
	tp := "kratos/test/ext"
	tn := "a.bird.json"
	tk := tp + "/" + tn
	tc := `{"a":1}`
	if _, err = client.KV().Put(&api.KVPair{Key: tk, Value: []byte(tc)}, nil); err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(tp))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(kvs), 1) {
		t.Errorf("len(kvs) is %d", len(kvs))
	}
	if !reflect.DeepEqual(tn, kvs[0].Key) {
		t.Errorf("kvs[0].Key is %s", kvs[0].Key)
	}
	if !reflect.DeepEqual(tc, string(kvs[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kvs[0].Value)
	}
	if !reflect.DeepEqual("json", kvs[0].Format) {
		t.Errorf("kvs[0].Format is %s", kvs[0].Format)
	}
}

func Test_source_Watch(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(testPath))
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		source config.Source
	}

	type args struct {
		key   string
		value string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		wantErr   bool
		deferFunc func(t *testing.T)
	}{
		{
			name:   "normal",
			fields: fields{source: source},
			args: args{
				key:   testKey,
				value: "test value",
			},
			want:    "test value",
			wantErr: false,
			deferFunc: func(t *testing.T) {
				_, err := client.KV().Delete(testKey, nil)
				if err != nil {
					t.Error(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}

			got, err := tt.fields.source.Watch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			time.Sleep(100 * time.Millisecond)
			_, err = client.KV().Put(&api.KVPair{Key: tt.args.key, Value: []byte(tt.args.value)}, nil)
			if err != nil {
				t.Error(err)
			}

			next, err := got.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(next) != 1 {
				t.Error("watch is error")
			}

			if !reflect.DeepEqual(string(next[0].Value), tt.want) {
				t.Errorf("Watch got = %v, want %v", string(next[0].Value), tt.want)
			}
		})
	}
}

func Test_source_Load(t *testing.T) {
	client, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(testPath))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		key   string
		value string
	}
	type fields struct {
		source config.Source
	}
	tests := []struct {
		name      string
		args      args
		fields    fields
		want      []*config.KeyValue
		wantErr   bool
		deferFunc func(t *testing.T)
	}{
		{
			name: "normal",
			args: args{
				key:   testKey,
				value: "test value",
			},
			fields: fields{
				source: source,
			},
			want: []*config.KeyValue{
				{
					Key:   "key",
					Value: []byte("test value"),
				},
			},
			deferFunc: func(t *testing.T) {
				_, err1 := client.KV().Delete(testKey, nil)
				if err1 != nil {
					t.Error(err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.deferFunc != nil {
				defer tt.deferFunc(t)
			}
			_, err = client.KV().Put(&api.KVPair{Key: tt.args.key, Value: []byte(tt.args.value)}, nil)
			if err != nil {
				t.Error(err)
			}
			got, err := tt.fields.source.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got[0], tt.want[0]) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
