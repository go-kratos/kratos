package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/service/main/favorite/conf"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/favorite-service-test.toml")
	flag.Set("conf", dir)
	err := conf.Init()
	if err != nil {
		fmt.Printf("conf.Init() error(%v)", err)
	}
	s = New(conf.Conf)
}

func Test_filter(t *testing.T) {
	Convey("filter", t, func() {
		var (
			name = "枪支迷药"
		)
		err := s.filter(context.TODO(), name)
		t.Logf("err:%v", err)
		So(err, ShouldEqual, ecode.FavFolderBanned)
	})
}
