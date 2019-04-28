package service

import (
	"container/ring"
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/job/main/sms/conf"
	"go-common/app/job/main/sms/model"
	smsmdl "go-common/app/service/main/sms/model"
	"go-common/library/net/trace"

	. "github.com/smartystreets/goconvey/convey"
)

var srv *Service

func init() {
	dir, _ := filepath.Abs("../cmd/sms-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func Test_ring(t *testing.T) {
	Convey("test ring", t, WithService(func(s *Service) {
		r := ring.New(3)
		r.Value = 0
		r = r.Next()
		r.Value = 1
		r = r.Next()
		r.Value = 2
		So(r.Len(), ShouldEqual, 3)
		for i := 0; i < 5; i++ {
			r = r.Next()
			t.Logf("%d", r.Value)
		}
	}))
}

func Test_Sms(t *testing.T) {
	Convey("sms", t, WithService(func(s *Service) {
		// http request会自动加trace header，不init trace的话,header value为空为会兴企reset
		trace.Init(s.c.Tracer)
		defer trace.Close()
		c := context.TODO()
		sl := &smsmdl.ModelSend{Mobile: "", Content: "您的账号正在哔哩哔哩2017动画角色人气大赏活动中进行领票操作，验证码为123456当日有效", Code: "whatever", Type: 1}
		p := s.smsp.Value.(model.Provider)
		_, err := p.SendSms(c, sl)
		So(err, ShouldBeNil)
		s.smsp.Ring = s.smsp.Next()
		p = s.smsp.Value.(model.Provider)
		_, err = p.SendSms(c, sl)
		So(err, ShouldBeNil)
		s.smsp.Ring = s.smsp.Next()
		p = s.smsp.Value.(model.Provider)
		_, err = p.SendSms(c, sl)
		So(err, ShouldBeNil)
		s.smsp.Ring = s.smsp.Next()
		p = s.smsp.Value.(model.Provider)
		_, err = p.SendSms(c, sl)
		So(err, ShouldBeNil)
	}))
}
