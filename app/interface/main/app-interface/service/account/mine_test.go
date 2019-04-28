package space

import (
	"context"
	"flag"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var s *Service

func init() {
	dir, _ := filepath.Abs("../../cmd/app-interface-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(3 * time.Second)
}

func Test_Mine(t *testing.T) {
	Convey("Test_Mine", t, func() {
		mine, err := s.Mine(context.Background(), 16840123, "ios", "", 0, 0)
		So(err, ShouldBeNil)
		for _, s := range mine.Sections {
			Printf("%+v\n", s.Title)
			for _, item := range s.Items {
				Printf("%+v\n", item.Title)
			}
		}
	})
}

func Test_MineIpad(t *testing.T) {
	Convey("Test_Mine", t, func() {
		mine, err := s.MineIpad(context.Background(), 16840123, "ios", "", 0, 2)
		So(err, ShouldBeNil)
		for _, item := range mine.IpadSections {
			Printf("%+v\n", item.Title)
		}
	})
}

func Test_MyInfo(t *testing.T) {
	Convey("Test_MyInfo", t, func() {
		myinfo, err := s.Myinfo(context.Background(), 1684013)
		So(err, ShouldBeNil)
		Printf("%+v", myinfo)
	})
}

func Test_LoadSidebar(t *testing.T) {
	Convey("Test_LoadSidebar", t, func() {
		Println("starting...")
		sections := s.sections(context.Background(), nil, nil, 1, 1, true, 1)
		for _, section := range sections {
			Println("section标题:" + section.Title)
			for _, item := range section.Items {
				Println(item.Title)
			}
		}
	})
}
