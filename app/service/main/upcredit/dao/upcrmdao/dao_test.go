package upcrmdao

import (
	"flag"
	"github.com/go-sql-driver/mysql"
	"go-common/app/service/main/upcredit/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.upcredit-service")
		flag.Set("conf_token", "e85677843358d6c5dcd7246e6e3fc2de")
		flag.Set("tree_id", "33287")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/upcredit-service.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func IgnoreErr(err error) error {
	if err == nil {
		return err
	}
	var e, _ = err.(*mysql.MySQLError)
	if e != nil {
		switch e.Number {
		case 1062:
			return nil
		}
	}
	return err
}
