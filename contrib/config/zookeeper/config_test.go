package zookeeper

import (
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
)

const testNamespace = "/kratos/test/config"

const testKey = "key"

func TestConfig(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*5)
	if err != nil {
		t.Fatal(err)
		return
	}

	fullPath := path.Join(testNamespace, testKey)

	if ok, _, _ := conn.Exists(fullPath); ok {
		if err = conn.Delete(fullPath, 1); err != nil {
			t.Fatal(err)
		}
	}

	if _, err = conn.Create(fullPath, []byte("test config"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatal(err)
	}

	s, err := New(
		conn,
		WithNamespace(testNamespace),
		WithKey(testKey),
	)
	if err != nil {
		t.Fatal(err)
	}

	kvs, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != "key" || string(kvs[0].Value) != "test config" {
		t.Fatal("config error")
	}

	w, err := s.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()

	if _, err = conn.Set(fullPath, []byte("new config"), 0); err != nil {
		t.Fatal(err)
	}

	if kvs, err = w.Next(); err != nil {
		t.Fatal(err)
	}

	if len(kvs) != 1 || kvs[0].Key != "key" || string(kvs[0].Value) != "new config" {
		t.Fatal("config error")
	}

	if err = conn.Delete(fullPath, 1); err != nil {
		t.Fatal(err)
	}
}

func TestExtToFormat(t *testing.T) {
	conn, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second*5)
	if err != nil {
		t.Fatal(err)
		return
	}
	namespace := "/kratos/test/ext"
	key := "a.json"
	value := `{"a":1}`

	fullPath := path.Join(namespace, key)

	if ok, _, _ := conn.Exists(fullPath); ok {
		if err = conn.Delete(fullPath, 0); err != nil {
			t.Fatal(err)
		}
	}

	if _, err = conn.Create(fullPath, []byte(value), 0, zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatal(err)
	}

	source, err := New(conn, WithNamespace(namespace), WithKey(key))
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
	if !reflect.DeepEqual(key, kvs[0].Key) {
		t.Errorf("kvs[0].Key is %s", kvs[0].Key)
	}
	if !reflect.DeepEqual(value, string(kvs[0].Value)) {
		t.Errorf("kvs[0].Value is %s", kvs[0].Value)
	}
	if !reflect.DeepEqual("json", kvs[0].Format) {
		t.Errorf("kvs[0].Format is %s", kvs[0].Format)
	}
}
