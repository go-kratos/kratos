package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDaoBlockAccount(t *testing.T) {
	convey.Convey("BlockAccount", t, func() {
		defer gock.OffAll()
		httpMock("GET", d.c.Property.BlockAccountURL).Reply(200).JSON(`{"code":0}`)
		err := d.BlockAccount(context.TODO(), 7593623, "1")
		convey.So(err, convey.ShouldBeNil)
	})
}
