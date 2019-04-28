package Service

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"encoding/json"
	"go-common/app/job/live/wallet/conf"
	"go-common/app/job/live/wallet/model"

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

func TestMergeData(t *testing.T) {
	Convey("Cache update", t, func() {
		once.Do(startService)
		time.Sleep(time.Second)

		m := &model.Message{}
		user := &model.User{}
		m.New, _ = json.Marshal(user)
		m.Old, _ = json.Marshal(user)
		s.mergeData(m.New, m.Old, "update")
	})
}
