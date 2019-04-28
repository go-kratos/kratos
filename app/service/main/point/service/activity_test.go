package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/service/main/point/dao"
	"go-common/app/service/main/point/model"
	xsql "go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

//以service层activityGiveTimes方法为例
func TestServiceactivityGiveTimes(t *testing.T) {
	convey.Convey("activityGiveTimes", t, func(ctx convey.C) {
		//被测方法与桩方法的变量及参数初始化
		var (
			c          = context.Background()
			mid        = int64(4780461)
			changeType = int(3)
			point      = int64(0)
			phs        []*model.PointHistory
			ph         = &model.PointHistory{
				ID:           13,
				Mid:          4780461,
				Point:        60,
				ChangeType:   1,
				PointBalance: 418,
			}
		)
		phs = append(phs, ph)
		//convey包裹调用service测试方法及断言部分
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//使用monkey包构造此service方法下所有依赖的dao层方法
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "SelPointHistory", func(_ *dao.Dao, _ context.Context, _ int64, _, _ xtime.Time) ([]*model.PointHistory, error) {
				return phs, nil
			})
			sendTime, err := s.activityGiveTimes(c, mid, changeType, point)
			ctx.Convey("Then err should be nil.sendTime should not be nil.", func(ctx convey.C) {
				t.Logf("sendTime:%+v", sendTime)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sendTime, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When dao return err", func(ctx convey.C) {
			//使用monkey包构造此service方法下调dao层失败的情况
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "SelPointHistory", func(_ *dao.Dao, _ context.Context, _ int64, _, _ xtime.Time) ([]*model.PointHistory, error) {
				return nil, fmt.Errorf("get history err")
			})
			_, err := s.activityGiveTimes(c, mid, changeType, point)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		//convey teardown部分(此处UnpatchAll解除所有Patch打桩绑定，确保后续测试流程不被打桩影响)
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}

func TestServiceactiveSendPoint(t *testing.T) {
	convey.Convey("activeSendPoint", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tx  *xsql.Tx
			phs []*model.PointHistory
			ph  = &model.PointHistory{
				ID:           13,
				Mid:          4780461,
				Point:        60,
				ChangeType:   1,
				PointBalance: 418,
			}
		)
		phs = append(phs, ph)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			//activeSendPoint方法中因调用了内部私有activityGiveTimes和updatePoint方法，无法运用monkey反射机制报panic，
			// 可以采用给私有方法的下一级打桩或将私有方法转为公有
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "SelPointHistory", func(_ *dao.Dao, _ context.Context, _ int64, _, _ xtime.Time) ([]*model.PointHistory, error) {
				return phs, nil
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(s.dao), "InsertPointHistory", func(_ *dao.Dao, _ context.Context, _ *xsql.Tx, ph *model.PointHistory) (int64, error) {
				return 0, nil
			})

			activePoint, err := s.activeSendPoint(c, tx, ph)
			ctx.Convey("Then err should be nil.activePoint should not be nil.", func(ctx convey.C) {
				t.Logf("activepoint:%+v", activePoint)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(activePoint, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			monkey.UnpatchAll()
		})
	})
}
