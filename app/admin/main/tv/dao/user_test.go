package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetByMId(t *testing.T) {
	d.saveUser()
	convey.Convey("GetByMId", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(100000000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			userInfo, err := d.GetByMId(c, mid)
			ctx.Convey("Then err should be nil.userInfo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(userInfo, convey.ShouldNotBeNil)
			})
		})
	})
	d.delUser()
}

func (d *Dao) saveUser() {
	userInfo := &model.TvUserInfo{
		MID:           100000000,
		RecentPayTime: 1544613229,
	}
	if err := d.DB.Save(userInfo).Error; err != nil {
		log.Error("SaveUser error(%v)", err)
	}
}

func (d *Dao) delUser() {
	userInfo := &model.TvUserInfo{
		MID: 100000000,
	}
	if err := d.DB.Where("mid = ?", userInfo.MID).Delete(userInfo).Error; err != nil {
		log.Error("DelUser error(%v)", err)
	}
}
