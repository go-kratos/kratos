package service

import (
	"flag"
	"math"
	"path/filepath"
	"sync"
	"testing"

	"go-common/app/service/main/push-strategy/conf"
	pushmdl "go-common/app/service/main/push/model"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	dir, _ := filepath.Abs("../cmd/push-strategy-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func TestCheckMidBySetting(t *testing.T) {
	var (
		mutex sync.Mutex
		biz   = s.c.BizID.Archive
		mid   = int64(91221505)
	)
	Convey("checkMidBySetting", t, func() {
		Convey("without user setting", func() {
			mutex.Lock()
			s.settings = make(map[int64]map[int]int)
			res := s.checkMidBySetting(biz, mid)
			So(res, ShouldBeTrue)
			mutex.Unlock()
		})
		Convey("not in limited business list", func() {
			mutex.Lock()
			s.settings = map[int64]map[int]int{mid: make(map[int]int)}
			res := s.checkMidBySetting(math.MaxInt32, mid)
			So(res, ShouldBeTrue)
			mutex.Unlock()
		})
		Convey("business no setting, should be pass", func() {
			mutex.Lock()
			s.settings = map[int64]map[int]int{mid: make(map[int]int)}
			res := s.checkMidBySetting(biz, mid)
			So(res, ShouldBeTrue)
			mutex.Unlock()
		})
		Convey("switch on, should be pass", func() {
			mutex.Lock()
			s.settings = map[int64]map[int]int{mid: {pushmdl.UserSettingArchive: pushmdl.SwitchOn}}
			res := s.checkMidBySetting(biz, mid)
			So(res, ShouldBeTrue)
			mutex.Unlock()
		})
		Convey("switch off, should not be pass", func() {
			mutex.Lock()
			s.settings = map[int64]map[int]int{mid: {pushmdl.UserSettingArchive: pushmdl.SwitchOff}}
			res := s.checkMidBySetting(biz, mid)
			So(res, ShouldBeFalse)
			mutex.Unlock()
		})
	})
}
