package block

import (
	"flag"
	"os"
	"testing"

	"go-common/app/admin/main/member/conf"
	"go-common/library/cache/memcache"
	xsql "go-common/library/database/sql"
	bm "go-common/library/net/http/blademaster"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.member-admin")
		flag.Set("conf_token", "c18eac2285e4e4a75a8672139c30d464")
		flag.Set("tree_id", "2135")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/member-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	hc := bm.NewClient(conf.Conf.HTTPClient.Read)
	d = New(conf.Conf, hc, memcache.NewPool(conf.Conf.BlockMemcache), xsql.NewMySQL(conf.Conf.BlockMySQL))
	os.Exit(m.Run())
}
