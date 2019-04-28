package appeal

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/appeal"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	s = New(conf.Conf, rpcdaos)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

var (
	mid = int64(27515256)
	cid = int64(123)
	aid = int64(10109136)
	tp  = "open"
	ip  = "127.0.0.1"
	c   = context.TODO()
)

func Test_List(t *testing.T) {
	Convey("List", t, WithService(func(s *Service) {
		all, open, closed, res, err := s.List(c, mid, 1, 10, tp, ip)
		So(err, ShouldBeNil)
		fmt.Println(all, open, closed, res)
	}))
}

func Test_Detail(t *testing.T) {
	Convey("Detail", t, WithService(func(s *Service) {
		res, err := s.Detail(c, mid, cid, ip)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func Test_State(t *testing.T) {
	Convey("State", t, WithService(func(s *Service) {
		err := s.State(c, mid, cid, 2, ip)
		So(err, ShouldBeNil)
	}))
}

func Test_Add(t *testing.T) {
	Convey("Add", t, WithService(func(s *Service) {
		qq := "2322"
		phone := "122333"
		email := "ddsds@qq.com"
		desc := "dedsds"
		attachments := "sddds"
		ap := &appeal.BusinessAppeal{}
		res, err := s.Add(c, mid, aid, qq, phone, email, desc, attachments, ip, ap)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func Test_Reply(t *testing.T) {
	Convey("Reply", t, WithService(func(s *Service) {
		event := int64(2)
		content := "122333"
		attachments := "sddds"
		err := s.Reply(c, mid, cid, event, content, attachments, ip)
		So(err, ShouldBeNil)
	}))
}

func Test_PhoneEmail(t *testing.T) {
	Convey("PhoneEmail", t, WithService(func(s *Service) {
		ck := ""
		res, err := s.PhoneEmail(c, ck, ip)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func Test_Star(t *testing.T) {
	Convey("Star", t, WithService(func(s *Service) {
		star := int64(2)
		err := s.Star(c, mid, cid, star, ip)
		So(err, ShouldBeNil)
	}))
}
