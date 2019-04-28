package upcrm

import (
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/job/main/up/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.mcn-admin")
		flag.Set("conf_token", "220af473858ad67f75586b66bece0e6b")
		flag.Set("tree_id", "58930")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}

	flag.Set("conf", "../../cmd/up-job.toml")

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestGetDueTask(t *testing.T) {
	var date, _ = time.Parse("2006-01-02", "2018-11-02")
	var res, err = d.GetDueTask(date)
	if err != nil {
		t.Errorf("err get tasks, err=%v", err)
		t.FailNow()
	}
	for i, v := range res {
		t.Logf("[%d]res=%+v", i, v)
	}
}
