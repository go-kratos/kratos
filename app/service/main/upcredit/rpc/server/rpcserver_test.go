package server

import (
	"flag"
	"path/filepath"
	"testing"
	"time"
	_ "time"

	"go-common/app/service/main/upcredit/conf"
	"go-common/app/service/main/upcredit/service"
	"go-common/library/net/rpc"
	xtime "go-common/library/time"
)

func init() {
	dir, _ := filepath.Abs("../../cmd/upcredit-service.toml")
	flag.Set("conf", dir)
}

func initSvrAndClient(t *testing.T) (client *rpc.Client, err error) {
	if err = conf.Init(); err != nil {
		t.Errorf("conf.Init() error(%v)", err)
		t.FailNow()
	}
	svr := service.New(conf.Conf)
	New(conf.Conf, svr)

	client = rpc.Dial("127.0.0.1:6079", xtime.Duration(time.Second), nil)
	return
}

func TestInfo(t *testing.T) {
	client, err := initSvrAndClient(t)
	defer client.Close()
	if err != nil {
		t.Errorf("rpc.Dial error(%v)", err)
		t.FailNow()
	}
	//time.Sleep(1 * time.Second)
}
