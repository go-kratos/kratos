package dao

import (
	"os"
	"testing"

	"go-common/app/job/main/dm/conf"
	"go-common/library/log"
)

var (
	testDao *Dao
)

func TestMain(m *testing.M) {
	conf.ConfPath = "../cmd/dm-job-test.toml"
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	testDao = New(conf.Conf)
	os.Exit(m.Run())
}
