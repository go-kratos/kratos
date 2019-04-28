package reply

import (
	"flag"
	"go-common/app/interface/main/reply/conf"
	"os"
	"testing"
)

var (
	d *Dao
	D *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.community.reply")
		flag.Set("conf_token", "54e85e3ab609f79ae908b9ea3e3f0775")
		flag.Set("tree_id", "2125")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/reply-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	D = d
	os.Exit(m.Run())
}
