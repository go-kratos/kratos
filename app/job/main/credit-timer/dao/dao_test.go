package dao

import (
	"context"
	"flag"
	"path/filepath"
	"strings"

	"go-common/app/job/main/credit-timer/conf"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/h2non/gock.v1"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func ctx() context.Context {
	return context.Background()
}
