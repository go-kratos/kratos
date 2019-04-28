package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpdateJuryExpired(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.UpdateJuryExpired(context.TODO(), 88889017, time.Now())
		So(err, ShouldBeNil)
	})
}

func Test_LoadConf(t *testing.T) {
	Convey("should return err be nil & map be nil", t, func() {
		m, err := d.LoadConf(context.TODO())
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	})
}
