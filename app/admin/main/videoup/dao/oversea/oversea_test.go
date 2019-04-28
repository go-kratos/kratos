package manager

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_UpArea(t *testing.T) {
	Convey("UpArea", t, WithDao(func(d *Dao) {
		_, err := d.UpPolicyRelation(context.TODO(), 1212121, 1)
		So(err, ShouldBeNil)
	}))
}
func Test_PolicyGroup(t *testing.T) {
	Convey("PolicyGroup", t, WithDao(func(d *Dao) {
		_, _, err := d.PolicyGroups(context.TODO(), 0, 0, 0, 0, 0, 0, "", "")
		So(err, ShouldBeNil)
	}))
}
func Test_PolicyRelation(t *testing.T) {
	Convey("PolicyRelation", t, WithDao(func(d *Dao) {
		_, err := d.PolicyRelation(context.TODO(), 1212121)
		So(err, ShouldBeNil)
	}))
}
func Test_PolicyGroupsByIds(t *testing.T) {
	Convey("PolicyGroupsByIds", t, WithDao(func(d *Dao) {
		_, err := d.PolicyGroupsByIds(context.TODO(), []int64{1212121})
		So(err, ShouldBeNil)
	}))
}

func Test_ArchiveGroups(t *testing.T) {
	Convey("ArchiveGroups", t, WithDao(func(d *Dao) {
		_, err := d.ArchiveGroups(context.TODO(), 1212121)
		So(err, ShouldBeNil)
	}))
}
func Test_PolicyItems(t *testing.T) {
	Convey("ItemsByGroup", t, WithDao(func(d *Dao) {
		_, err := d.PolicyItems(context.Background(), 1)
		So(err, ShouldBeNil)
	}))
}
func Test_ZoneIDs(t *testing.T) {
	Convey("ZoneIDs", t, WithDao(func(d *Dao) {
		_, err := d.ZoneIDs(context.TODO(), []int64{1212121})
		So(err, ShouldBeNil)
	}))
}
