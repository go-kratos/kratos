package service

import (
	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_groupparam(t *testing.T) {
	initd()

	expect := map[string]string{
		"1#ai:pushlist_follow_recent": "follow",
		"1#ai:pushlist_play_recent":   "play",
		"1#ai:pushlist_offline_up":    "offline",
		"2#special":                   "special",
	}
	convey.Convey("推送的group参数", t, func() {
		for k, g := range s.dao.FanGroups {
			group := s.getGroupParam(g)
			convey.So(group, convey.ShouldEqual, expect[k])
		}
	})
}

func Test_usersettingfilter(t *testing.T) {
	initd()

	mid := int64(11111111)
	s.userSettings[mid] = &model.Setting{Type: model.PushTypeForbid}
	allow := s.filterUserSetting(mid, model.RelationSpecial)
	convey.Convey("usersettings filter关闭开关，则排除", t, func() {
		convey.So(allow, convey.ShouldEqual, false)
	})

	s.userSettings[mid] = nil
	allow = s.filterUserSetting(mid, model.RelationSpecial)
	convey.Convey("usersettings filter未设置开关，则允许", t, func() {
		convey.So(allow, convey.ShouldEqual, true)
	})

	s.userSettings[mid] = &model.Setting{Type: model.PushTypeAttention}
	allow = s.filterUserSetting(mid, model.RelationSpecial)
	convey.Convey("usersettings filter设置未所有关注，则允许", t, func() {
		convey.So(allow, convey.ShouldEqual, true)
	})
}

func Test_ispgc(t *testing.T) {
	arc := new(model.Archive)
	convey.Convey("pgc稿件判断", t, func() {
		arc.Attribute = int32(110336)
		convey.So(s.isPGC(arc), convey.ShouldEqual, true)

		arc.Attribute = int32(16512)
		convey.So(s.isPGC(arc), convey.ShouldEqual, false)
	})
}

func TestServicefansByAbtest(t *testing.T) {
	initd()
	group := &dao.FanGroup{
		Hitby:       "ab_test",
		HBaseTable:  "push_archive_ab_test",
		HBaseFamily: []string{"cf"},
	}
	fans := []int64{1, 2, 3, 4, 5, 6}
	convey.Convey("fansByAbtest", t, func() {
		exists, notExists := s.fansByAbtest(group, fans)
		t.Logf("exists(%v)", exists)
		t.Logf("notExists(%v)", notExists)
	})
}
