package mcndao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/mcn/conf"

	"github.com/bluele/gcache"
	"go-common/app/interface/main/mcn/dao/global"
)

var (
	d          *Dao
	localcache = gcache.New(1024).Simple().Build()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.mcn-interface")
		flag.Set("conf_token", "49e4671bafbf93059aeb602685052ca0")
		flag.Set("tree_id", "58909")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../../cmd/mcn-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	global.Init(conf.Conf)
	d = New(conf.Conf, localcache)
	os.Exit(m.Run())
}
