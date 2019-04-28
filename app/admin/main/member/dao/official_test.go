package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOfficial(t *testing.T) {
	convey.Convey("Official", t, func() {
		o, err := d.Official(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
	})
}

func TestOfficialEdit(t *testing.T) {
	convey.Convey("OfficialEdit", t, func() {
		o, err := d.OfficialEdit(context.Background(), 123, 1, "title", "desc")
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
	})
}

func TestOfficials(t *testing.T) {
	convey.Convey("Officials", t, func() {
		o, total, err := d.Officials(context.Background(), 123, []int8{1}, time.Time{}, time.Now(), 1, 20)
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(total, convey.ShouldBeGreaterThan, 0)
	})
}

func TestOfficialDoc(t *testing.T) {
	convey.Convey("OfficialDoc", t, func() {
		o, err := d.OfficialDoc(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
	})
}

func TestOfficialDocs(t *testing.T) {
	convey.Convey("OfficialDocs", t, func() {
		o, total, err := d.OfficialDocs(context.Background(), 123, []int8{1}, []int8{1}, "", time.Time{}, time.Now(), 1, 20)
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(total, convey.ShouldBeGreaterThan, 0)
	})
}

func TestOfficialDocAudit(t *testing.T) {
	convey.Convey("OfficialDocAudit", t, func() {
		err := d.OfficialDocAudit(context.Background(), 123, 1, "guan", true, "")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestOfficialDocEdit(t *testing.T) {
	convey.Convey("OfficialDocEdit", t, func() {
		err := d.OfficialDocEdit(context.Background(), 123, "guan", 1, 1, "title", "desc", "extra", "12121", true)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestOfficialDocsByMids(t *testing.T) {
	convey.Convey("OfficialDocsByMids", t, func() {
		res, err := d.OfficialDocsByMids(context.Background(), []int64{2, 3, 10, 12})
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(res), convey.ShouldBeGreaterThan, 0)
	})
}

func TestOfficialDocSubmit(t *testing.T) {
	convey.Convey("OfficialDocSubmit", t, func() {
		err := d.OfficialDocSubmit(context.Background(), 123, "guan", 1, 1, "title", "desc", "extra", "12121", true, "admin")
		convey.So(err, convey.ShouldBeNil)
	})
}
