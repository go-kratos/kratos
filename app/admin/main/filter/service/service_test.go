package service

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s   *Service
	ctx = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	defer s.Close()
	os.Exit(m.Run())
}

func TestCheck(t *testing.T) {
	Convey("reg", t, func() {
		err := s.checkReg(model.RegMode, "abc.{0,2}cde")
		So(err, ShouldBeNil)
		err = s.checkReg(model.RegMode, "abc.*cde")
		So(err, ShouldEqual, ecode.FilterRegexpError1)
		err = s.checkReg(model.RegMode, "(abc||cde)")
		So(err, ShouldEqual, ecode.FilterRegexpError2)
		err = s.checkReg(model.RegMode, "(dasfdsa")
		So(err, ShouldEqual, ecode.FilterIllegalRegexp)
	})

	Convey("sample check", t, func() {
		err := s.checkBlackSample(model.RegMode, ".*")
		So(err, ShouldEqual, ecode.FilterBlackSampleHit)
		err = s.checkBlackSample(model.RegMode, "test.{1,2}test")
		So(err, ShouldBeNil)

		err = s.checkWhiteSample(model.RegMode, ".*")
		So(err, ShouldEqual, ecode.FilterWhiteSampleHit)
		err = s.checkWhiteSample(model.RegMode, "test.{1,2}test")
		So(err, ShouldBeNil)
	})
}

func TestFilter(t *testing.T) {
	Convey("filter rule admin", t, func() {
		var (
			areas           = []string{"common"}
			rules           = []string{"test134"}
			tpIDs           = []int64{0}
			adminName       = "muyang"
			adminID   int64 = 233
			stime           = time.Now()
			etime           = stime.Add(time.Hour)
			level           = &model.AreaLevel{
				Level: 20,
				Area: map[string]int8{
					"common": 30,
				},
			}
			filterInfo *model.FilterInfo
		)
		err := s.AdminAdd(ctx, areas, rules, level, "test comment", adminName, model.RegMode, tpIDs, adminID, stime.Unix(), etime.Unix(), 0, 0)
		So(err, ShouldBeNil)

		filterInfos, count, err := s.AdminSearch(ctx, "test134", "common", "", "", 30, 0, 0, 1, 1)
		So(err, ShouldBeNil)
		So(filterInfos, ShouldNotBeEmpty)
		So(count, ShouldEqual, 1)
		So(filterInfos[0].ID, ShouldBeGreaterThan, 0)

		filterInfo, err = s.AdminRuleByID(ctx, filterInfos[0].ID)
		So(err, ShouldBeNil)
		So(filterInfo, ShouldNotBeNil)

		err = s.AdminEdit(ctx, areas, "test321", "test comment", "test reason", adminName, model.RegMode, level, tpIDs, filterInfo.ID, adminID, stime.Unix(), etime.Unix(), 0, 0)
		So(err, ShouldBeNil)

		logs, err := s.AdminLog(ctx, filterInfo.ID)
		So(err, ShouldBeNil)
		So(logs, ShouldNotBeEmpty)

		err = s.AdminDel(ctx, filterInfos[0].ID, adminID, "test delete", adminName)
		So(err, ShouldBeNil)
	})
}

func TestWhite(t *testing.T) {
	Convey("white admin", t, func() {
		var (
			areas           = []string{"common"}
			content         = "test_white_233"
			tpIDs           = []int64{0}
			adminName       = "muyang"
			adminID   int64 = 233
		)
		err := s.AddAreaWhite(ctx, content, model.RegMode, areas, tpIDs, adminID, adminName, "test comment")
		So(err, ShouldBeNil)

		whiteInofs, count, err := s.SearchWhite(ctx, content, areas[0], 1, 1)
		So(err, ShouldBeNil)
		So(whiteInofs, ShouldNotBeEmpty)
		So(count, ShouldEqual, 1)
		So(whiteInofs[0].ID, ShouldBeGreaterThan, 0)

		whiteInfo, err := s.WhiteInfo(ctx, whiteInofs[0].ID)
		So(err, ShouldBeNil)
		So(whiteInfo, ShouldNotBeNil)

		logs, err := s.WhiteEditLog(ctx, whiteInfo.ID)
		So(err, ShouldBeNil)
		So(logs, ShouldNotBeEmpty)

		err = s.DeleteWhite(ctx, whiteInfo.ID, adminID, adminName, "test delte")
		So(err, ShouldBeNil)
	})
}

func TestKey(t *testing.T) {
	Convey("key admin", t, func() {
		var (
			areas           = []string{"common"}
			content         = "test_key_233"
			key             = "test:233"
			adminName       = "muyang"
			adminID   int64 = 233
			stime           = time.Now()
			etime           = stime.Add(time.Hour)
		)
		err := s.AddKey(ctx, areas, key, content, "test comment", adminName, model.RegMode, 20, adminID, stime.Unix(), etime.Unix())
		So(err, ShouldBeNil)

		count, keyInfos, err := s.SearchKey(ctx, key, "test comment", 1, 1, 0)
		So(err, ShouldBeNil)
		So(keyInfos, ShouldNotBeEmpty)
		So(count, ShouldEqual, 1)
		So(keyInfos[0].ID, ShouldBeGreaterThan, 0)

		err = s.EditKey(ctx, key, keyInfos[0].ID, areas, model.StrMode, "test_key_322", 30, stime.Unix(), etime.Unix(), adminID, adminName, "test edit", "test edit")
		So(err, ShouldBeNil)

		err = s.DelKeyFid(ctx, key, keyInfos[0].ID, adminID, "test delete", adminName, "test delete")
		So(err, ShouldBeNil)
	})
}
