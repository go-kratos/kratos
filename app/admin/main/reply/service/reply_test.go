package service

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/reply/conf"
	"go-common/app/admin/main/reply/model"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReply(t *testing.T) {
	c := context.Background()
	Convey("reply del and recover", t, WithService(func(s *Service) {
		sub, err := s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		rp, err := s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		err = s.adminDeleteReply(c, 23, []int64{5464686, 5464686}, []int64{894733824, 894733816}, 0, 1, 0, true, "test", "test", 0, 0)
		So(err, ShouldBeNil)
		rp2, err := s.reply(c, 5464686, 894733824)
		So(err, ShouldBeNil)
		So(rp2.State, ShouldEqual, 3)
		rp2, err = s.reply(c, 5464686, 894733816)
		So(err, ShouldBeNil)
		So(rp2.State, ShouldEqual, 3)
		rp2, err = s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		So((rp2.RCount - rp.RCount), ShouldEqual, -2)
		sub2, err := s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So((sub2.ACount - sub.ACount), ShouldEqual, -2)
		err = s.AdminRecoverReply(c, 23, "asd", 5464686, 894733824, 1, "test")
		So(err, ShouldBeNil)
		err = s.AdminRecoverReply(c, 23, "sad", 5464686, 894733816, 1, "test")
		So(err, ShouldBeNil)
		rp2, err = s.reply(c, 5464686, 894733824)
		So(err, ShouldBeNil)
		So(rp2.State, ShouldEqual, 0)
		rp2, err = s.reply(c, 5464686, 894733816)
		So(err, ShouldBeNil)
		So(rp2.State, ShouldEqual, 0)
		rp2, err = s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		So((rp2.RCount - rp.RCount), ShouldEqual, 0)
		sub2, err = s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So((sub2.ACount - sub.ACount), ShouldEqual, 0)
		alog, err := s.dao.AdminLog(c, 894733824)
		So(err, ShouldBeNil)
		So(alog.State, ShouldEqual, 3)
		alog, err = s.dao.AdminLog(c, 894733816)
		So(err, ShouldBeNil)
		So(alog.State, ShouldEqual, 3)
	}))

	Convey("reply top", t, WithService(func(s *Service) {
		err := s.AddTop(c, 23, "asd", 5464686, 894717392, 1, 1)
		So(err, ShouldBeNil)
		rp, err := s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		So(rp.IsTop(), ShouldBeTrue)
		sub, err := s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So(sub.AttrVal(model.SubAttrTopAdmin), ShouldEqual, model.AttrYes)
		err = s.AddTop(c, 23, "asd", 5464686, 894717392, 1, 0)
		So(err, ShouldBeNil)
		rp, err = s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		So(rp.IsTop(), ShouldBeFalse)
		sub, err = s.subject(c, 5464686, 1)
		So(err, ShouldBeNil)
		So(sub.AttrVal(model.SubAttrTopAdmin), ShouldEqual, model.AttrNo)
	}))

	Convey("reply pass", t, WithService(func(s *Service) {
		s.dao.UpdateReplyState(c, 5464686, 894717384, model.StatePending)
		err := s.adminPassReply(c, 23, "sds", []int64{5464686}, []int64{894717384}, 1, "test")
		So(err, ShouldBeNil)
		rp, err := s.reply(c, 5464686, 894717392)
		So(err, ShouldBeNil)
		So(rp.State, ShouldEqual, model.StateNormal)
	}))
}

func TestAddReplyConfig(t *testing.T) {
	var (
		bs       []byte
		err      error
		id       int64
		typ      = int32(1)
		oid      = int64(1)
		adminID  = int64(1)
		category = int32(1)
		operator = string("管理员001")
		c        = context.Background()
		config   = &model.Config{}
	)
	s := New(conf.Conf)
	config.Oid = oid
	config.Type = typ
	config.Category = category
	config.AdminID = adminID
	config.Operator = operator
	configValue := map[string]int64{
		"showentry": 0,
		"showadmin": 1,
	}
	if bs, err = json.Marshal(configValue); err == nil {
		config.Config = string(bs)
	}
	if id, err = s.AddReplyConfig(c, config); err != nil {
		t.Errorf("d.AddConfig error(%v)", err)
	}
	t.Logf("d.AddReplyConfig result(%d)", id)
}
