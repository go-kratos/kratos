package dsn

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	xtime "go-common/library/time"
)

type config struct {
	Network   string         `dsn:"network"`
	Addresses []string       `dsn:"address"`
	Username  string         `dsn:"username"`
	Password  string         `dsn:"password"`
	Timeout   xtime.Duration `dsn:"query.timeout"`
	Sub       Sub            `dsn:"query.sub"`
	Def       string         `dsn:"query.def,hello"`
}

type Sub struct {
	Foo int `dsn:"query.foo"`
}

func TestBind(t *testing.T) {
	var cfg config
	rawdsn := "tcp://root:toor@172.12.23.34,178.23.34.45?timeout=1s&sub.foo=1&hello=world"
	dsn, err := Parse(rawdsn)
	if err != nil {
		t.Fatal(err)
	}
	values, err := dsn.Bind(&cfg)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(values, url.Values{"hello": {"world"}}) {
		t.Errorf("unexpect values get %v", values)
	}
	cfg2 := config{
		Network:   "tcp",
		Addresses: []string{"172.12.23.34", "178.23.34.45"},
		Password:  "toor",
		Username:  "root",
		Sub:       Sub{Foo: 1},
		Timeout:   xtime.Duration(time.Second),
		Def:       "hello",
	}
	if !reflect.DeepEqual(cfg, cfg2) {
		t.Errorf("unexpect config get %v, expect %v", cfg, cfg2)
	}
}

type config2 struct {
	Network string         `dsn:"network"`
	Address string         `dsn:"address"`
	Timeout xtime.Duration `dsn:"query.timeout"`
}

func TestUnix(t *testing.T) {
	var cfg config2
	rawdsn := "unix:///run/xxx.sock?timeout=1s&sub.foo=1&hello=world"
	dsn, err := Parse(rawdsn)
	if err != nil {
		t.Fatal(err)
	}
	_, err = dsn.Bind(&cfg)
	if err != nil {
		t.Error(err)
	}
	cfg2 := config2{
		Network: "unix",
		Address: "/run/xxx.sock",
		Timeout: xtime.Duration(time.Second),
	}
	if !reflect.DeepEqual(cfg, cfg2) {
		t.Errorf("unexpect config2 get %v, expect %v", cfg, cfg2)
	}
}
