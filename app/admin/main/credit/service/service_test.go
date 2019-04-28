package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/credit/conf"
	"go-common/app/admin/main/credit/model/blocked"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

// func CleanCache() {
// 	//c := context.TODO()
// 	//pool := redis.NewPool(conf.Conf.Redis.Config)
// 	//pool.Get(c).Do("FLUSHDB")
// }

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

// func WithService(f func(s *Service)) func() {
// 	return func() {
// 		Reset(func() { CleanCache() })
// 		f(s)
// 	}
// }

func TestService_loadConfig(t *testing.T) {
	s = New(conf.Conf)
	s.loadConfig()
	s.loadManager()
	//t.Logf("config (%+v)", s.c.Judge)
}

func Test_Infos(t *testing.T) {
	arg := &blocked.ArgBlockedSearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.Infos(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_InfosEx(t *testing.T) {
	arg := &blocked.ArgBlockedSearch{}
	Convey("return someting", t, func() {
		list, err := s.InfosEx(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
	})
}

func Test_Publishs(t *testing.T) {
	arg := &blocked.ArgPublishSearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.Publishs(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_FormatCSV(t *testing.T) {
	var records [][]string
	record := []string{"sss", "ssss"}
	records = append(records, record)
	Convey("return someting", t, func() {
		buf := s.FormatCSV(records)
		So(buf, ShouldNotBeNil)
	})
}

func Test_Upload(t *testing.T) {
	Convey("return someting", t, func() {
		_, err := s.Upload(context.TODO(), "blocked_info", "", 12313, nil)
		So(err, ShouldBeNil)
	})
}

func Test_AddJury(t *testing.T) {
	arg := &blocked.ArgAddJurys{MIDs: []int64{111, 22}, OPID: 111, Day: 111, Send: 1}
	Convey("return someting", t, func() {
		err := s.AddJury(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func Test_Cases(t *testing.T) {
	arg := &blocked.ArgCaseSearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.Cases(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_Opinions(t *testing.T) {
	arg := &blocked.ArgOpinionSearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.Opinions(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_Jurys(t *testing.T) {
	arg := &blocked.ArgJurySearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.Jurys(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}

func Test_JurysEx(t *testing.T) {
	arg := &blocked.ArgJurySearch{}
	Convey("return someting", t, func() {
		res, err := s.JurysEx(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_KPIPoint(t *testing.T) {
	arg := &blocked.ArgKpiPointSearch{}
	Convey("return someting", t, func() {
		list, pager, err := s.KPIPoint(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeNil)
		So(pager, ShouldNotBeNil)
	})
}
