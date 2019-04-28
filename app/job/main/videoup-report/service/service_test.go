package service

import (
	"context"
	"flag"
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/videoup-report/conf"
	"go-common/app/job/main/videoup-report/model/archive"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/videoup-report-job.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_loadType(t *testing.T) {
	Convey("loadType", t, func() {
		s.loadType()
	})
}

func Test_VideoAudit(t *testing.T) {
	Convey("VideoAudit", t, func() {
		_, err := s.VideoAudit(context.Background(), time.Now(), time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_Ping(t *testing.T) {
	Convey("Ping", t, func() {
		err := s.Ping(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_TaskTooksByHalfHour(t *testing.T) {
	Convey("TaskTooksByHalfHour", t, func() {
		_, err := s.TaskTooksByHalfHour(context.Background(), time.Now(), time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_AddArchiveHotRecheck(t *testing.T) {
	Convey("AddArchiveHotRecheck", t, func() {
		time.Sleep(time.Second)
		err := s.addHotRecheck()
		So(err, ShouldBeNil)
	})
}

func Test_SecondRound(t *testing.T) {
	Convey("SecondRound", t, func() {
		m := &archive.VideoupMsg{
			Route:    "second_round",
			Aid:      24320325,
			FromList: "hot_review",
		}
		err := s.secondRound(context.Background(), m)
		So(err, ShouldBeNil)
	})
}

func Test_SecondRoundCancelMission(t *testing.T) {
	Convey("SecondRound", t, func() {
		m := &archive.VideoupMsg{
			Route:     "second_round",
			Aid:       17191032,
			MissionID: 1,
		}
		err := s.secondRound(context.Background(), m)
		So(err, ShouldBeNil)
	})
}

func Test_SecondFormat(t *testing.T) {
	Convey("SecondFormat", t, func() {
		m := 5556
		format := secondsFormat(m)
		So(format, ShouldNotBeNil)
	})
}
