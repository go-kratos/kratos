package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/admin/main/usersuit/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestMcJury .

func Test_pingMC(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.pingMC(context.Background())
		So(err, ShouldBeNil)
	})
}

func Test_SetMedalOwnersache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		nos := make([]*model.MedalOwner, 0)
		no := &model.MedalOwner{}
		no.ID = 1
		no.MID = 1
		no.NID = 2
		no.CTime = xtime.Time(time.Now().Second())
		no.MTime = xtime.Time(time.Now().Second())
		nos = append(nos, no)
		err := d.SetMedalOwnersache(context.Background(), 1, nos)
		So(err, ShouldBeNil)
	})
	Convey("should return err be nil", t, func() {
		res, err := d.MedalOwnersCache(context.Background(), 1)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func Test_DelMedalOwnersCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelMedalOwnersCache(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_MedalActivatedCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		nid, err := d.MedalActivatedCache(context.Background(), 1)
		fmt.Printf("%+v\n", nid)
		So(nid, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func Test_SetMedalActivatedCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.SetMedalActivatedCache(context.Background(), 1, 22)
		So(err, ShouldBeNil)
	})
}

func Test_DelMedalActivatedCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelMedalActivatedCache(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_PopupCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		nid, err := d.PopupCache(context.Background(), 1)
		fmt.Printf("%+v\n", nid)
		So(nid, ShouldNotBeNil)
		So(nid, ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)
	})
}

func Test_SetPopupCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.SetPopupCache(context.Background(), 1, 3)
		So(err, ShouldBeNil)
	})
}
func Test_DelPopupCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelPopupCache(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}

func Test_RedPointCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		nid, err := d.RedPointCache(context.Background(), 1)
		fmt.Printf("%+v\n", nid)
		So(nid, ShouldNotBeNil)
		So(nid, ShouldBeGreaterThanOrEqualTo, 0)
		So(err, ShouldBeNil)
	})
}

func Test_SetRedPointCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.SetRedPointCache(context.Background(), 1, 3)
		So(err, ShouldBeNil)
	})
}
func Test_DelRedPointCache(t *testing.T) {
	Convey("should return err be nil", t, func() {
		err := d.DelRedPointCache(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}
