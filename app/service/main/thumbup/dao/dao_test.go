package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/thumbup/conf"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.thumbup-service")
		flag.Set("conf_token", "VhnSEtd0oymsNQaDUYuEknoWu2mVOOVK")
		flag.Set("tree_id", "7720")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/thumbup-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

// INSERT INTO `bilibili_likes`.`counts_01`(`id`, `mtime`, `ctime`, `business_id`, `origin_id`, `message_id`, `likes_count`, `dislikes_count`, `likes_change`, `dislikes_change`) VALUES (98, '2018-06-13 14:51:11', '2018-06-13 14:51:11', 1, 99901, 8888, 1, 2, 3, 4);
// INSERT INTO `bilibili_likes`.`counts_01`(`id`, `mtime`, `ctime`, `business_id`, `origin_id`, `message_id`, `likes_count`, `dislikes_count`, `likes_change`, `dislikes_change`) VALUES (99, '2018-06-13 14:56:56', '2018-06-13 14:56:56', 1, 101, 8888, 0, 0, 0, 0);
// INSERT INTO `bilibili_likes`.`likes`(`id`, `mtime`, `ctime`, `business_id`, `origin_id`, `message_id`, `mid`, `type`) VALUES (0, '2018-11-01 18:03:28', '2018-11-01 18:03:28', 1, 1, 1, 1, 1);
// INSERT INTO `bilibili_likes`.`counts`(`id`, `mtime`, `ctime`, `business_id`, `origin_id`, `message_id`, `likes_count`, `dislikes_count`, `likes_change`, `dislikes_change`, `up_mid`) VALUES (0, '2018-11-03 12:16:25', '2018-11-03 12:16:25', 1, 1, 1, 1, 1, 0, 0, 0);
