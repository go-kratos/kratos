package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRelation(t *testing.T) {
	convey.Convey("Relation", t, func() {
		res, err := s.Relation(context.TODO(), 1, 2)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRelations(t *testing.T) {
	convey.Convey("Relations", t, func() {
		res, err := s.Relations(context.TODO(), 1, []int64{2, 3})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestRichRelations2(t *testing.T) {
	convey.Convey("RichRelations2", t, func() {
		res, err := s.RichRelations2(context.TODO(), 1, []int64{2, 3})
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestBlacks(t *testing.T) {
	convey.Convey("Blacks", t, func() {
		res, err := s.Blacks(context.TODO(), 1)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestAttentions(t *testing.T) {
	convey.Convey("Attentions", t, func() {
		res, err := s.Attentions(context.TODO(), 1)
		convey.So(res, convey.ShouldNotBeNil)
		convey.So(err, convey.ShouldBeNil)
	})
}
