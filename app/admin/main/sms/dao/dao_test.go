package dao

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/sms/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/sms-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		d.Ping(context.TODO())
	})
}

func Test_Close(t *testing.T) {
	Convey("close", t, func() {
		d.Close()
	})
}
