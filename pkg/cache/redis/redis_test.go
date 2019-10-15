package redis

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/bilibili/kratos/pkg/container/pool"
	xtime "github.com/bilibili/kratos/pkg/time"
)

func TestRedis(t *testing.T) {
	testSet(t, testPool)
	testSend(t, testPool)
	testGet(t, testPool)
	testErr(t, testPool)
	if err := testPool.Close(); err != nil {
		t.Errorf("redis: close error(%v)", err)
	}
	conn, err := NewConn(testConfig)
	if err != nil {
		t.Errorf("redis: new conn error(%v)", err)
	}
	if err := conn.Close(); err != nil {
		t.Errorf("redis: close error(%v)", err)
	}
}

func testSet(t *testing.T, p *Pool) {
	var (
		key   = "test"
		value = "test"
		conn  = p.Get(context.TODO())
	)
	defer conn.Close()
	if reply, err := conn.Do("set", key, value); err != nil {
		t.Errorf("redis: conn.Do(SET, %s, %s) error(%v)", key, value, err)
	} else {
		t.Logf("redis: set status: %s", reply)
	}
}

func testSend(t *testing.T, p *Pool) {
	var (
		key    = "test"
		value  = "test"
		expire = 1000
		conn   = p.Get(context.TODO())
	)
	defer conn.Close()
	if err := conn.Send("SET", key, value); err != nil {
		t.Errorf("redis: conn.Send(SET, %s, %s) error(%v)", key, value, err)
	}
	if err := conn.Send("EXPIRE", key, expire); err != nil {
		t.Errorf("redis: conn.Send(EXPIRE key(%s) expire(%d)) error(%v)", key, expire, err)
	}
	if err := conn.Flush(); err != nil {
		t.Errorf("redis: conn.Flush error(%v)", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := conn.Receive(); err != nil {
			t.Errorf("redis: conn.Receive error(%v)", err)
			return
		}
	}
	t.Logf("redis: set value: %s", value)
}

func testGet(t *testing.T, p *Pool) {
	var (
		key  = "test"
		conn = p.Get(context.TODO())
	)
	defer conn.Close()
	if reply, err := conn.Do("GET", key); err != nil {
		t.Errorf("redis: conn.Do(GET, %s) error(%v)", key, err)
	} else {
		t.Logf("redis: get value: %s", reply)
	}
}

func testErr(t *testing.T, p *Pool) {
	conn := p.Get(context.TODO())
	if err := conn.Close(); err != nil {
		t.Errorf("redis: close error(%v)", err)
	}
	if err := conn.Err(); err == nil {
		t.Errorf("redis: err not nil")
	} else {
		t.Logf("redis: err: %v", err)
	}
}

func BenchmarkRedis(b *testing.B) {
	conf := &Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         testRedisAddr,
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	conf.Config = &pool.Config{
		Active:      10,
		Idle:        5,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	benchmarkPool := NewPool(conf)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn := benchmarkPool.Get(context.TODO())
			if err := conn.Close(); err != nil {
				b.Errorf("redis: close error(%v)", err)
			}
		}
	})
	if err := benchmarkPool.Close(); err != nil {
		b.Errorf("redis: close error(%v)", err)
	}
}

var testRedisCommands = []struct {
	args     []interface{}
	expected interface{}
}{
	{
		[]interface{}{"PING"},
		"PONG",
	},
	{
		[]interface{}{"SET", "foo", "bar"},
		"OK",
	},
	{
		[]interface{}{"GET", "foo"},
		[]byte("bar"),
	},
	{
		[]interface{}{"GET", "nokey"},
		nil,
	},
	{
		[]interface{}{"MGET", "nokey", "foo"},
		[]interface{}{nil, []byte("bar")},
	},
	{
		[]interface{}{"INCR", "mycounter"},
		int64(1),
	},
	{
		[]interface{}{"LPUSH", "mylist", "foo"},
		int64(1),
	},
	{
		[]interface{}{"LPUSH", "mylist", "bar"},
		int64(2),
	},
	{
		[]interface{}{"LRANGE", "mylist", 0, -1},
		[]interface{}{[]byte("bar"), []byte("foo")},
	},
}

func TestNewRedis(t *testing.T) {
	type args struct {
		c       *Config
		options []DialOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"new_redis",
			args{
				testConfig,
				make([]DialOption, 0),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRedis(tt.args.c, tt.args.options...)
			if r == nil {
				t.Errorf("NewRedis() error, got nil")
				return
			}
			err := r.Close()
			if err != nil {
				t.Errorf("Close() error %v", err)
			}
		})
	}
}

func TestRedis_Do(t *testing.T) {
	r := NewRedis(testConfig)
	r.Do(context.TODO(), "FLUSHDB")

	for _, cmd := range testRedisCommands {
		actual, err := r.Do(context.TODO(), cmd.args[0].(string), cmd.args[1:]...)
		if err != nil {
			t.Errorf("Do(%v) returned error %v", cmd.args, err)
			continue
		}
		if !reflect.DeepEqual(actual, cmd.expected) {
			t.Errorf("Do(%v) = %v, want %v", cmd.args, actual, cmd.expected)
		}
	}
	err := r.Close()
	if err != nil {
		t.Errorf("Close() error %v", err)
	}
}

func TestRedis_Conn(t *testing.T) {

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		p       *Redis
		args    args
		wantErr bool
		g       int
		c       int
	}{
		{
			"Close",
			NewRedis(&Config{
				Config: &pool.Config{
					Active: 1,
					Idle:   1,
				},
				Name:         "test_get",
				Proto:        "tcp",
				Addr:         testRedisAddr,
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
			NewRedis(&Config{
				Config: &pool.Config{
					Active: 1,
					Idle:   1,
				},
				Name:         "test_get_out",
				Proto:        "tcp",
				Addr:         testRedisAddr,
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
				got := tt.p.Conn(tt.args.ctx)
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

func BenchmarkRedisDoPing(b *testing.B) {
	r := NewRedis(testConfig)
	defer r.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := r.Do(context.Background(), "PING"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedisDoSET(b *testing.B) {
	r := NewRedis(testConfig)
	defer r.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := r.Do(context.Background(), "SET", "a", "b"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedisDoGET(b *testing.B) {
	r := NewRedis(testConfig)
	defer r.Close()
	r.Do(context.Background(), "SET", "a", "b")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := r.Do(context.Background(), "GET", "b"); err != nil {
			b.Fatal(err)
		}
	}
}
