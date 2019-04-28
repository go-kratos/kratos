package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_AssetDBConnect(t *testing.T) {
	var (
		c        = context.Background()
		host     = "172.16.33.205"
		port     = "3308"
		user     = "test"
		password = "test"
	)

	Convey("AssetDBConnect", t, WithService(func(s *Service) {
		res, err := svr.AssetDBConnect(c, host, port, user, password)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_AssetDBAdd(t *testing.T) {
	var (
		c           = context.Background()
		name        = "bilibili_reply"
		description = "评论数据库"
		host        = "172.16.33.205"
		port        = "3308"
		user        = "test"
		password    = "test"
	)

	Convey("AssetDBAdd", t, WithService(func(s *Service) {
		res, err := svr.AssetDBAdd(c, name, description, host, port, user, password)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_AssetTableFields(t *testing.T) {
	var (
		c     = context.Background()
		db    = "bilibili_reply"
		regex = "reply_([0-9]{1,3})"
	)

	Convey("AssetTableFields", t, WithService(func(s *Service) {
		res, _, err := svr.AssetTableFields(c, db, regex)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_BusinessUpdate(t *testing.T) {
	var (
		c     = context.Background()
		name  = "dm"
		field = "description"
		value = "弹幕11"
	)

	Convey("BusinessUpdate", t, WithService(func(s *Service) {
		res, err := svr.BusinessUpdate(c, name, field, value)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
