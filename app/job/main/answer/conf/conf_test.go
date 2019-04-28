package conf

import (
	"flag"
	"testing"
)

func TestInit(t *testing.T) {
	flag.Set("conf", "../cmd/answer-job-test.toml")
	flag.Parse()
	var err error
	if err = Init(); err != nil {
		t.Fatal(err)
	}
	if Conf == nil {
		t.Fatal()
	}
}
