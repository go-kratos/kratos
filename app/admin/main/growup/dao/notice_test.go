package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/growup/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NoticeSql(t *testing.T) {
	Convey("growup-admin", t, WithMysql(func(d *Dao) {
		var (
			a = &model.Notice{Title: "title", Type: 1, Platform: 1, Link: "www.bilibili.com", Status: 1}
		)
		_, err := d.InsertNotice(context.Background(), a)
		So(err, ShouldBeNil)
	}))

	Convey("growup-admin", t, WithMysql(func(d *Dao) {
		var (
			query = ""
		)
		_, err := d.NoticeCount(context.Background(), query)
		So(err, ShouldBeNil)
	}))

	Convey("growup-admin", t, WithMysql(func(d *Dao) {
		var (
			query       = ""
			from, limit = 0, 1000
		)
		res, err := d.Notices(context.Background(), query, from, limit)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))

	Convey("growup-admin", t, WithMysql(func(d *Dao) {
		var (
			kv       = "is_deleted = 1"
			id int64 = 1
		)
		_, err := d.UpdateNotice(context.Background(), kv, id)
		So(err, ShouldBeNil)
	}))
}
