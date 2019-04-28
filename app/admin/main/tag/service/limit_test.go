package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceLimit(t *testing.T) {
	var (
		mid   int64 = 35152246
		cname       = "shunza"
		ps    int32 = 1
		pn    int32 = 20
	)
	Convey("LimitUserDel", func() {
		testSvc.LimitUserDel(context.TODO(), mid)
	})
	Convey("LimitUserAdd", func() {
		testSvc.LimitUserAdd(context.TODO(), mid, cname)
	})
	Convey("LimitUsers", func() {
		testSvc.LimitUsers(context.TODO(), ps, pn)
	})
}
