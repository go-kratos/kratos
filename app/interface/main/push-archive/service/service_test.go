package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	time2 "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

var s *Service

func initd() {
	dir, _ := filepath.Abs("../cmd/push-archive-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = New(conf.Conf)
}

//普通关注的3个实验组比列均为5%，尾号分别为：00~04, 05~09, 10~14
func initd2() {
	dir, _ := filepath.Abs("../cmd/push-archive-test.toml")
	flag.Set("conf", dir)
	conf.Init()

	conf.Conf.ArcPush.UpperLimitExpire = time2.Duration(2 * time.Second)

	ps := []conf.Proportion{}
	for i := 0; i < 5; i++ {
		p := conf.Proportion{
			Proportion:          "0.05",
			ProportionStartFrom: fmt.Sprintf("%d", i*5),
		}
		ps = append(ps, p)
	}
	conf.Conf.ArcPush.Proportions = ps
	s = New(conf.Conf)
}

func Test_todaytime(t *testing.T) {
	initd()

	todayTime, err := s.getTodayTime("06:10:00")
	convey.Convey("当日定时时间", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
	t.Logf("todaytime(%s) unix(%d)", todayTime.Format("2006-01-02 15:04:05"), todayTime.Unix())
}

func Test_forbidTime(t *testing.T) {
	initd()

	t.Logf("the forbidtimes(%v)\n", s.ForbidTimes)
	forbid, err := s.isForbidTime()
	convey.Convey("禁止时间判断", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(forbid, convey.ShouldEqual, false)
	})
}

func Test_deadline(t *testing.T) {
	initd()

	deadline, err := s.getDeadline()
	convey.Convey("最近第3天的凌晨", t, func() {
		convey.So(err, convey.ShouldBeNil)
	})
	t.Logf("the deadline(%s) unix(%d)", deadline.Format("2006-01-02 15:04:05"), deadline.Unix())
}

func Test_proportion(t *testing.T) {
	initd2()

	fans := map[int64]int{
		100014: model.RelationAttention,
		100015: model.RelationAttention,
		100016: model.RelationAttention,
		100017: model.RelationAttention,
		100038: model.RelationAttention,
		100029: model.RelationAttention,
		100070: model.RelationAttention,
		100071: model.RelationAttention,
		100072: model.RelationAttention,
		100073: model.RelationSpecial,
		10034:  model.RelationSpecial,
		10038:  model.RelationSpecial,
	}
	expected := map[int]int{
		1: 4,
		2: 3,
	}

	attentions, specials := s.dao.FansByProportion(121231, fans)
	convey.Convey("根据比列过滤各组", t, func() {
		convey.So(len(attentions), convey.ShouldEqual, expected[1])
		convey.So(len(specials), convey.ShouldEqual, expected[2])
	})
}

func Test_Group(t *testing.T) {
	initd2()
	upper := int64(27515256)

	//upper主所有粉丝
	fans, err := s.dao.Fans(context.TODO(), upper, false)
	convey.Convey("获取所有粉丝", t, func() {
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(fans), convey.ShouldEqual, 37)
	})
	attentions, specials := s.dao.FansByProportion(upper, fans)
	t.Logf("attentions(%d), specials(%d)", len(attentions), len(specials))

	convey.Convey("没有hitby，则没有命中任何分组", t, func() {
		for _, g := range s.dao.FanGroups {
			g.Hitby = ""
		}
		list := s.group(upper, fans)
		convey.So(len(list), convey.ShouldEqual, 0)
	})

	convey.Convey("hitby=default，只命中2组", t, func() {
		s.dao.GroupOrder = []string{"1#attention", "1#ai:pushlist_play_recent", "2#special"}
		for _, g := range s.dao.FanGroups {
			g.Hitby = model.GroupDataTypeDefault
		}
		list := s.group(upper, fans)
		convey.So(len(list), convey.ShouldEqual, 3)
		convey.So(len(list["1#attention"]), convey.ShouldEqual, len(attentions))
		convey.So(len(list["1#ai:pushlist_play_recent"]), convey.ShouldEqual, 0)
		convey.So(len(list["2#special"]), convey.ShouldEqual, len(specials))
	})

	convey.Convey("hitby=hbase，根据表命中获取", t, func() {
		s.dao.GroupOrder = []string{"1#ai:pushlist_play_recent", "1#ai:pushlist_offline_up", "1#ai:pushlist_follow_recent", "1#attention", "2#special"}
		for _, g := range s.dao.FanGroups {
			if g.RelationType == 1 {
				g.Hitby = "hbase"
			}
		}
		list := s.group(upper, fans)
		t.Logf("list(%+v)", list)
	})

	//map[1#ai:pushlist_follow_recent:[]
	// 2#special:[21231134 4235023 1232032 2089809 88889018]
	// 1#ai:pushlist_play_recent:[27515303 27515311 27515317 27515300 27515401 27515306]
	// 1#ai:pushlist_offline_up:[]]
	expect := map[string]int{
		"1#ai:pushlist_follow_recent": 1,
		"1#ai:pushlist_play_recent":   5,
		"1#ai:pushlist_offline_up":    1,
		"2#special":                   5,
	}
	convey.Convey("实验组粉丝优先级分组", t, func() {
		hit := s.group(upper, fans)
		t.Logf("the hits(%v)", hit)
		for gkey, f := range hit {
			convey.So(len(f), convey.ShouldEqual, expect[gkey])
		}
	})
}

func Test_manytimes(t *testing.T) {
	initd2()

	upper := int64(27515256)
	mids1, len1, err1 := s.fans(upper, 1020, false)
	time.Sleep(time.Second * 2)
	mids2, len2, err2 := s.fans(upper, 1030, false)
	convey.Convey(" 同一up主多次获取过滤后的粉丝，第二次比第一次少", t, func() {
		convey.So(err1, convey.ShouldBeNil)
		convey.So(err2, convey.ShouldBeNil)
		convey.So(len1, convey.ShouldBeGreaterThan, 0)
		convey.So(len2, convey.ShouldBeGreaterThan, 0)
		convey.So(len2, convey.ShouldBeLessThanOrEqualTo, len1)
	})
	t.Logf("the mids1(%v)\n the mids2(%v)\n", mids1, mids2)
}
