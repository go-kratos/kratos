package etcd

import (
	"context"
	"reflect"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

const testKey = "/kratos/test/config"

func TestConfig(t *testing.T) {
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
	if _, err = client.Put(context.Background(), testKey, "test config"); err != nil {
		t.Fatal(err)
	}

	source, err := New(client, WithPath(testKey))
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := source.Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != testKey || string(kvs[0].Value) != "test config" {
		t.Fatal("config error")
	}

	w, err := source.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()

	if _, err = client.Put(context.Background(), testKey, "new config"); err != nil {
		t.Error(err)
	}

	if kvs, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != testKey || string(kvs[0].Value) != "new config" {
		t.Fatal("config error")
	}

	if _, err := client.Delete(context.Background(), testKey); err != nil {
		t.Error(err)
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

func TestEtcdWithPath(t *testing.T) {
	tests := []struct {
		name   string
		fields string
		want   string
	}{
		{
			name:   "default",
			fields: testKey,
			want:   testKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &options{
				ctx: context.Background(),
			}

			got := WithPath(tt.fields)
			got(options)

			if options.path != tt.want {
				t.Errorf("WithPath(tt.fields) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdWithPrefix(t *testing.T) {
	tests := []struct {
		name   string
		fields bool
		want   bool
	}{
		{
			name:   "default",
			fields: false,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &options{
				ctx: context.Background(),
			}

			got := WithPrefix(tt.fields)
			got(options)

			if options.prefix != tt.want {
				t.Errorf("WithPrefix(tt.fields) = %v, want %v", got, tt.want)
			}
		})
	}
}
