package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/history/model"

	. "github.com/smartystreets/goconvey/convey"
)

// TestService_History
func TestService_History(t *testing.T) {
	var (
		c           = context.TODO()
		mid   int64 = 14771787
		aid   int64 = 5463823
		aids        = []int64{5463823}
		sid   int64 = 5730
		cid   int64 = 97791
		epid  int64 = 97922
		pro   int64 = 155
		rtime int64 = 1490958549
		tp    int8  = 1
		dt    int8  = 2
		pn          = 1
		ps          = 100
		now         = time.Now().Unix()
		h           = &model.History{Aid: aid, Unix: now, Sid: sid, Epid: epid, Cid: cid, Pro: pro, TP: int8(tp), DT: int8(dt)}
	)
	Convey("history ", t, WithService(func(s *Service) {
		Convey("history AddHistory ", func() {
			err := s.AddHistory(c, mid, rtime, h)
			So(err, ShouldBeNil)
		})
		Convey("history Progress", func() {
			_, err := s.Progress(c, mid, aids)
			So(err, ShouldBeNil)
		})
		Convey("history DelHistory", func() {
			err := s.DelHistory(c, mid, aids, 3)
			So(err, ShouldBeNil)
		})
		Convey("history ClearHistory", func() {
			err := s.ClearHistory(c, mid, []int8{3})
			So(err, ShouldBeNil)
		})
		Convey("history  Videos", func() {
			_, err := s.Videos(c, mid, pn, ps, 3)
			So(err, ShouldBeNil)
		})
		Convey("history AVHistories", func() {
			_, err := s.AVHistories(c, mid)
			So(err, ShouldBeNil)
		})
		Convey("history Histories", func() {
			_, err := s.Histories(c, mid, 1, 2, 3)
			So(err, ShouldBeNil)
		})
		Convey("history  SetShadow", func() {
			err := s.SetShadow(c, mid, 1)
			So(err, ShouldBeNil)
		})
		Convey("history Shadow", func() {
			_, err := s.Shadow(c, mid)
			So(err, ShouldBeNil)
		})
		Convey("history Manager", func() {
			_, err := s.ManagerHistory(c, false, mid)
			So(err, ShouldBeNil)
		})
	}))
}

func TestService_AddHistory(t *testing.T) {
	var (
		c           = context.TODO()
		mid   int64 = 14771787
		aid   int64 = 5463823
		sid   int64 = 5730
		cid   int64 = 97791
		epid  int64 = 97922
		pro   int64 = 155
		rtime int64 = 1490958549
		tp    int8  = 1
		dt    int8  = 2
		now         = time.Now().Unix()
		h           = &model.History{Aid: aid, Unix: now, Sid: sid, Epid: epid, Cid: cid, Pro: pro, TP: int8(tp), DT: int8(dt)}
	)
	Convey("history ", t, WithService(func(s *Service) {
		Convey("history AddHistory ", func() {
			err := s.AddHistory(c, mid, rtime, h)
			So(err, ShouldBeNil)
		})
	}))
}
