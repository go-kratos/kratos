package conf

import (
	"flag"
	"testing"
)

func TestInit(t *testing.T) {
	flag.Set("conf", "../cmd/spy-service-test.toml")
	flag.Parse()
	var err error
	if err = Init(); err != nil {
		t.Fatal(err)
	}
	if Conf == nil {
		t.Fatal()
	}
	if Conf.Property.Punishment == nil {
		t.Fatal()
	}
	if Conf.Property.Score == nil {
		t.Fatal()
	}
	t.Log(Conf.Property.AutoBlockSwitch)

	t.Log(Conf.Property.White.Tels)
}
