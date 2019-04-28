package dao

import (
	"context"
	"flag"
	"go-common/app/job/main/growup/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "mobile.studio.growup-job")
		flag.Set("conf_token", "8781e02680f40996bc01eb1248ac2ac9")
		flag.Set("tree_id", "14716")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/growup-job.toml")
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
