package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagSlice(t *testing.T) {
	convey.Convey("Test TagSlice", t, func() {
		_ = s.TagSlice(int8(0))
		convey.ShouldBeNil(nil)
	})
}

func TestTagMap(t *testing.T) {
	convey.Convey("Test TagMap", t, func() {
		_ = s.TagMap(int8(0), int32(0))
		convey.ShouldBeNil(nil)
	})
}

func TestTags(t *testing.T) {
	convey.Convey("Test Tags", t, func() {
		_, err := s.tags(context.Background())
		convey.ShouldBeNil(err)
	})
}

func TestTag3(t *testing.T) {
	convey.Convey("Test Tag3", t, func() {
		_ = s.Tag3(int64(0), int64(0))
		convey.ShouldBeNil(nil)
	})
}

func TestTags3(t *testing.T) {
	convey.Convey("Test Tags3", t, func() {
		_, err := s.tags3(context.Background())
		convey.ShouldBeNil(err)
	})
}
