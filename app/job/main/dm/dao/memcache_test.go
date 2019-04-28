package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubjectCache(t *testing.T) {
	Convey("subject cache", t, func() {
		sub, err := testDao.SubjectCache(context.TODO(), 1, 1221)
		So(err, ShouldBeNil)
		So(sub, ShouldNotBeNil)
	})
}

func TestSetSubjectCache(t *testing.T) {
	sub := &model.Subject{
		Type: 1,
		Oid:  1221,
	}
	Convey("add subject cache", t, func() {
		err := testDao.SetSubjectCache(context.TODO(), sub)
		So(err, ShouldBeNil)
	})
}

func TestDurationCache(t *testing.T) {
	Convey("", t, func() {
		_, err := testDao.DurationCache(context.TODO(), 1221)
		So(err, ShouldBeNil)
	})
}

func TestSetDurationCachee(t *testing.T) {
	Convey("", t, func() {
		err := testDao.SetDurationCache(context.TODO(), 1221, 10000)
		So(err, ShouldBeNil)
	})
}
