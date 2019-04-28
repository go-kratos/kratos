package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDao_UpdateAccFace(t *testing.T) {
	Convey("UpdateAccFace", t, func() {
		Convey("When everything goes positive", func() {
			defer gock.OffAll()
			httpMock("POST", _updateFaceURL).Reply(200).JSON(`{"code":0}`)
			err := d.UpdateAccFace(context.TODO(), 1, "1456")
			So(err, ShouldBeNil)
		})
		Convey("When everything goes negative", func() {
			defer gock.OffAll()
			httpMock("POST", _updateFaceURL).Reply(200).JSON(`{"code":500}`)
			err := d.UpdateAccFace(context.TODO(), 1, "1456")
			Convey("Then err should be nil.", func() {
				So(err, ShouldNotBeNil)
			})
		})

	})
}
