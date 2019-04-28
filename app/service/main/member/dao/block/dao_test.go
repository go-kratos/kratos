package block

import (
	"flag"
	"os"
	"testing"

	"go-common/app/service/main/member/conf"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.member-service")
		flag.Set("conf_token", "ef70dbff7ee115ce242c67e633b21c29")
		flag.Set("tree_id", "2137")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	c := conf.Conf
	d = New(c,
		sql.NewMySQL(c.BlockMySQL),
		memcache.NewPool(c.BlockMemcache),
		bm.NewClient(c.HTTPClient),
		nil,
	)
	m.Run()
	os.Exit(0)
}
