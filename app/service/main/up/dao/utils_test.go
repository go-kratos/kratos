package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoConcatCondition(t *testing.T) {
	var (
		conditions Condition
	)
	convey.Convey("ConcatCondition", t, func(ctx convey.C) {
		conditionStr, args, hasOperator := ConcatCondition(conditions)
		ctx.Convey("Then conditionStr,args,hasOperator should not be nil.", func(ctx convey.C) {
			ctx.So(hasOperator, convey.ShouldNotBeNil)
			ctx.So(args, convey.ShouldBeNil)
			ctx.So(conditionStr, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAndCondition(t *testing.T) {
	var (
		conditions Condition
	)
	convey.Convey("AndCondition", t, func(ctx convey.C) {
		result := AndCondition(conditions)
		ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoOrCondition(t *testing.T) {
	var (
		conditions Condition
	)
	convey.Convey("OrCondition", t, func(ctx convey.C) {
		result := OrCondition(conditions)
		ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoaddLogicOperator(t *testing.T) {
	var (
		operator   = int(0)
		conditions Condition
	)
	convey.Convey("addLogicOperator", t, func(ctx convey.C) {
		result := addLogicOperator(operator, conditions)
		ctx.Convey("Then result should not be nil.", func(ctx convey.C) {
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}
func TestSplit(t *testing.T) {
	var (
		start = 0
		end   = 100
		size  = 10
		f     = func(start int, end int) {}
	)
	convey.Convey("split", t, func(ctx convey.C) {
		Split(start, end, size, f)
	})
}
