package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoConcatCondition(t *testing.T) {
	convey.Convey("ConcatCondition", t, func(ctx convey.C) {
		var (
			conditions Condition = Condition{Key: "id", Operator: "=", Value: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			conditionStr, args, hasOperator := ConcatCondition(conditions)
			ctx.Convey("Then conditionStr,args,hasOperator should not be nil.", func(ctx convey.C) {
				ctx.So(hasOperator, convey.ShouldNotBeNil)
				ctx.So(args, convey.ShouldNotBeNil)
				ctx.So(conditionStr, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAndCondition(t *testing.T) {
	convey.Convey("AndCondition", t, func(ctx convey.C) {
		var (
			conditions Condition
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := AndCondition(conditions)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoOrCondition(t *testing.T) {
	convey.Convey("OrCondition", t, func(ctx convey.C) {
		var (
			conditions Condition
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := OrCondition(conditions)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoaddLogicOperator(t *testing.T) {
	convey.Convey("addLogicOperator", t, func(ctx convey.C) {
		var (
			operator   = int(0)
			conditions Condition
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result := addLogicOperator(operator, conditions)
			ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
