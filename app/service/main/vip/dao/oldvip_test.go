package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoOldFrozenChange(t *testing.T) {
	convey.Convey("OldFrozenChange success", t, func() {
		defer gock.OffAll()
		httpMock("GET", _frozenChange).Reply(200).JSON(`{"code":0}`)
		err := d.OldFrozenChange(context.TODO(), 7593623)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("OldFrozenChange fail", t, func() {
		defer gock.OffAll()
		httpMock("GET", _frozenChange).Reply(200).JSON(`{"code":-400}`)
		err := d.OldFrozenChange(context.TODO(), 7593623)
		convey.So(err, convey.ShouldNotBeNil)
	})
}
