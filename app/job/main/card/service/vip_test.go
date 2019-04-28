package service

import (
	"context"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/card/conf"
	"go-common/app/job/main/card/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	c = context.TODO()
	s *Service
)

func init() {
	var (
		err error
	)
	flag.Set("conf", "../cmd/test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	c = context.Background()
	if s == nil {
		s = New(conf.Conf)
	}
	time.Sleep(time.Second)
}

// go test  -test.v -test.run TestChangeEquipTime
func TestChangeEquipTime(t *testing.T) {
	Convey("TestChangeEquipTime ", t, func() {
		So(s.ChangeEquipTime(c, &model.VipReq{
			Mid:            1,
			VipType:        2,
			VipStatus:      1,
			VipOverdueTime: time.Now().Unix(),
		}), ShouldBeNil)
	})
}
