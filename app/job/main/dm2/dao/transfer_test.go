package dao

import (
	"context"
	"testing"

	"go-common/app/job/main/dm2/conf"
	"go-common/app/job/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransfers(t *testing.T) {
	var (
		d = New(conf.Conf)
		c = context.TODO()
	)
	Convey("test transferJob", t, func() {
		_, err := d.Transfers(c, model.StatInit)
		So(err, ShouldBeNil)
	})
}

func TestDmIndexs(t *testing.T) {
	var (
		d = New(conf.Conf)
		c = context.TODO()
	)
	Convey("test DmIndexs", t, func() {
		_, _, _, err := d.DMIndexs(c, 1, 1012, 0, 10)
		So(err, ShouldBeNil)
	})
}

func TestUpdateTransfer(t *testing.T) {
	var (
		d     = New(conf.Conf)
		c     = context.TODO()
		trans = &model.Transfer{
			ID:      265,
			FromCid: 233,
			ToCid:   1221,
			Dmid:    333,
			Mid:     1,
			Offset:  0.00,
			State:   0,
		}
	)
	Convey("test update job", t, func() {
		d.UpdateTransfer(c, trans)
	})
}
