package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/dm/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReduceMoral(t *testing.T) {
	arg := &model.ReduceMoral{
		UID:        150781,
		Moral:      -1,
		Origin:     2,
		Reason:     1,
		ReasonType: 1,
		Operator:   "zhang",
		IsNotify:   1,
		Remark:     "dm admin test",
	}
	Convey("", t, func() {
		err := testDao.ReduceMoral(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestBlockUser1(t *testing.T) {
	arg := &model.BlockUser{
		UID:             150781,
		BlockForever:    1,
		BlockTimeLength: 0,
		BlockRemark:     model.BlockReason[5],
		ReasonType:      5,
		Operator:        "zhang",
		OriginType:      2,
		Moral:           10,
		OriginURL:       "aaaaa",
		OriginContent:   "test delete1",
		OriginTitle:     "test title",
		IsNotify:        0,
	}
	Convey("", t, func() {
		err := testDao.BlockUser(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestBlockUser2(t *testing.T) {
	arg := &model.BlockUser{
		UID:             150781,
		BlockForever:    0,
		BlockTimeLength: 5,
		BlockRemark:     model.BlockReason[5],
		ReasonType:      5,
		Operator:        "zhang",
		OriginType:      2,
		Moral:           10,
		OriginURL:       "aaaaa",
		OriginContent:   "test delete2",
		OriginTitle:     "test title",
		IsNotify:        0,
	}
	Convey("", t, func() {
		err := testDao.BlockUser(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}
