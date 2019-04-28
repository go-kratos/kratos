package service

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/videoup/conf"
	"testing"
)

func TestTypeTopParent(t *testing.T) {
	err := conf.Init()
	if err != nil {
		return
	}
	s := New(conf.Conf)
	Convey("test TypeTopParent", t, func() {
		_, err := s.TypeTopParent(int16(1808))
		So(err, ShouldNotBeNil)
	})
}

// TestArchiveRound 测试商单稿件round
func TestPorderArchiveRound(t *testing.T) {
	var (
		c              = context.TODO()
		aid      int64 = 5464730 //稿件id
		mid      int64 = 254386  //up主id
		typeID   int16 = 22      //分区id
		nowRound int8  = 10      //二审提交
		newState int8  = -40     //定时发布
		resRound       = 21      //最终返回的round结果
	)
	err := conf.Init()
	if err != nil {
		return
	}
	s := New(conf.Conf)
	Convey("test TestPorderArchiveRound", t, func() {
		round := s.archiveRound(c, nil, aid, mid, typeID, nowRound, newState, false)
		//round == 21
		So(round, ShouldEqual, resRound)
	})
}

func TestStringHandler(t *testing.T) {
	var res string
	delimiter := ","
	s1 := "t1,t2"
	s2 := "t1"
	s3 := "t2"
	s4 := "t3"
	s5 := "t1,t2,t3"
	s6 := "t1,t3,t4,t5"
	s7 := "t1,t2,t3,t4,t5"

	Convey("StringHandler", t, func() {
		//增删空字符串
		res = StringHandler(s1, "", delimiter, false)
		So(res, ShouldEqual, s1)
		res = StringHandler(s1, "", delimiter, true)
		So(res, ShouldEqual, s1)

		//增删重复字符串
		res = StringHandler(s1, s2, delimiter, false)
		So(res, ShouldEqual, s1)
		res = StringHandler(s1, s2, delimiter, true)
		So(res, ShouldEqual, s3)

		//增删不重复字符串
		res = StringHandler(s1, s4, delimiter, false)
		So(res, ShouldEqual, s5)
		res = StringHandler(s1, s4, delimiter, true)
		So(res, ShouldEqual, s1)

		//增删多个重复，且多个不重复字符串
		res = StringHandler(s5, s6, delimiter, false)
		So(res, ShouldEqual, s7)
		res = StringHandler(s5, s6, delimiter, true)
		So(res, ShouldEqual, s3)
	})
}

// TestSplitInts
func TestSplitInts(t *testing.T) {
	var (
		str = " 123,334343\n,\t1\r11"
	)
	err := conf.Init()
	if err != nil {
		return
	}
	s := New(conf.Conf)
	Convey("test TestSplitInts", t, func() {
		ids, err := s.SplitInts(str)
		fmt.Print(ids)
		So(ids, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
