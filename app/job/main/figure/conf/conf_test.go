package conf

import (
	"flag"
	"testing"
)

func TestInit(t *testing.T) {
	flag.Parse()
	var err error
	if err = Init(); err != nil {
		t.Fatal(err)
	}
	if Conf == nil {
		t.Fatal()
	}
}
