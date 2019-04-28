package conf

import (
	"flag"
	"testing"
)

func TestInit(t *testing.T) {
	flag.Set("conf", "../cmd/figure-service-test.toml")
	var err error
	if err = Init(); err != nil {
		t.Fatal(err)
	}
	if Conf == nil {
		t.Fatal()
	}
	if Conf.Verify == nil {
		t.Fatal()
	}
	if Conf.Verify.HTTPClient == nil {
		t.Fatal()
	}
}
