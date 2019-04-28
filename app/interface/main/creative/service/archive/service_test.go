package archive

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"

	"go-common/app/interface/main/creative/conf"
	arcmdl "go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/order"

	"go-common/app/interface/main/creative/service"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	p *service.Public
)

func init() {
	dir, _ := filepath.Abs("../../cmd/creative.toml")
	flag.Set("conf", dir)
	conf.Init()
	rpcdaos := service.NewRPCDaos(conf.Conf)
	p = service.New(conf.Conf, rpcdaos)
	s = New(conf.Conf, rpcdaos, p)
	time.Sleep(time.Second)
}
func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func ctx() context.Context {
	return context.Background()
}

func Test_Oasis(t *testing.T) {
	Convey("should get user oasis info", t, func() {
		res, err := s.Oasis(ctx(), 27515256, "127.0.0.1")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
func Test_Order(t *testing.T) {
	var (
		c      = ctx()
		err    error
		MID    = int64(27515256)
		orders []*order.Order
		oasis  *order.Oasis
	)
	Convey("ExecuteOrders", t, WithService(func(s *Service) {
		orders, err = s.ExecuteOrders(c, MID, "127.0.0.1")
		So(err, ShouldBeNil)
		spew.Dump(orders)
	}))
	Convey("Oasis", t, WithService(func(s *Service) {
		oasis, err = s.Oasis(c, MID, "127.0.0.1")
		So(err, ShouldBeNil)
		spew.Dump(oasis)
	}))
}

func Test_ServiceBasic(t *testing.T) {
	var (
		c   = context.Background()
		err error
	)
	Convey("ServiceClose", t, WithService(func(s *Service) {
		s.Close()
	}))
	Convey("Ping", t, WithService(func(s *Service) {
		err = s.Ping(c)
		So(err, ShouldBeNil)
	}))
	Convey("AppModuleShowMap", t, WithService(func(s *Service) {
		var res map[string]bool
		mid := int64(91513044)
		res = s.AppModuleShowMap(mid, false)
		So(res, ShouldNotBeNil)
	}))
}

func Test_NetSafe(t *testing.T) {
	var (
		c   = context.Background()
		err error
		nid int64
		md5 string
	)
	nid = 123
	md5 = "iamamd5string"
	Convey("AddNetSafeMd5", t, WithService(func(s *Service) {
		err = s.AddNetSafeMd5(c, nid, md5)
		So(err, ShouldBeNil)
	}))
	Convey("NotifyNetSafe", t, WithService(func(s *Service) {
		err = s.NotifyNetSafe(c, nid)
		So(err, ShouldBeNil)
	}))
}

func Test_DescFormat(t *testing.T) {
	var (
		c                 = context.Background()
		err               error
		typeid, copyright int64
		langStr, ip       string
		desc              *arcmdl.DescFormat
		af                []*arcmdl.AppFormat
		length            int
	)
	Convey("DescFormat", t, WithService(func(s *Service) {
		desc, err = s.DescFormat(c, typeid, copyright, langStr, ip)
		So(err, ShouldBeNil)
		So(desc, ShouldNotBeNil)
	}))
	Convey("DescFormatForApp", t, WithService(func(s *Service) {
		desc, length, err = s.DescFormatForApp(c, typeid, copyright, langStr, ip)
		So(err, ShouldBeNil)
		So(desc, ShouldNotBeNil)
		So(length, ShouldNotBeNil)
		So(length, ShouldBeGreaterThanOrEqualTo, 0)
	}))
	Convey("AppFormats", t, WithService(func(s *Service) {
		af, err = s.AppFormats(c)
		So(err, ShouldBeNil)
		So(af, ShouldNotBeNil)
		So(len(af), ShouldBeGreaterThanOrEqualTo, 0)
	}))
}
func Test_Arc(t *testing.T) {
	var (
		c          = context.Background()
		err        error
		mid, aid   int64
		ak, ck, ip = "", "", ""
		ap         *arcmdl.SimpleArchiveVideos
	)
	Convey("SimpleArchiveVideos", t, WithService(func(s *Service) {
		ap, err = s.SimpleArchiveVideos(c, mid, aid, ak, ck, ip)
		So(err, ShouldBeNil)
		So(ap, ShouldNotBeNil)
	}))
}

func TestArchiveBIZsByTime(t *testing.T) {
	Convey("BIZsByTime", t, WithService(func(s *Service) {
		var (
			c     = context.Background()
			start = time.Unix(1544000000, 0)
			end   = time.Unix(1544008000, 0)
			tp    = int8(2)
		)
		// mock bizs
		Convey("When everything gose positive", WithService(func(s *Service) {
			bizs, err := s.BIZsByTime(c, &start, &end, tp)
			Convey("Then err should be nil.bizs should not be nil.", func(ctx C) {
				So(err, ShouldBeNil)
				So(bizs, ShouldNotBeNil)
				fmt.Println(bizs[0])
			})
		}))
	}))
}
