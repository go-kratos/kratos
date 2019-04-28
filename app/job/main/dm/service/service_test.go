package service

import (
	"os"
	"testing"

	"go-common/app/job/main/dm/conf"
	"go-common/library/log"
)

var testSvc *Service

func TestMain(m *testing.M) {
	conf.ConfPath = "../cmd/dm-job-test.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	testSvc = New(conf.Conf)
	os.Exit(m.Run())
}
