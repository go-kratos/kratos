package dao

import (
	"testing"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ApplyCount(t *testing.T) {
	mid := int64(1)
	pending := 0
	Convey("add apply", t, func() {
		err := d.AddApply(ctx(), mid, "content", "category")
		So(err, ShouldBeNil)
		Convey("get apply", func() {
			author, err := d.RawAuthor(ctx(), mid)
			So(err, ShouldBeNil)
			So(author, ShouldResemble, &model.AuthorLimit{State: pending})
		})
		Convey("apply count", func() {
			res, err := d.ApplyCount(ctx())
			So(err, ShouldBeNil)
			So(res, ShouldBeGreaterThan, 0)
		})
		Convey("add twice should be ok", func() {
			err := d.AddApply(ctx(), mid, "content", "category")
			So(err, ShouldBeNil)
			author, err := d.RawAuthor(ctx(), mid)
			So(err, ShouldBeNil)
			So(author, ShouldResemble, &model.AuthorLimit{State: pending})
		})
	})
}

func Test_Authors(t *testing.T) {
	Convey("should get authors", t, func() {
		res, err := d.Authors(ctx())
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func Test_Identity(t *testing.T) {
	Convey("should return identity when no error", t, func() {
		httpMock("get", _verifyAPI).Reply(200).JSON(`{"ts":1514341945,"code":0,"data":{"identify":1,"phone":0}}`)
		res, err := d.Identify(ctx(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &model.Identify{Identify: 1, Phone: 0})
	})
	Convey("should return error when code != 0", t, func() {
		httpMock("get", _verifyAPI).Reply(200).JSON(`{"ts":1514341945,"code":-3,"data":null}`)
		_, err := d.Identify(ctx(), 1)
		So(err, ShouldEqual, ecode.SignCheckErr)
	})
}
