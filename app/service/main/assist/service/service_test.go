package service

import (
	"context"
	"flag"
	"go-common/app/service/main/assist/conf"
	"go-common/app/service/main/assist/model/assist"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/assist-service.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Check(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		MID       = int64(2089809)
		assistMid = int64(2089810)
	)
	Convey("checkBanned", t, WithService(func(s *Service) {
		err = s.checkBanned(c, MID)
		So(err, ShouldBeNil)
	}))
	Convey("checkFollow", t, WithService(func(s *Service) {
		err = s.checkFollow(c, MID, assistMid)
		So(err, ShouldBeNil)
	}))
	Convey("checkIdentify", t, WithService(func(s *Service) {
		err = s.checkIdentify(c, assistMid)
		So(err, ShouldBeNil)
	}))
	Convey("checkIsAssist", t, WithService(func(s *Service) {
		err = s.checkIsAssist(c, MID, assistMid)
		So(err, ShouldBeNil)
	}))
	Convey("checkIsNotAssist", t, WithService(func(s *Service) {
		err = s.checkIsNotAssist(c, MID, assistMid)
		So(err, ShouldBeNil)
	}))
}

func Test_Limit(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		MID       = int64(2089809)
		assistMid = int64(2089810)
	)
	Convey("checkMaxAssistCnt", t, WithService(func(s *Service) {
		err = s.checkMaxAssistCnt(c, MID)
		So(err, ShouldBeNil)
	}))
	Convey("checkTotalLimit", t, WithService(func(s *Service) {
		err = s.checkTotalLimit(c, MID)
		So(err, ShouldBeNil)
	}))
	Convey("checkSameLimit", t, WithService(func(s *Service) {
		err = s.checkSameLimit(c, MID, assistMid)
		So(err, ShouldBeNil)
	}))
}

func Test_Log(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		MID       = int64(2089809)
		assistMid = int64(2089810)
		logID     = int64(123)
		objID     = int64(1)
		tp        = int64(2)
		act       = int64(3)
		subID     = int64(4)
		logInfo   = &assist.Log{}
		logs      = []*assist.Log{}
		cnt       int64
		stime     = time.Now().Add(-time.Hour * 72)
		etime     = time.Now()
		objIDStr  = "hash_object_id"
		detail    = "detail string info"
	)
	Convey("CancelLog", t, WithService(func(s *Service) {
		err = s.CancelLog(c, MID, assistMid, logID)
		So(err, ShouldBeNil)
	}))
	Convey("LogInfo", t, WithService(func(s *Service) {
		logInfo, err = s.LogInfo(c, logID, MID, assistMid)
		So(err, ShouldBeNil)
		So(logInfo, ShouldNotBeNil)
	}))
	Convey("LogObj", t, WithService(func(s *Service) {
		logInfo, err = s.LogObj(c, MID, objID, tp, act)
		So(err, ShouldBeNil)
		So(logInfo, ShouldNotBeNil)
	}))
	Convey("LogCnt", t, WithService(func(s *Service) {
		cnt, err = s.LogCnt(c, MID, assistMid, stime, etime)
		So(err, ShouldBeNil)
		So(cnt, ShouldNotBeNil)
		So(cnt, ShouldBeGreaterThanOrEqualTo, 0)
	}))
	Convey("AddLog", t, WithService(func(s *Service) {
		err = s.AddLog(c, MID, assistMid, tp, act, subID, objIDStr, detail)
		So(err, ShouldBeNil)
	}))
	Convey("Logs", t, WithService(func(s *Service) {
		logs, err = s.Logs(c, MID, assistMid, stime, etime, 1, 20)
		So(err, ShouldBeNil)
		So(len(logs), ShouldBeGreaterThanOrEqualTo, 0)
	}))
}
