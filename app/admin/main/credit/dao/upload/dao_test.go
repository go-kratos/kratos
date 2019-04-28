package upload

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"go-common/app/admin/main/credit/conf"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var d *Dao

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	d = New(conf.Conf)
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		f(d)
	}
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

func ctx() context.Context {
	return context.Background()
}

func Test_Upload(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := d.Upload(context.TODO(), "blocked_info", "", 12313, nil)
		So(err, ShouldBeNil)
	})
}
