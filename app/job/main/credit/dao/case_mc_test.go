package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoDelCaseInfoCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelCaseInfoCache(context.TODO(), 12)
		So(err, ShouldBeNil)
	})
}

func TestDaoDelVoteCaseCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelVoteCaseCache(context.TODO(), 12, 11)
		So(err, ShouldBeNil)
	})
}
