package dao

import (
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/job/main/aegis/conf"
	"go-common/app/job/main/aegis/model"
)

var (
	d     *Dao
	task1 = &model.Task{ID: 1, BusinessID: 1, FlowID: 1, UID: 1, Weight: 3, Ctime: model.IntTime(time.Now().Unix())}
	task2 = &model.Task{ID: 2, BusinessID: 1, FlowID: 1, UID: 1, Weight: 2, Ctime: model.IntTime(time.Now().Unix())}
	task3 = &model.Task{ID: 3, BusinessID: 1, FlowID: 1, UID: 1, Weight: 1, Ctime: model.IntTime(time.Now().Unix())}
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-job")
		flag.Set("conf_token", "aed3cc21ca345ffc284c6036da32352b")
		flag.Set("tree_id", "61819")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/aegis-job.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
