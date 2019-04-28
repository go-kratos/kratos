package dao

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/interface/bbq/app-bbq/conf"
	"testing"
)

var (
	d *Dao
)

func init() {
	flag.Set("conf", "../cmd/")
	flag.Set("conf_name", "test.toml")
	flag.Set("app_setting_name", "app_setting.toml")
	flag.Set("deploy.env", "uat")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func TestUserLike(t *testing.T) {
	res, err := d.CheckUserLike(context.Background(), 88895104, []int64{74850, 2222})
	fmt.Printf("res: %v", res)
	if err != nil && len(res) == 0 {
		t.Errorf("user like fail: err(%v)", err)
	}
}
