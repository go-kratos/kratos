package memcache

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"go-common/library/container/pool"
	xtime "go-common/library/time"
)

var testMemcacheAddr = "127.0.0.1:11211"

var testConfig = &Config{
	Config: &pool.Config{
		Active:      10,
		Idle:        10,
		IdleTimeout: xtime.Duration(time.Second),
		WaitTimeout: xtime.Duration(time.Second),
		Wait:        false,
	},
	Proto:        "tcp",
	DialTimeout:  xtime.Duration(time.Second),
	ReadTimeout:  xtime.Duration(time.Second),
	WriteTimeout: xtime.Duration(time.Second),
}

func init() {
	if addr := os.Getenv("TEST_MEMCACHE_ADDR"); addr != "" {
		testMemcacheAddr = addr
	}
	testConfig.Addr = testMemcacheAddr
}

func TestMain(m *testing.M) {
	testClient = New(testConfig)
	m.Run()
	testClient.Close()
	os.Exit(0)
}

func ExampleConn_set() {
	var (
		err    error
		value  []byte
		conn   Conn
		expire int32 = 100
		p            = struct {
			Name string
			Age  int64
		}{"golang", 10}
	)
	cnop := DialConnectTimeout(time.Duration(time.Second))
	rdop := DialReadTimeout(time.Duration(time.Second))
	wrop := DialWriteTimeout(time.Duration(time.Second))
	if value, err = json.Marshal(p); err != nil {
		fmt.Println(err)
		return
	}
	if conn, err = Dial("tcp", testMemcacheAddr, cnop, rdop, wrop); err != nil {
		fmt.Println(err)
		return
	}
	// FlagRAW test
	itemRaw := &Item{
		Key:        "test_raw",
		Value:      value,
		Expiration: expire,
	}
	if err = conn.Set(itemRaw); err != nil {
		fmt.Println(err)
		return
	}
	// FlagGzip
	itemGZip := &Item{
		Key:        "test_gzip",
		Value:      value,
		Flags:      FlagGzip,
		Expiration: expire,
	}
	if err = conn.Set(itemGZip); err != nil {
		fmt.Println(err)
		return
	}
	// FlagGOB
	itemGOB := &Item{
		Key:        "test_gob",
		Object:     p,
		Flags:      FlagGOB,
		Expiration: expire,
	}
	if err = conn.Set(itemGOB); err != nil {
		fmt.Println(err)
		return
	}
	// FlagJSON
	itemJSON := &Item{
		Key:        "test_json",
		Object:     p,
		Flags:      FlagJSON,
		Expiration: expire,
	}
	if err = conn.Set(itemJSON); err != nil {
		fmt.Println(err)
		return
	}
	// FlagJSON | FlagGzip
	itemJSONGzip := &Item{
		Key:        "test_jsonGzip",
		Object:     p,
		Flags:      FlagJSON | FlagGzip,
		Expiration: expire,
	}
	if err = conn.Set(itemJSONGzip); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
}

func ExampleConn_get() {
	var (
		err   error
		item2 *Item
		conn  Conn
		p     struct {
			Name string
			Age  int64
		}
	)
	cnop := DialConnectTimeout(time.Duration(time.Second))
	rdop := DialReadTimeout(time.Duration(time.Second))
	wrop := DialWriteTimeout(time.Duration(time.Second))
	if conn, err = Dial("tcp", testMemcacheAddr, cnop, rdop, wrop); err != nil {
		fmt.Println(err)
		return
	}
	if item2, err = conn.Get("test_raw"); err != nil {
		fmt.Println(err)
	} else {
		if err = conn.Scan(item2, &p); err != nil {
			fmt.Printf("FlagRAW conn.Scan error(%v)\n", err)
			return
		}
	}
	// FlagGZip
	if item2, err = conn.Get("test_gzip"); err != nil {
		fmt.Println(err)
	} else {
		if err = conn.Scan(item2, &p); err != nil {
			fmt.Printf("FlagGZip conn.Scan error(%v)\n", err)
			return
		}
	}
	// FlagGOB
	if item2, err = conn.Get("test_gob"); err != nil {
		fmt.Println(err)
	} else {
		if err = conn.Scan(item2, &p); err != nil {
			fmt.Printf("FlagGOB conn.Scan error(%v)\n", err)
			return
		}
	}
	// FlagJSON
	if item2, err = conn.Get("test_json"); err != nil {
		fmt.Println(err)
	} else {
		if err = conn.Scan(item2, &p); err != nil {
			fmt.Printf("FlagJSON conn.Scan error(%v)\n", err)
			return
		}
	}
	// Output:
}

func ExampleConn_getMulti() {
	var (
		err  error
		conn Conn
		res  map[string]*Item
		keys = []string{"test_raw", "test_gzip"}
		p    struct {
			Name string
			Age  int64
		}
	)
	cnop := DialConnectTimeout(time.Duration(time.Second))
	rdop := DialReadTimeout(time.Duration(time.Second))
	wrop := DialWriteTimeout(time.Duration(time.Second))
	if conn, err = Dial("tcp", testMemcacheAddr, cnop, rdop, wrop); err != nil {
		fmt.Println(err)
		return
	}
	if res, err = conn.GetMulti(keys); err != nil {
		fmt.Printf("conn.GetMulti(%v) error(%v)", keys, err)
		return
	}
	for _, v := range res {
		if err = conn.Scan(v, &p); err != nil {
			fmt.Printf("conn.Scan error(%v)\n", err)
			return
		}
		fmt.Println(p)
	}
	// Output:
	//{golang 10}
	//{golang 10}
}
