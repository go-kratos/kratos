package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/growup/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_BlockSql(t *testing.T) {
	Convey("query apply_at from up_info_video by mid", t, WithMysql(func(d *Dao) {
		mid := int64(1011)
		_, err := d.ApplyAt(context.Background(), mid)
		So(err, ShouldBeNil)
	}))

	Convey("insert block up to blacklist", t, WithMysql(func(d *Dao) {
		v := &model.Blocked{MID: int64(1011), Nickname: "hello", OriginalArchiveCount: 10, MainCategory: 1, Fans: 100, ApplyAt: xtime.Time(time.Now().Unix())}
		_, err := d.InsertBlocked(context.Background(), v)
		So(err, ShouldBeNil)
	}))

	Convey("update blocked is_deleted", t, WithMysql(func(d *Dao) {
		var (
			mid int64 = 1011
			del       = 1
		)
		_, err := d.UpdateBlockedState(context.Background(), mid, del)
		So(err, ShouldBeNil)
	}))

	Convey("del blocked", t, WithMysql(func(d *Dao) {
		var (
			mid int64 = 1011
		)
		_, err := d.DelFromBlocked(context.Background(), mid)
		So(err, ShouldBeNil)
	}))

	Convey("get blocked count", t, WithMysql(func(d *Dao) {
		var (
			query = "1 = 1"
		)
		_, err := d.BlockCount(context.Background(), query)
		So(err, ShouldBeNil)
	}))

	Convey("query blocked user in black list", t, WithMysql(func(d *Dao) {
		var (
			query = "1 = 1"
		)
		_, err := d.QueryFromBlocked(context.Background(), query)
		So(err, ShouldBeNil)
	}))

	Convey("check user is blocked", t, WithMysql(func(d *Dao) {
		var (
			mid int64 = 1011
		)
		_, err := d.Blocked(context.Background(), mid)
		So(err, ShouldBeNil)
	}))
}
