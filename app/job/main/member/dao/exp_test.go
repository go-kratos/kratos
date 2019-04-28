package dao

import (
	"context"
	"testing"
	"time"

	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SelExp(t *testing.T) {
	Convey("selExp", t, func() {
		exp, err := d.SelExp(context.TODO(), 2233)
		So(err, ShouldBeNil)
		So(exp, ShouldNotBeEmpty)
	})
}

func Test_UpdateExpAped(t *testing.T) {
	Convey("updateExpAped", t, func() {
		if _, err := d.UpdateExpAped(context.TODO(), 2233, 666, 1); err != nil {
			So(err, ShouldBeNil)
		}
	})
}

func Test_InitExp(t *testing.T) {
	Convey("InitExp", t, func() {
		if err := d.InitExp(context.Background(), 1111); err != nil {
			So(err, ShouldNotBeNil)
		}
	})
}

func Test_UpdateExpFlag(t *testing.T) {
	Convey("UpdateExpFlag", t, func() {
		if err := d.UpdateExpFlag(context.Background(), 1111, 123, xtime.Time(time.Now().Unix())); err != nil {
			So(err, ShouldNotBeNil)
		}
	})
}
