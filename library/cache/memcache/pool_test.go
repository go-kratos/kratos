package memcache

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"go-common/library/container/pool"
	xtime "go-common/library/time"
)

var p *Pool
var config *Config

func init() {
	testMemcacheAddr := "127.0.0.1:11211"
	if addr := os.Getenv("TEST_MEMCACHE_ADDR"); addr != "" {
		testMemcacheAddr = addr
	}
	config = &Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         testMemcacheAddr,
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	config.Config = &pool.Config{
		Active:      10,
		Idle:        5,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
}

var itempool = &Item{
	Key:        "testpool",
	Value:      []byte("testpool"),
	Flags:      0,
	Expiration: 60,
	cas:        0,
}
var itempool2 = &Item{
	Key:        "test_count",
	Value:      []byte("0"),
	Flags:      0,
	Expiration: 1000,
	cas:        0,
}

type testObject struct {
	Mid   int64
	Value []byte
}

var largeValue = &Item{
	Key:        "large_value",
	Flags:      FlagGOB | FlagGzip,
	Expiration: 1000,
	cas:        0,
}

var largeValueBoundary = &Item{
	Key:        "large_value",
	Flags:      FlagGOB | FlagGzip,
	Expiration: 1000,
	cas:        0,
}

func prepareEnv2() {
	if p != nil {
		return
	}
	p = NewPool(config)
}

func TestPoolSet(t *testing.T) {
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	// set
	if err := conn.Set(itempool); err != nil {
		t.Errorf("memcache: set error(%v)", err)
	} else {
		t.Logf("memcache: set value: %s", item.Value)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolGet(t *testing.T) {
	prepareEnv2()
	key := "testpool"
	conn := p.Get(context.Background())
	defer conn.Close()
	// get
	if res, err := conn.Get(key); err != nil {
		t.Errorf("memcache: get error(%v)", err)
	} else {
		t.Logf("memcache: get value: %s", res.Value)
	}
	if _, err := conn.Get("not_found"); err != ErrNotFound {
		t.Errorf("memcache: expceted err is not found but got: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolGetMulti(t *testing.T) {
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	s := []string{"testpool", "test1"}
	// get
	if res, err := conn.GetMulti(s); err != nil {
		t.Errorf("memcache: gets error(%v)", err)
	} else {
		t.Logf("memcache: gets value: %d", len(res))
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolTouch(t *testing.T) {
	prepareEnv2()
	key := "testpool"
	conn := p.Get(context.Background())
	defer conn.Close()
	// touch
	if err := conn.Touch(key, 10); err != nil {
		t.Errorf("memcache: touch error(%v)", err)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolIncrement(t *testing.T) {
	prepareEnv2()
	key := "test_count"
	conn := p.Get(context.Background())
	defer conn.Close()
	// set
	if err := conn.Set(itempool2); err != nil {
		t.Errorf("memcache: set error(%v)", err)
	} else {
		t.Logf("memcache: set value: 0")
	}
	// incr
	if res, err := conn.Increment(key, 1); err != nil {
		t.Errorf("memcache: incr error(%v)", err)
	} else {
		t.Logf("memcache: incr n: %d", res)
		if res != 1 {
			t.Errorf("memcache: expected res=1 but got %d", res)
		}
	}
	// decr
	if res, err := conn.Decrement(key, 1); err != nil {
		t.Errorf("memcache: decr error(%v)", err)
	} else {
		t.Logf("memcache: decr n: %d", res)
		if res != 0 {
			t.Errorf("memcache: expected res=0 but got %d", res)
		}
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolErr(t *testing.T) {
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
	if err := conn.Err(); err == nil {
		t.Errorf("memcache: err not nil")
	} else {
		t.Logf("memcache: err: %v", err)
	}
}

func TestPoolCompareAndSwap(t *testing.T) {
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	key := "testpool"
	//cas
	if r, err := conn.Get(key); err != nil {
		t.Errorf("conn.Get() error(%v)", err)
	} else {
		r.Value = []byte("shit")
		if err := conn.CompareAndSwap(r); err != nil {
			t.Errorf("conn.Get() error(%v)", err)
		}
		r, _ := conn.Get("testpool")
		if r.Key != "testpool" || !bytes.Equal(r.Value, []byte("shit")) || r.Flags != 0 {
			t.Error("conn.Get() error, value")
		}
		if err := conn.Close(); err != nil {
			t.Errorf("memcache: close error(%v)", err)
		}
	}
}

func TestPoolDel(t *testing.T) {
	prepareEnv2()
	key := "testpool"
	conn := p.Get(context.Background())
	defer conn.Close()
	// delete
	if err := conn.Delete(key); err != nil {
		t.Errorf("memcache: delete error(%v)", err)
	} else {
		t.Logf("memcache: delete key: %s", key)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func BenchmarkMemcache(b *testing.B) {
	c := &Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         "127.0.0.1:11211",
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	c.Config = &pool.Config{
		Active:      10,
		Idle:        5,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	p = NewPool(c)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn := p.Get(context.Background())
			if err := conn.Close(); err != nil {
				b.Errorf("memcache: close error(%v)", err)
			}
		}
	})
	if err := p.Close(); err != nil {
		b.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolSetLargeValue(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < 4000000; i++ {
		b.WriteByte(1)
	}
	obj := &testObject{}
	obj.Mid = 1000
	obj.Value = b.Bytes()
	largeValue.Object = obj
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	// set
	if err := conn.Set(largeValue); err != nil {
		t.Errorf("memcache: set error(%v)", err)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolGetLargeValue(t *testing.T) {
	prepareEnv2()
	key := largeValue.Key
	conn := p.Get(context.Background())
	defer conn.Close()
	// get
	var err error
	if _, err = conn.Get(key); err != nil {
		t.Errorf("memcache: large get error(%+v)", err)
	}
}

func TestPoolGetMultiLargeValue(t *testing.T) {
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	s := []string{largeValue.Key, largeValue.Key}
	// get
	if res, err := conn.GetMulti(s); err != nil {
		t.Errorf("memcache: gets error(%v)", err)
	} else {
		t.Logf("memcache: gets value: %d", len(res))
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolSetLargeValueBoundary(t *testing.T) {
	var b bytes.Buffer
	for i := 0; i < _largeValue; i++ {
		b.WriteByte(1)
	}
	obj := &testObject{}
	obj.Mid = 1000
	obj.Value = b.Bytes()
	largeValueBoundary.Object = obj
	prepareEnv2()
	conn := p.Get(context.Background())
	defer conn.Close()
	// set
	if err := conn.Set(largeValueBoundary); err != nil {
		t.Errorf("memcache: set error(%v)", err)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolGetLargeValueBoundary(t *testing.T) {
	prepareEnv2()
	key := largeValueBoundary.Key
	conn := p.Get(context.Background())
	defer conn.Close()
	// get
	var err error
	if _, err = conn.Get(key); err != nil {
		t.Errorf("memcache: large get error(%v)", err)
	}
}

func TestPoolAdd(t *testing.T) {
	var (
		key  = "test_add"
		item = &Item{
			Key:        key,
			Value:      []byte("0"),
			Flags:      0,
			Expiration: 60,
			cas:        0,
		}
		conn = p.Get(context.Background())
	)
	defer conn.Close()
	prepareEnv2()
	conn.Delete(key)
	if err := conn.Add(item); err != nil {
		t.Errorf("memcache: add error(%v)", err)
	}
	if err := conn.Add(item); err != ErrNotStored {
		t.Errorf("memcache: add error(%v)", err)
	}
}
