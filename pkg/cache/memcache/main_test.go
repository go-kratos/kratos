package memcache

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"go-common/library/container/pool"
	"go-common/library/testing/lich"
	xtime "go-common/library/time"
)

var testConnASCII Conn
var testMemcache *Memcache
var testPool *Pool
var testMemcacheAddr string

func setupTestConnASCII(addr string) {
	var err error
	cnop := DialConnectTimeout(time.Duration(2 * time.Second))
	rdop := DialReadTimeout(time.Duration(2 * time.Second))
	wrop := DialWriteTimeout(time.Duration(2 * time.Second))
	testConnASCII, err = Dial("tcp", addr, cnop, rdop, wrop)
	if err != nil {
		log.Fatal(err)
	}
	testConnASCII.Delete("test")
	testConnASCII.Delete("test1")
	testConnASCII.Delete("test2")
	if err != nil {
		log.Fatal(err)
	}
}

func setupTestMemcache(addr string) {
	testConfig := &Config{
		Config: &pool.Config{
			Active:      10,
			Idle:        10,
			IdleTimeout: xtime.Duration(time.Second),
			WaitTimeout: xtime.Duration(time.Second),
			Wait:        false,
		},
		Addr:         addr,
		Proto:        "tcp",
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	testMemcache = New(testConfig)
}

func setupTestPool(addr string) {
	config := &Config{
		Name:         "test",
		Proto:        "tcp",
		Addr:         addr,
		DialTimeout:  xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
	}
	config.Config = &pool.Config{
		Active:      10,
		Idle:        5,
		IdleTimeout: xtime.Duration(90 * time.Second),
	}
	testPool = NewPool(config)
}

func TestMain(m *testing.M) {
	flag.Set("f", "./test/docker-compose.yaml")
	flag.Parse()
	if err := lich.Setup(); err != nil {
		panic(err)
	}
	testMemcacheAddr = "127.0.0.1:11211"
	setupTestConnASCII(testMemcacheAddr)
	setupTestMemcache(testMemcacheAddr)
	setupTestPool(testMemcacheAddr)
	// TODO: add setupexample?
	testExampleAddr = testMemcacheAddr

	ret := m.Run()
	lich.Teardown()
	os.Exit(ret)
}
