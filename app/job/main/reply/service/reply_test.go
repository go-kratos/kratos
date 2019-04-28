package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFilterViolationMsg(t *testing.T) {
	Convey("TestFilterViolationMsg", t, func() {
		res := filterViolationMsg("123456789评论过虑违规内容评论过虑违规内容")
		t.Log(res)
	})
}

func TestTopicReg(t *testing.T) {
	s := Service{}
	c := context.Background()
	Convey("TestAtReg", t, func() {
		topics := s.regTopic(c, "#你懂 得##222#")
		So(len(topics), ShouldEqual, 2)
		So(topics[0], ShouldEqual, "你懂 得")
		So(topics[1], ShouldEqual, "222")
		topics = s.regTopic(c, "#你懂 \n得##22@有人艾特2#")
		So(len(topics), ShouldEqual, 0)
		topics = s.regTopic(c, "#你懂 \n得#哈哈哈#22@有人艾特2#")
		So(len(topics), ShouldEqual, 1)
		So(topics[0], ShouldEqual, "哈哈哈")
		topics = s.regTopic(c, "#  ##	##你懂得")
		So(len(topics), ShouldEqual, 0)
		topics = s.regTopic(c, "热热# ##！%……&（）（）*（）*（）&*……&……%……￥%##同一套##协助特大号哈哈哈嘎嘎协助特大号哈哈哈嘎嘎协助特大号哈哈哈ee120##协助特大号哈哈哈嘎嘎协助特大号哈哈哈嘎嘎协助特大号哈哈哈ee12##@1r##tet##899##5677#")
		So(len(topics), ShouldEqual, 5)
		topics = s.regTopic(c, "#我是大佬你是谁你是大佬嘛哈哈啊#123#")
		So(len(topics), ShouldEqual, 1)
		topics = s.regTopic(c, "#2😁3#123#3😁3##2😁3#")
		So(len(topics), ShouldEqual, 1)
		So(topics[0], ShouldEqual, "123")
		topics = s.regTopic(c, " http://t.bilibili.com/av111111#reply#haha #didi")
		So(len(topics), ShouldEqual, 0)
		topics = s.regTopic(c, " http://t.bilibili.com/av111111#reply#haha #didi# http://t.baidu.com/av111111#reply#haha")
		So(len(topics), ShouldEqual, 2)
		So(topics[0], ShouldEqual, "didi")
		So(topics[1], ShouldEqual, "reply")
		topics = s.regTopic(c, "asdasd#av1000#33333#vc11111#44444#cv1111#55555#")
		So(len(topics), ShouldEqual, 3)

	})
}

func TestAtReg(t *testing.T) {
	Convey("TestAtReg", t, func() {
		ss := _atReg.FindAllStringSubmatch("@aa:hh@bb,cc", 10)
		So(len(ss), ShouldEqual, 2)
		So(ss[0][1], ShouldEqual, "aa")
		So(ss[1][1], ShouldEqual, "bb")
		ss = _atReg.FindAllStringSubmatch("@aa@bb", 10)
		So(len(ss), ShouldEqual, 2)
		So(ss[0][1], ShouldEqual, "aa")
		So(ss[1][1], ShouldEqual, "bb")
		ss = _atReg.FindAllStringSubmatch("@aa  @bb", 10)
		So(len(ss), ShouldEqual, 2)
		So(ss[0][1], ShouldEqual, "aa")
		So(ss[1][1], ShouldEqual, "bb")
		ss = _atReg.FindAllStringSubmatch("@aa  bb@cc;@dd:sa", 10)
		So(len(ss), ShouldEqual, 3)
		So(ss[0][1], ShouldEqual, "aa")
		So(ss[1][1], ShouldEqual, "cc;")
		So(ss[2][1], ShouldEqual, "dd")
	})
}
