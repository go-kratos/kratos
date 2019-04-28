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
	InitAPI()
	InitGrpc(conf.Conf)
}

//group=qa01 DEPLOY_ENV=uat go test -run TestIncrDMNum
func TestIncrDMNum(t *testing.T) {
	IncrDMNum(context.TODO(), 5392, 2)
}
