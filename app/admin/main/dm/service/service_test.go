package service

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/dm/conf"
	"go-common/app/admin/main/dm/model/oplog"
	"go-common/library/log"
	manager "go-common/library/queue/databus/report"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	svr *Service
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	dir, _ := filepath.Abs("../cmd/dm-admin-test.toml")
	if err = flag.Set("conf", dir); err != nil {
		panic(err)
	}
	if err = conf.Init(); err != nil {
		panic(err)
	}
	// log.Init(nil)
	svr = New(conf.Conf)
	// manager log init
	manager.InitManager(conf.Conf.ManagerLog)
	os.Exit(m.Run())
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(svr)
	}
}

func TestOpLog(t *testing.T) {
	var (
		c                  = context.TODO()
		cid          int64 = 12
		typ          int32 = 1
		operator     int64 = 219
		dmids              = []int64{719150137}
		subject            = "status"
		originVal          = "1"
		currentVal         = "2"
		remark             = "备注"
		source             = oplog.SourceManager
		operatorType       = oplog.OperatorAdmin
	)
	Convey("OpLog", t, WithService(func(s *Service) {
		err := svr.OpLog(c, cid, operator, typ, dmids, subject, originVal, currentVal, remark, source, operatorType)
		So(err, ShouldBeNil)
	}))
}

func TestQueryOpLogs(t *testing.T) {
	var (
		c          = context.TODO()
		dmid int64 = 719150137
	)
	Convey("QueryOpLogs", t, WithService(func(s *Service) {
		objs, err := svr.QueryOpLogs(c, dmid)
		log.Info("%v", objs)
		So(err, ShouldBeNil)
	}))
}
