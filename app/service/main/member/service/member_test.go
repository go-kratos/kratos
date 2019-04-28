package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/service/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestSetOfficialDoc(t *testing.T) {
	var (
		c              = context.TODO()
		argOfficialDoc = &model.ArgOfficialDoc{
			Mid:        1234567,
			Role:       1,
			Title:      "test123515",
			Name:       "test12345",
			CreditCode: "12dfrtg12",
		}
	)
	convey.Convey("SetOfficialDoc", t, func(ctx convey.C) {
		err := s.SetOfficialDoc(c, argOfficialDoc)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestOfficialDoc(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("SetOfficialDoc", t, func(ctx convey.C) {
		officialDoc, err := s.OfficialDoc(c, 777777)
		fmt.Printf("%v", officialDoc.SubmitTime)
		ctx.Convey("Error should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		ctx.Convey("Result should not be nil", func(ctx convey.C) {
			ctx.So(officialDoc, convey.ShouldNotBeNil)
		})
	})
}
