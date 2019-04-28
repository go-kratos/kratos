package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/job/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDMInfos(t *testing.T) {
	Convey("", t, func() {
		dms, err := testDao.DMInfos(context.TODO(), 1, 1221)
		if err != nil {
			fmt.Println(err)
		}
		So(err, ShouldBeNil)
		So(dms, ShouldNotBeEmpty)
	})
}

func TestDMHides(t *testing.T) {
	Convey("", t, func() {
		dms, err := testDao.DMHides(context.TODO(), 1, 1221, 1)
		if err != nil {
			fmt.Println(err)
		}
		So(err, ShouldBeNil)
		So(dms, ShouldNotBeEmpty)
	})
}

func TestUpdateDM(t *testing.T) {
	Convey("", t, func() {
		dm := &model.DM{
			ID:    1234555,
			Oid:   1221,
			State: 3,
		}
		_, err := testDao.UpdateDM(context.TODO(), dm)
		So(err, ShouldBeNil)
	})
}

func TestUpdateDMStates(t *testing.T) {
	Convey("", t, func() {
		dmids := []int64{709150141, 709150142}
		_, err := testDao.UpdateDMStates(context.TODO(), 1221, dmids, 1)
		So(err, ShouldBeNil)
	})
}
