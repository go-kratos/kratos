package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/creative/conf"
	accapi "go-common/app/service/main/account/api"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/creative-admin.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_Profile(t *testing.T) {
	mid := int64(27515256)
	Convey("Profile", t, WithService(func(s *Service) {
		var (
			pfl *accapi.ProfileStatReply
			err error
		)
		pfl, err = s.ProfileStat(context.TODO(), mid)
		time.Sleep(time.Millisecond * 100)
		So(err, ShouldBeNil)
		So(pfl, ShouldNotBeNil)
	}))
}

func Test_SearchKeywords(t *testing.T) {
	Convey("SearchKeywords", t, WithService(func(s *Service) {
		s.SearchKeywords()
	}))
}
