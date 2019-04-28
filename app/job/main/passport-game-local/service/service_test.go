package service

import (
	"sync"
	"testing"

	"go-common/app/job/main/passport-game-local/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func TestNew(t *testing.T) {
	once.Do(startService)
	Convey("ping should ok", t, func() {
		So(s.Ping, ShouldBeNil)
	})
}
