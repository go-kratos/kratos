package consul

import (
	"testing"

	"github.com/hashicorp/consul/api"
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
