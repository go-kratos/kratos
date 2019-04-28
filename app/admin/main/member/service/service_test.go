package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/admin/main/member/conf"

	"github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	flag.Parse()
	flag.Set("conf", "../cmd/member-admin-test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}

	s = New(conf.Conf)
}

func TestPing(t *testing.T) {
	convey.Convey("Ping", t, func() {
		err := s.Ping(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestFaceCheck(t *testing.T) {
	convey.Convey("faceCheck", t, func() {
		etime := time.Now().Unix()
		stime := time.Now().AddDate(0, 0, -2).Unix()
		err := s.faceAuditAI(context.Background(), stime, etime)
		convey.So(err, convey.ShouldBeNil)
	})
}
