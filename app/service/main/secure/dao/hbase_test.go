package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/secure/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAddLocs(t *testing.T) {
	Convey("TestAddLocs", t, func() {
		err := d.AddLocs(context.TODO(), 2, 3, 11)
		d.AddLocs(context.TODO(), 2, 4, 12)
		d.AddLocs(context.TODO(), 2, 4, 13)
		d.AddLocs(context.TODO(), 2, 4, 14)
		d.AddLocs(context.TODO(), 2, 3, 15)
		So(err, ShouldBeNil)

	})
}

func TestLocs(t *testing.T) {
	Convey("TestLocs", t, func() {
		locs, err := d.Locs(context.TODO(), 2)
		So(err, ShouldBeNil)
		So(locs, ShouldNotBeNil)
	})
}

func TestAddEcpt(t *testing.T) {
	Convey("TestAddEcpt", t, func() {
		err := d.AddException(context.TODO(), &model.Log{Mid: 2, Time: xtime.Time(1111), IP: 222, LocationID: 3, Location: "aa"})
		So(err, ShouldBeNil)

		err = d.AddException(context.TODO(), &model.Log{Mid: 2, Time: xtime.Time(1311), IP: 222, LocationID: 3, Location: "aa"})
		So(err, ShouldBeNil)

		err = d.AddException(context.TODO(), &model.Log{Mid: 2333, Time: xtime.Time(1111), IP: 222, LocationID: 3, Location: "aa"})
		So(err, ShouldBeNil)
	})
}

func TestAddFeedBack(t *testing.T) {
	Convey("TestAddFeedBack", t, func() {
		err := d.AddFeedBack(context.TODO(), &model.Log{Mid: 2, Time: xtime.Time(1111), IP: 222, Type: 2, LocationID: 3, Location: "aa"})
		if err != nil {
			t.Errorf("test hbase add  err %v", err)
		}
	})
}

func TestExcep(t *testing.T) {
	Convey("TestExcep", t, func() {
		locs, err := d.ExceptionLoc(context.TODO(), 2)
		if err != nil {
			t.Errorf("test hbase add  err %v", err)
		} else {
			for _, l := range locs {
				t.Logf("locs %v", l)
			}
		}
	})
}
