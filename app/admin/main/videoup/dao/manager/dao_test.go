package manager

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func TestPing(t *testing.T) {
	Convey("Ping", t, WithDao(func(d *Dao) {
		err := d.Ping(context.TODO())
		So(err, ShouldBeNil)
	}))
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-admin")
		flag.Set("conf_token", "gRSfeavV7kJdY9875Gf29pbd2wrdKZ1a")
		flag.Set("tree_id", "2307")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/videoup-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
