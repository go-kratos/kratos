package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/spy/conf"
	"go-common/app/admin/main/spy/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func init() {
	var err error
	dir, _ := filepath.Abs("../cmd/spy-admin-test.toml")
	flag.Set("conf", dir)
	if err = conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}

func TestSetting(t *testing.T) {
	Convey("test checkSettingVal", t, func() {
		var err = s.checkSettingVal(model.AutoBlock, "abc")
		So(err, ShouldEqual, ecode.SpySettingValTypeError)
		err = s.checkSettingVal(model.AutoBlock, "3")
		So(err, ShouldEqual, ecode.SpySettingValueOutOfRange)
		err = s.checkSettingVal(model.AutoBlock, "1")
		So(err, ShouldBeNil)
		err = s.checkSettingVal(model.LessBlockScore, "abc")
		So(err, ShouldEqual, ecode.SpySettingValTypeError)
		err = s.checkSettingVal(model.LessBlockScore, "45")
		So(err, ShouldEqual, ecode.SpySettingValueOutOfRange)
		err = s.checkSettingVal(model.LessBlockScore, "1")
		So(err, ShouldBeNil)
		err = s.checkSettingVal(model.LimitBlockCount, "abc")
		So(err, ShouldEqual, ecode.SpySettingValTypeError)
		err = s.checkSettingVal(model.LimitBlockCount, "-2")
		So(err, ShouldEqual, ecode.SpySettingValueOutOfRange)
		err = s.checkSettingVal(model.LimitBlockCount, "20000")
		So(err, ShouldBeNil)
		err = s.checkSettingVal("unknown prop", "abc")
		So(err, ShouldEqual, ecode.SpySettingUnknown)
	})
	Convey("test Setting", t, func() {
		var (
			list []*model.Setting
			err  error
		)
		list, err = s.SettingList(c)
		So(err, ShouldBeNil)
		So(list, ShouldNotBeEmpty)
		var (
			name = "go test"
			prop = list[0].Property
			val  = list[0].Val
		)
		err = s.UpdateSetting(c, name, prop, val)
		So(err, ShouldBeNil)
		err = s.UpdateSetting(c, name, model.LessBlockScore, "100")
		So(err, ShouldEqual, ecode.SpySettingValueOutOfRange)
	})
}

// go test  -test.v -test.run TestStat
func TestStat(t *testing.T) {
	var (
		state    int8  = 1
		id       int64 = 3
		isdel    int8  = 1
		tid      int64 = 1
		tmid     int64 = 1
		ty       int8  = 1
		count    int64 = 10
		operater       = "admin"
		pn             = 1
		ps             = 8
	)
	Convey("test UpdateState", t, func() {
		err := s.UpdateState(c, state, id, operater)
		So(err, ShouldBeNil)
	})
	Convey("test UpdateStatQuantity", t, func() {
		err := s.UpdateStatQuantity(c, count, id, operater)
		So(err, ShouldBeNil)
	})
	Convey("test DeleteStat", t, func() {
		err := s.DeleteStat(c, isdel, id, operater)
		So(err, ShouldBeNil)
	})
	Convey("test StatPage", t, func() {
		page, err := s.StatPage(c, tmid, tid, ty, pn, ps)
		So(err, ShouldBeNil)
		fmt.Println("page", page)
	})
}
