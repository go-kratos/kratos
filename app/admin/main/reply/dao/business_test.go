package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/reply/model"

	. "github.com/smartystreets/goconvey/convey"
)

func delBusiness(d *Dao, id int64) error {
	_delBsSQL := "Delete from business where id = %d"
	_, err := d.db.Exec(context.Background(), fmt.Sprintf(_delBsSQL, id))
	if err != nil {
		return err
	}
	return nil
}

func TestBusiness(t *testing.T) {
	var (
		id1 int64
		id2 int64
		err error
		bs1 = model.Business{
			Type:   244,
			Name:   "TestName",
			Appkey: "TestAppKey",
			Remark: "TestRemark",
			Alias:  "TestAlias",
		}
		bs2 = model.Business{
			Type:   245,
			Name:   "TestName",
			Appkey: "TestAppKey",
			Remark: "TestRemark",
			Alias:  "TestAlias2",
		}
		b  []*model.Business
		bu *model.Business
	)
	c := context.Background()
	Convey("test dao business", t, WithDao(func(d *Dao) {
		id1, err = d.InBusiness(c, bs1.Type, bs1.Name, bs1.Appkey, bs1.Remark, bs1.Alias)
		So(err, ShouldBeNil)
		So(id1, ShouldBeGreaterThan, 0)
		defer delBusiness(d, id1)

		id2, err = d.InBusiness(c, bs2.Type, bs2.Name, bs2.Appkey, bs2.Remark, bs2.Alias)
		So(err, ShouldBeNil)
		So(id2, ShouldBeGreaterThan, 0)
		defer delBusiness(d, id2)
		Convey("list business", WithDao(func(d *Dao) {
			b, err = d.ListBusiness(c, 0)
			So(err, ShouldBeNil)
			So(len(b), ShouldBeGreaterThan, 0)
		}))
		Convey("update business", WithDao(func(d *Dao) {
			bu, err = d.Business(c, 245)
			So(err, ShouldBeNil)
			So(bu.Name, ShouldEqual, "TestName")

			_, err = d.UpBusiness(c, "TestChangeName", bu.Appkey, bu.Remark, bu.Alias, bu.Type)
			So(err, ShouldBeNil)

			bu, err = d.Business(c, 245)
			So(err, ShouldBeNil)
			So(bu.Name, ShouldEqual, "TestChangeName")
		}))
		Convey("update business state", WithDao(func(d *Dao) {
			_, err = d.UpBusinessState(c, 1, bu.Type)
			So(err, ShouldBeNil)

			bu, err = d.Business(c, 245)
			So(err, ShouldBeNil)
			So(bu.Type, ShouldEqual, 245)
		}))
	}))
}
