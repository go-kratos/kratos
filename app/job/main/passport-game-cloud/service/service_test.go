package service

import (
	"sync"
	"testing"

	"go-common/app/job/main/passport-game-cloud/conf"

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
	Convey("game app ids should ok", t, func() {
		So(s.gameAppIDs[0], ShouldEqual, _gameAppID)

		t.Logf("s.gameAppIDs: %v", s.gameAppIDs)
	})
}
