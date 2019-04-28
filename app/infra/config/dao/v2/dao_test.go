package v2

import (
	"flag"
	"go-common/app/infra/config/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	// if os.Getenv("DEPLOY_ENV") != "" {
	// 	// flag.Set("app_id", "")
	// 	// flag.Set("conf_token", "")
	// 	// flag.Set("tree_id", "")
	// 	// flag.Set("conf_version", "docker-1")
	// 	// flag.Set("deploy_env", "fat1")
	// 	// flag.Set("conf_host", "config.bilibili.co")
	// 	// flag.Set("conf_path", "/tmp")
	// 	// flag.Set("region", "sh")
	// 	// flag.Set("zone", "sh001")
	// } else {
	// 	flag.Set("conf", "../../cmd/config-service-example.toml")
	// }
	flag.Set("conf", "../../cmd/config-service-example.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}
