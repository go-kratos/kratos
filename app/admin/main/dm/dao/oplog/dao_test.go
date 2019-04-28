package oplog

import (
	"testing"

	"go-common/app/admin/main/dm/conf"
	"go-common/library/log"
)

var (
	dao *Dao
)

func TestMain(m *testing.M) {
	var err error
	conf.ConfPath = "../../cmd/dm-admin-test.toml"
	if err = conf.Init(); err != nil {
		log.Error("conf.Init(%v)", err)
		return
	}
	dao = New(conf.Conf)
	m.Run()
	//log.Close()
}
