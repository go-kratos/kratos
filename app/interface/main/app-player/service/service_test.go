package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"go-common/app/interface/main/app-player/conf"
	"go-common/app/interface/main/app-player/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/app-player-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

func Test_Playurl(t *testing.T) {
	var params = &model.Param{
		AID:     1,
		MobiApp: "android",
		CID:     1,
		Qn:      32,
	}
	Convey("Test_Playurl", t, func() {
		s.Playurl(context.TODO(), 1, params, 1, "dajskldasjkl", "")
	})
}
