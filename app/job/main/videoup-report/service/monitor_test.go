package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/job/main/videoup-report/model/archive"
	"testing"
)

func Test_ArcMonitorStats(t *testing.T) {
	Convey("ArcMonitorStats", t, func() {
		//var o *archive.Archive
		o := &archive.Video{
			ID:    10018978,
			State: -30,
		}
		n := &archive.Video{
			ID:    10018978,
			State: -1,
		}
		err := s.hdlMonitorVideo(n, o)
		So(err, ShouldBeNil)
	})
}

func TestMonitorNotify(t *testing.T) {
	Convey("MonitorNotify", t, func() {
		tpl := s.email.MonitorNotifyTemplate("", "http://uat-manager.bilibili.co/#!/video/list?monitor_list=1_1_5", []string{"liusiming@bilibili.com"})
		s.email.PushToRedis(context.TODO(), tpl)
		s.email.Start("f_mail_list_fast")
	})
}
