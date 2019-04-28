package dao

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"go-common/app/service/main/history/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.history-service")
		flag.Set("conf_token", "10f1bb6e589c42e7e1ee2560aff96b81")
		flag.Set("tree_id", "56699")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		fmt.Println("") // 存在才能pass
		flag.Set("conf", "../cmd/history-service-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

// INSERT INTO `bilibili_history`.`histories`(`mtime`, `ctime`, `mid`, `business_id`, `kid`, `aid`, `sid`, `epid`, `sub_type`, `cid`, `device`, `progress`, `view_at`) VALUES ('2018-08-27 03:03:50', '2018-08-27 03:01:29', 1, 4, 2, 3, 4, 5, 7, 6, 8, 9, 10);
