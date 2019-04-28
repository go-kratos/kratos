package dao

import (
	"context"
	"flag"
	"go-common/app/service/live/live-dm/conf"
	"path/filepath"
	"testing"
)

func init() {
	dir, _ := filepath.Abs("../cmd/test.toml")
	flag.Set("conf", dir)
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
}

// DEPLOY_ENV=uat go test -run TestIsWhietListUID
func TestIsWhietListUID(t *testing.T) {
	dao := New(conf.Conf)
	if isWhite := dao.IsWhietListUID(context.TODO(), "1111"); isWhite {
		t.Error("白名单判断失败")
	}
	if isWhite := dao.IsWhietListUID(context.TODO(), "13269933"); !isWhite {
		t.Error("白名单判断失败")
	}
}
