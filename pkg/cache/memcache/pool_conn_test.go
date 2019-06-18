package memcache

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/container/pool"
	xtime "github.com/bilibili/kratos/pkg/time"
)

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

func TestPoolSet(t *testing.T) {
	conn := testPool.Get(context.Background())
	defer conn.Close()
	// set
	if err := conn.Set(itempool); err != nil {
		t.Errorf("memcache: set error(%v)", err)
	} else {
		t.Logf("memcache: set value: %s", itempool.Value)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("memcache: close error(%v)", err)
	}
}

func TestPoolGet(t *testing.T) {
	key := "testpool"
	conn := testPool.Get(context.Background())
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
	conn := testPool.Get(context.Background())
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
	key := "testpool"
	conn := testPool.Get(context.Background())
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
	key := "test_count"
	conn := testPool.Get(context.Background())
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
	conn := testPool.Get(context.Background())
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
	conn := testPool.Get(context.Background())
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
	key := "testpool"
	conn := testPool.Get(context.Background())
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
		Addr:         testMemcacheAddr,
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	c.Config = &pool.Config{
		Active:      10,
		Idle:        5,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	testPool = NewPool(c)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn := testPool.Get(context.Background())
			if err := conn.Close(); err != nil {
				b.Errorf("memcache: close error(%v)", err)
			}
		}
	})
	if err := testPool.Close(); err != nil {
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
	conn := testPool.Get(context.Background())
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
	key := largeValue.Key
	conn := testPool.Get(context.Background())
	defer conn.Close()
	// get
	var err error
	if _, err = conn.Get(key); err != nil {
		t.Errorf("memcache: large get error(%+v)", err)
	}
}

func TestPoolGetMultiLargeValue(t *testing.T) {
	conn := testPool.Get(context.Background())
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
	conn := testPool.Get(context.Background())
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
	key := largeValueBoundary.Key
	conn := testPool.Get(context.Background())
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
		conn = testPool.Get(context.Background())
	)
	defer conn.Close()
	conn.Delete(key)
	if err := conn.Add(item); err != nil {
		t.Errorf("memcache: add error(%v)", err)
	}
	if err := conn.Add(item); err != ErrNotStored {
		t.Errorf("memcache: add error(%v)", err)
	}
}

func TestNewPool(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		wantPanic bool
	}{
		{
			"NewPoolIllegalDialTimeout",
			args{
				&Config{
					Name:         "test_illegal_dial_timeout",
					Proto:        "tcp",
					Addr:         testMemcacheAddr,
					DialTimeout:  xtime.Duration(-time.Second),
					ReadTimeout:  xtime.Duration(time.Second),
					WriteTimeout: xtime.Duration(time.Second),
				},
			},
			nil,
			true,
		},
		{
			"NewPoolIllegalReadTimeout",
			args{
				&Config{
					Name:         "test_illegal_read_timeout",
					Proto:        "tcp",
					Addr:         testMemcacheAddr,
					DialTimeout:  xtime.Duration(time.Second),
					ReadTimeout:  xtime.Duration(-time.Second),
					WriteTimeout: xtime.Duration(time.Second),
				},
			},
			nil,
			true,
		},
		{
			"NewPoolIllegalWriteTimeout",
			args{
				&Config{
					Name:         "test_illegal_write_timeout",
					Proto:        "tcp",
					Addr:         testMemcacheAddr,
					DialTimeout:  xtime.Duration(time.Second),
					ReadTimeout:  xtime.Duration(time.Second),
					WriteTimeout: xtime.Duration(-time.Second),
				},
			},
			nil,
			true,
		},
		{
			"NewPool",
			args{
				&Config{
					Name:         "test_new",
					Proto:        "tcp",
					Addr:         testMemcacheAddr,
					DialTimeout:  xtime.Duration(time.Second),
					ReadTimeout:  xtime.Duration(time.Second),
					WriteTimeout: xtime.Duration(time.Second),
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("wantPanic recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			if gotP := NewPool(tt.args.cfg); gotP == nil {
				t.Error("NewPool() failed, got nil")
			}
		})
	}
}

func TestPool_Get(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		p       *Pool
		args    args
		wantErr bool
		n       int
	}{
		{
			"Get",
			NewPool(&Config{
				Config: &pool.Config{
					Active: 3,
					Idle:   2,
				},
				Name:         "test_get",
				Proto:        "tcp",
				Addr:         testMemcacheAddr,
				DialTimeout:  xtime.Duration(time.Second),
				ReadTimeout:  xtime.Duration(time.Second),
				WriteTimeout: xtime.Duration(time.Second),
			}),
			args{context.TODO()},
			false,
			3,
		},
		{
			"GetExceededPoolSize",
			NewPool(&Config{
				Config: &pool.Config{
					Active: 3,
					Idle:   2,
				},
				Name:         "test_get_out",
				Proto:        "tcp",
				Addr:         testMemcacheAddr,
				DialTimeout:  xtime.Duration(time.Second),
				ReadTimeout:  xtime.Duration(time.Second),
				WriteTimeout: xtime.Duration(time.Second),
			}),
			args{context.TODO()},
			true,
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 1; i <= tt.n; i++ {
				got := tt.p.Get(tt.args.ctx)
				if reflect.TypeOf(got) == reflect.TypeOf(errConn{}) {
					if !tt.wantErr {
						t.Errorf("got errConn, export Conn")
					}
					return
				} else {
					if tt.wantErr {
						if i > tt.p.c.Active {
							t.Errorf("got Conn, export errConn")
						}
					}
				}
			}
		})
	}
}

func TestPool_Close(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		p       *Pool
		args    args
		wantErr bool
		g       int
		c       int
	}{
		{
			"Close",
			NewPool(&Config{
				Config: &pool.Config{
					Active: 1,
					Idle:   1,
				},
				Name:         "test_get",
				Proto:        "tcp",
				Addr:         testMemcacheAddr,
				DialTimeout:  xtime.Duration(time.Second),
				ReadTimeout:  xtime.Duration(time.Second),
				WriteTimeout: xtime.Duration(time.Second),
			}),
			args{context.TODO()},
			false,
			3,
			3,
		},
		{
			"CloseExceededPoolSize",
			NewPool(&Config{
				Config: &pool.Config{
					Active: 1,
					Idle:   1,
				},
				Name:         "test_get_out",
				Proto:        "tcp",
				Addr:         testMemcacheAddr,
				DialTimeout:  xtime.Duration(time.Second),
				ReadTimeout:  xtime.Duration(time.Second),
				WriteTimeout: xtime.Duration(time.Second),
			}),
			args{context.TODO()},
			true,
			5,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 1; i <= tt.g; i++ {
				got := tt.p.Get(tt.args.ctx)
				if err := got.Close(); err != nil {
					if !tt.wantErr {
						t.Error(err)
					}
				}
				if i <= tt.c {
					if err := got.Close(); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}
