package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_List(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.dao.List(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("list not exist", t, func() {
		res, err := s.dao.List(context.TODO(), 999)
		So(err, ShouldBeNil)
		So(res, ShouldBeNil)
	})
}

func Test_rawListArticles(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.rawListArticles(context.TODO(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_ListInfo(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.ListInfo(context.TODO(), 821)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("null data", t, func() {
		res, err := s.ListInfo(context.TODO(), 999999999)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	})
}

func Test_Lists(t *testing.T) {
	Convey("get data", t, func() {
		res, err := s.Lists(context.TODO(), []int64{3})
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		So(res[3].ImageURL, ShouldNotBeEmpty)
	})
	Convey("null data", t, func() {
		res, err := s.Lists(context.TODO(), []int64{})
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	})
}
