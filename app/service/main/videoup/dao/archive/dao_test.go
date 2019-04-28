package archive

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/service/main/videoup/conf"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-service")
		flag.Set("conf_token", "4b62721602981eb3635dba3b0d866ac5")
		flag.Set("tree_id", "2308")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../../cmd/videoup-service.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	//conf.Init()
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
