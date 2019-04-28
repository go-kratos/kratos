package service

import (
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go-common/app/job/live/dao-anchor-job/internal/conf"
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}
func TestMinuteDataToCacheList(t *testing.T) {
	Convey("testMinuteDataToCacheList", t, func() {
		s.minuteDataToCacheList()
	})

}
func TestMinuteDataToDB(t *testing.T) {
	Convey("TestMinuteDataToDB", t, func() {
		s.minuteDataToDB()
	})

}
