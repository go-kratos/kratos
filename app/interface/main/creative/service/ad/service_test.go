package ad

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	gM "go-common/app/interface/main/creative/model/game"
	"path/filepath"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/service"
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

func Test_GameList(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		localHost = "127.0.0.1"
		glist     *gM.ListWithPager
	)
	Convey("GameList", t, WithService(func(s *Service) {
		glist, err = s.GameList(c, "", "", 1, 20, localHost)
		So(err, ShouldBeNil)
		So(glist, ShouldNotBeNil)
		So(glist.Pn, ShouldEqual, 1)
		So(glist.Ps, ShouldEqual, 20)
		So(glist.List[0].GameBaseID, ShouldBeGreaterThanOrEqualTo, glist.List[1].GameBaseID)
	}))
}
