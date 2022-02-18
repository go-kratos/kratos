package consul

import (
	"reflect"
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
