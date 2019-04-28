package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/ep/saga/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/saga-admin-test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New()
	defer d.Close()
	os.Exit(m.Run())
}
