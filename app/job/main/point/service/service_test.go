package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/point/conf"
	"go-common/app/job/main/point/model"
	"go-common/library/log"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/point-job.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestAddPoint
func TestAddPoint(t *testing.T) {
	Convey("TestAddPoint", t, func() {
		p := &model.VipPoint{
			Mid:          111,
			PointBalance: 1,
			Ver:          1,
		}
		err := s.AddPoint(c, p)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdatePoint
func TestUpdatePoint(t *testing.T) {
	Convey("TestUpdatePoint", t, func() {
		p := &model.VipPoint{
			Mid:          27515415,
			PointBalance: 2,
			Ver:          260,
		}
		oldp := &model.VipPoint{
			Ver: 259,
		}
		err := s.UpdatePoint(c, p, oldp)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestChangeHistory
func TestChangeHistory(t *testing.T) {
	Convey("TestChangeHistory", t, func() {
		h := &model.VipPointChangeHistoryMsg{
			Mid:          111,
			Point:        2,
			OrderID:      "wqwqe22112",
			ChangeType:   10,
			ChangeTime:   "2018-03-14 18:30:57",
			RelationID:   "11",
			PointBalance: 111,
		}
		err := s.AddPointHistory(c, h)
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestUpdatePointByHistory
func TestUpdatePointByHistory(t *testing.T) {
	Convey("TestUpdatePointByHistory", t, func() {
		var (
			count   int64
			err     error
			hs      []*model.VipPointChangeHistoryMsg
			history *model.VipPointChangeHistoryMsg
		)
		history = &model.VipPointChangeHistoryMsg{
			Mid:          int64(26645632),
			Point:        int64(150),
			OrderID:      "201804121752494042821812",
			ChangeType:   3,
			ChangeTime:   "2018-04-12 17:52:49",
			RelationID:   "seasonId:6339",
			PointBalance: int64(150),
			Remark:       "欢迎来到实力至上主义的教室",
			Operator:     "",
		}
		hs = append(hs, history)
		history = &model.VipPointChangeHistoryMsg{
			Mid:          int64(26645632),
			Point:        int64(100),
			OrderID:      "201804121755088867934612",
			ChangeType:   3,
			ChangeTime:   "2018-04-12 17:55:09",
			RelationID:   "seasonId:1699",
			PointBalance: int64(250),
			Remark:       "四月是你的谎言",
			Operator:     "",
		}
		hs = append(hs, history)
		for _, history := range hs {
			if count, err = s.dao.HistoryCount(c, history.Mid, history.OrderID); err != nil {
				log.Error("update point(%v) history(%v)", err, history)
				continue
			}
			So(err, ShouldBeNil)
			So(count == 0, ShouldBeTrue)
			changeTime, err := time.ParseInLocation("2006-01-02 15:04:05", history.ChangeTime, time.Local)
			So(err, ShouldBeNil)
			ph := &model.PointHistory{
				Mid:          history.Mid,
				Point:        history.Point,
				OrderID:      history.OrderID,
				ChangeType:   int(history.ChangeType),
				ChangeTime:   xtime.Time(changeTime.Unix()),
				RelationID:   history.RelationID,
				PointBalance: history.PointBalance,
				Remark:       history.Remark,
				Operator:     history.Operator,
			}
			s.updatePointWithHistory(c, ph)
			So(err, ShouldBeNil)
		}
	})
}

// go test  -test.v -test.run TestFixData
func TestFixData(t *testing.T) {
	Convey("TestFixData", t, func() {
		err := s.fixdata("2018-04-12 18:20:00")
		So(err, ShouldBeNil)
	})
}

// go test  -test.v -test.run TestFixData
func TestNotify(t *testing.T) {
	Convey("TestNotify", t, func() {
		msg := &model.VipPointChangeHistoryMsg{
			Mid: 1,
		}
		err := s.Notify(context.TODO(), msg)
		So(err, ShouldBeNil)
	})
}
