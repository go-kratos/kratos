package pendant

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoUpEquipMID(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.UpEquipMID(context.Background(), 11)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoUpEquipExpires(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.UpEquipExpires(context.TODO(), 11, 123121)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoPendantEquipMID(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.PendantEquipMID(context.TODO(), 11)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoExpireEquipPendant(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.ExpireEquipPendant(context.TODO(), 11)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}

func TestDaoPendantEquipGidPid(t *testing.T) {
	Convey("should return err be nil", t, func() {
		res, err := d.PendantEquipGidPid(context.TODO(), 11)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	})
}
