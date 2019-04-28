package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bouk/monkey"

	"go-common/app/service/main/member/dao"
	"go-common/app/service/main/member/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestNonage(t *testing.T) {
	convey.Convey("", t, func() {
		var (
			code18   = "340702199110120012"
			code15   = "130503670401001"
			birthday time.Time
			f        bool
			err      error
		)
		birthday, _, err = ParseIdentity(code18)
		convey.So(err, convey.ShouldBeNil)
		convey.So(birthday, convey.ShouldResemble, time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local))
		f, err = isAdult(birthday, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(f, convey.ShouldBeTrue)

		birthday, _, err = ParseIdentity(code15)
		convey.So(err, convey.ShouldBeNil)
		convey.So(birthday, convey.ShouldResemble, time.Date(1967, time.Month(04), 01, 0, 0, 0, 0, time.Local))
		f, err = isAdult(birthday, time.Now())
		convey.So(err, convey.ShouldBeNil)
		convey.So(f, convey.ShouldBeTrue)

		_, err = isAdult(time.Now(), birthday)
		convey.So(err, convey.ShouldNotBeNil)

		type bd struct {
			day      time.Time
			anchor   time.Time
			expected bool
		}

		var (
			birthdays = []bd{
				{time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local), time.Date(2008, time.Month(10), 12, 0, 0, 0, 0, time.Local), false},
				{time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local), time.Date(2009, time.Month(9), 30, 0, 0, 0, 0, time.Local), false},
				{time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local), time.Date(2009, time.Month(10), 11, 0, 0, 0, 0, time.Local), false},
				{time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local), time.Date(2009, time.Month(10), 12, 0, 0, 0, 0, time.Local), true},
				{time.Date(1991, time.Month(10), 12, 0, 0, 0, 0, time.Local), time.Date(2009, time.Month(11), 1, 0, 0, 0, 0, time.Local), true},

				{time.Date(2000, time.Month(01), 01, 0, 0, 0, 0, time.Local), time.Date(2017, time.Month(12), 31, 0, 0, 0, 0, time.Local), false},
				{time.Date(2000, time.Month(01), 01, 0, 0, 0, 0, time.Local), time.Date(2018, time.Month(01), 01, 0, 0, 0, 0, time.Local), true},
				{time.Date(2000, time.Month(01), 01, 0, 0, 0, 0, time.Local), time.Date(2018, time.Month(02), 28, 0, 0, 0, 0, time.Local), true},
				{time.Date(2000, time.Month(01), 01, 0, 0, 0, 0, time.Local), time.Date(2018, time.Month(12), 31, 0, 0, 0, 0, time.Local), true},
				{time.Date(2000, time.Month(01), 01, 0, 0, 0, 0, time.Local), time.Date(2019, time.Month(01), 01, 0, 0, 0, 0, time.Local), true},
			}
		)
		for _, b := range birthdays {
			f, err = isAdult(b.day, b.anchor)
			convey.So(err, convey.ShouldBeNil)
			convey.So(f, convey.ShouldEqual, b.expected)
		}

		var (
			encrytedCode18 = []byte("y8Gxa0TUFCb0Tbtbw99Fm9si0zubpwVC6c3End/B4OnVYGesi7Y0rHOvD+wQusQ2+cmBjiOMZ0bvpfmH26b58RRGt9a8DDJno1zMpstyJWfdC75g0cgppgVWwkmGr9EfiUFSyKGHJFPFpT/YIyk+0kAA4jcD+kNQ23/DV2tpsrI=")
			codeBytes      []byte
		)
		codeBytes, err = s.realnameCryptor.CardDecrypt(encrytedCode18)
		convey.So(err, convey.ShouldBeNil)
		fmt.Println(string(codeBytes))
	})
}

func TestServicerealnameAlipayApply(t *testing.T) {
	convey.Convey("realnameAlipayApply", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(s.mbDao), "RealnameAlipayApply", func(_ *dao.Dao, _ context.Context, _ int64) (*model.RealnameAlipayApply, error) {
				return nil, nil
			})
			defer guard.Unpatch()
			info, err := s.realnameAlipayApply(c, mid)
			ctx.Convey("Then err should be nil,info should not be nil", func(ctx convey.C) {
				ctx.So(info, convey.ShouldNotBeNil)
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
