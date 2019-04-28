package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/admin/main/tv/model"

	. "github.com/smartystreets/goconvey/convey"
)

func getMids(limit int) (res []int64, err error) {
	var (
		db  = d.DB.Where("deleted = 0")
		ups []*model.Upper
	)
	if err = db.Limit(limit).Find(&ups).Error; err != nil {
		fmt.Println("pickMid err ", err)
	}
	for _, v := range ups {
		res = append(res, v.MID)
	}
	return
}

func TestDao_UpList(t *testing.T) {
	Convey("TestDao_UpList", t, WithDao(func(d *Dao) {
		mids, errGet := getMids(5)
		if errGet != nil {
			fmt.Println("empty mids")
			return
		}
		res, pager, err := d.UpList(1, 1, mids)
		So(err, ShouldBeNil)
		So(pager, ShouldNotBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_VerifyIds(t *testing.T) {
	Convey("TestDao_VerifyIds", t, WithDao(func(d *Dao) {
		mids, errGet := getMids(25)
		if errGet != nil {
			fmt.Println("empty mids")
			return
		}
		okMids, err := d.VerifyIds(mids)
		So(err, ShouldBeNil)
		So(len(mids), ShouldBeGreaterThanOrEqualTo, len(okMids))
		data, _ := json.Marshal(okMids)
		fmt.Println(string(data))
	}))
}

func TestDao_AuditIds(t *testing.T) {
	Convey("TestDao_AuditIds", t, WithDao(func(d *Dao) {
		mids, errGet := getMids(15)
		if errGet != nil {
			fmt.Println("empty mids")
			return
		}
		err := d.AuditIds(mids, 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_DelCache(t *testing.T) {
	Convey("TestDao_DelCache", t, WithDao(func(d *Dao) {
		var (
			mid = int64(88895270)
			ctx = context.Background()
		)
		err := d.SetUpMetaCache(ctx, &model.UpMC{
			MID: mid,
		})
		So(err, ShouldBeNil)
		err = d.DelCache(ctx, mid)
		So(err, ShouldBeNil)
	}))
}
