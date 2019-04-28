package service

import (
	"flag"
	"go-common/app/interface/main/videoup/conf"
	"path/filepath"
	"time"

	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"testing"
	"unicode/utf8"
)

var (
	s *Service
)

func init() {
	dir, _ := filepath.Abs("../cmd/videoup.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		Reset(func() {})
		f(s)
	}
}

func Test_checkAddStaff(t *testing.T) {
	var (
		c         = context.Background()
		mid int64 = 222
		sf1       = &archive.Staff{
			Title: "a",
			Mid:   123,
		}
		sf2 = &archive.Staff{
			Title: "ab",
			Mid:   124,
		}
		ap = &archive.ArcParam{
			Copyright: archive.CopyrightCopy,
			Staffs:    []*archive.Staff{sf1, sf2},
		}
		err error
	)
	Convey("checkAddStaff", t, WithService(func(s *Service) {
		err = s.checkAddStaff(c, ap, mid, "")
		So(err, ShouldEqual, ecode.VideoupStaffCopyright)
	}))
}

func Test_getStaffChanges(t *testing.T) {
	var (
		sf1 = &archive.Staff{
			Title: "a",
			Mid:   123,
		}
		sf2 = &archive.Staff{
			Title: "ab",
			Mid:   124,
		}
		sf3 = &archive.Staff{
			Title: "abv",
			Mid:   124,
		}
		os = []*archive.Staff{sf1, sf2}
		ns = []*archive.Staff{sf1, sf3}
	)
	Convey("checkAddStaff", t, WithService(func(s *Service) {
		changes, _ := s.getStaffChanges(os, ns)
		So(changes, ShouldNotBeNil)
	}))
}

func Test_CheckStaffReg(t *testing.T) {
	var err error
	v := &archive.Staff{
		Title: "配音",
		Mid:   123,
	}
	Convey("CheckStaffReg", t, WithService(func(s *Service) {
		tl := utf8.RuneCountInString(v.Title)
		if tl > 4 || tl == 0 {
			err = ecode.VideoupStaffTitleLength
		}
		if !_staffNameReg.MatchString(v.Title) {
			err = ecode.VideoupStaffTitleLength
		}
		So(err, ShouldBeNil)
	}))
}
