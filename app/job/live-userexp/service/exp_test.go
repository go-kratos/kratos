package service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"encoding/json"
	"go-common/app/job/live-userexp/conf"
	"go-common/app/job/live-userexp/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = New(conf.Conf)
}

func TestLevelCacheUpdate(t *testing.T) {
	Convey("Cache update", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)

		m := &model.Message{}
		exp := &model.Exp{}
		m.New, _ = json.Marshal(exp)
		m.Old, _ = json.Marshal(exp)
		s.levelCacheUpdate(m.New, m.Old)
	})
}
