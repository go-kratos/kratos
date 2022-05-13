package etcd

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func TestConfig(t *testing.T) {
	cases := map[string]string{
		"/kratos/test/db.yaml":    "db config",
		"/kratos/test/cache.yaml": "cache config",
		"/kratos/test/app.yaml":   "app config",
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	for key, val := range cases {
		if _, err = client.Put(context.Background(), key, val); err != nil {
			t.Fatal(err)
		}
	}

	var keys []string
	for key := range cases {
		keys = append(keys, key)
	}

	source, err := New(client, WithPath(keys...))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}

	for _, kv := range kvs {
		if val, exist := cases[kv.Key]; exist && string(kv.Value) != val {
			t.Fatalf("%q config error. expected: %q, actual: %q", kv.Key, kv.Value, val)
		}
	}

	w, err := source.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()

	for key, val := range cases {
		if _, err = client.Put(context.Background(), key, strings.ToUpper(val)); err != nil {
			t.Error(err)
		}
	}

	if kvs, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	for _, kv := range kvs {
		if val, exist := cases[kv.Key]; exist && string(kv.Value) != strings.ToUpper(val) {
			t.Fatalf("%q config error. expected: %q, actual: %q", kv.Key, kv.Value, val)
		}
	}

	for key := range cases {
		if _, err := client.Delete(context.Background(), key); err != nil {
			t.Error(err)
		}
	}
}

func TestExtToFormat(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	tp := "/kratos/test/ext"
	tn := "a.bird.json"
	tk := tp + "/" + tn
	tc := `{"a":1}`
	if _, err = client.Put(context.Background(), tk, tc); err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(tp), WithPrefix(true))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(len(kvs), 1) {
		t.Errorf("len(kvs) = %d", len(kvs))
	}
	if !reflect.DeepEqual(tk, kvs[0].Key) {
		t.Errorf("kvs[0].Key is %s", kvs[0].Key)
	}
	if !reflect.DeepEqual(tc, string(kvs[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kvs[0].Value)
	}
	if !reflect.DeepEqual("json", kvs[0].Format) {
		t.Errorf("kvs[0].Format is %s", kvs[0].Format)
	}
}
