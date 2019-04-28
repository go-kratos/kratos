package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/growup/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func CleanMysql() {
}

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "mobile.studio.growup-admin")
		flag.Set("conf_token", "ac1fd397cbc33eb60541e8734844bdd5")
		flag.Set("tree_id", "13583")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/growup-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func WithMysql(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanMysql() })
		f(d)
	}
}

// Exec do exec
func Exec(c context.Context, sql string) (rows int64, err error) {
	res, err := d.rddb.Exec(c, sql)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
