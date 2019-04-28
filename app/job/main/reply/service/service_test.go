package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/job/main/reply/conf"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ser *Service
)

func CleanCache() {

}

func TestMain(m *testing.M) {
	dir, _ := filepath.Abs("../cmd/reply-job-test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		panic("conf init err:" + err.Error())
	}
	ser = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func WithService(f func(ser *Service)) func() {
	return func() {
		Reset(func() { CleanCache() })
		f(ser)
	}
}
func TestClose(t *testing.T) {
	err := ser.Close()
	So(err, ShouldBeNil)
}

func TestPing(t *testing.T) {
	err := ser.Ping(context.Background())
	So(err, ShouldBeNil)
}
