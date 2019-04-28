package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/growup/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "mobile.studio.growup-interface")
		flag.Set("conf_token", "c68ad4f01bc8c39a3fa6242623e79ffb")
		flag.Set("tree_id", "13584")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/growup-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func Exec(c context.Context, sql string) (rows int64, err error) {
	res, err := d.db.Exec(c, sql)
	if err != nil {
		return
	}
	return res.RowsAffected()
}
