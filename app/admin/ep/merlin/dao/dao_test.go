package dao

import (
	"context"
	"flag"
	"path/filepath"
	"strings"

	"go-common/app/admin/ep/merlin/conf"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
	c context.Context
)

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
	d.httpClient.SetTransport(gock.DefaultTransport)
	c = ctx()
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func ctx() context.Context {
	return context.Background()
}

func WithPaasToken(f func()) func() {
	return func() {
		d.paasToken(c)
		f()
	}
}
