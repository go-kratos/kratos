package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/reply/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

var nowTs = time.Now().Unix()

func Test_OnlineNoticeConflict(t *testing.T) {
	c := context.Background()
	nt := model.Notice{
		Plat:       model.PlatAndroid,
		Condition:  model.ConditionGT,
		Build:      1113,
		Title:      "测试",
		Status:     model.StatusOffline,
		Content:    "测试内容",
		Link:       "http://www.bilibili.com",
		StartTime:  xtime.Time(nowTs),
		EndTime:    xtime.Time(nowTs + 10*3600),
		ClientType: "",
	}
	nt2 := model.Notice{
		Plat:       model.PlatAndroid,
		Condition:  model.ConditionGT,
		Build:      1000,
		Title:      "测试2",
		Status:     model.StatusOffline,
		Content:    "测试内容2",
		Link:       "http://www.bilibili.com",
		StartTime:  xtime.Time(nowTs - 5*3600),
		EndTime:    xtime.Time(nowTs + 5*3600),
		ClientType: "android",
	}
	Convey("test notice data conflict ", t, WithService(func(s *Service) {
		id1, err := s.CreateNotice(c, &nt)
		So(err, ShouldBeNil)
		id2, err := s.CreateNotice(c, &nt2)
		So(err, ShouldBeNil)

		defer func() {
			s.DeleteNotice(c, uint32(id1))
			s.DeleteNotice(c, uint32(id2))
		}()

		err = s.UpdateNoticeStatus(c, model.StatusOnline, uint32(id1))
		So(err, ShouldBeNil)
		err = s.UpdateNoticeStatus(c, model.StatusOnline, uint32(id2))
		So(err, ShouldNotBeNil)
	}))
}
