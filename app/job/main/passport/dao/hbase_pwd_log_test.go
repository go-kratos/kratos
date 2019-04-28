package dao

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"go-common/app/job/main/passport/model"
	"go-common/library/database/hbase.v2"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tsuna/gohbase/hrpc"
)

func TestDao_AddPwdLogHBase(t *testing.T) {
	convey.Convey("AddPwdLogHBase", t, func(ctx convey.C) {
		var (
			c = context.Background()
			a = &model.PwdLog{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.pwdLogHBase), "PutStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ map[string]map[string][]byte, _ ...func(hrpc.Call) error) (res *hrpc.Result, err error) {
				return nil, nil
			})
			defer mock.Unpatch()
			err := d.AddPwdLogHBase(c, a)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When error", func(ctx convey.C) {
			mock := monkey.PatchInstanceMethod(reflect.TypeOf(d.pwdLogHBase), "PutStr", func(_ *hbase.Client, _ context.Context, _, _ string, _ map[string]map[string][]byte, _ ...func(hrpc.Call) error) (res *hrpc.Result, err error) {
				return nil, fmt.Errorf("error")
			})
			defer mock.Unpatch()
			err := d.AddPwdLogHBase(c, a)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDao_rowKeyPwdLog(t *testing.T) {
	convey.Convey("rowKeyPwdLog", t, func() {
		res := rowKeyPwdLog(0, 0)
		convey.So(res, convey.ShouldNotBeNil)
	})
}
