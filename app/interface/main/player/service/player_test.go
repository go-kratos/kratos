package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/player/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var svf *Service

func WithService(f func(s *Service)) func() {
	return func() {
		dir, _ := filepath.Abs("../cmd/player-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if svf == nil {
			svf = New(conf.Conf)
		}
		time.Sleep(2 * time.Second)
		f(svf)
	}
}

func TestMidCrc(t *testing.T) {
	expects := map[int64]string{
		8167601: "2425c296",
		123456:  "0972d361",
	}
	for mid, crc := range expects {
		if midCrc(mid) != crc {
			t.Errorf("crc %v expect %s got %s", mid, crc, midCrc(mid))
		}
	}
}

func TestService_Carousel(t *testing.T) {
	Convey("carousel len should > 0", t, WithService(func(svf *Service) {
		carousel, err := svf.Carousel(context.Background())
		So(err, ShouldBeNil)
		So(len(carousel), ShouldBeGreaterThan, 0)
	}))
}

func TestService_Player(t *testing.T) {
	Convey("player should return without err", t, WithService(func(svf *Service) {
		player, err := svf.Player(context.Background(), 0, 10097666, 10108404, "127.0.0.1", "", time.Now())
		So(err, ShouldBeNil)
		So(len(player), ShouldBeGreaterThan, 0)
	}))
	Convey("player should return without err", t, WithService(func(svf *Service) {
		player, err := svf.Player(context.Background(), 0, 10010666, 10108404, "127.0.0.1", "", time.Now())
		So(err, ShouldBeNil)
		So(len(player), ShouldBeGreaterThan, 0)
	}))
}
