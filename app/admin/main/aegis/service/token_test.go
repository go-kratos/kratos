package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/net"

	"github.com/smartystreets/goconvey/convey"
	"time"
)

func TestServiceGetTokenList(t *testing.T) {
	var (
		c  = context.TODO()
		pm = &net.ListTokenParam{
			NetID: 1,
		}
	)
	convey.Convey("GetTokenList", t, func(ctx convey.C) {
		result, err := s.GetTokenList(c, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)
		})
	})
}

func TestServiceShowToken(t *testing.T) {
	var (
		c  = context.TODO()
		id = int64(20)
	)
	convey.Convey("ShowToken", t, func(ctx convey.C) {
		n, err := s.ShowToken(c, id)
		ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(n, convey.ShouldNotBeNil)
		})
	})
}

func TestServicecheckTokenUnique(t *testing.T) {
	var (
		netID   = int64(0)
		name    = ""
		compare = int8(0)
		value   = ""
	)
	convey.Convey("checkTokenUnique", t, func(ctx convey.C) {
		err, msg := s.checkTokenUnique(cntx, netID, name, compare, value)
		ctx.Convey("Then err should be nil.msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceAddToken(t *testing.T) {
	var (
		c  = context.TODO()
		no = &net.Token{
			NetID:   1,
			ChName:  "创建时间",
			Name:    "time",
			Compare: net.TokenCompareAssign,
			Value:   time.Now().String(),
			Type:    net.TokenTypeInt16,
		}
	)
	convey.Convey("AddToken", t, func(ctx convey.C) {
		id, err, msg := s.AddToken(c, no)
		ctx.Convey("Then err should be nil.id,msg should not be nil.", func(ctx convey.C) {
			ctx.So(msg, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(id, convey.ShouldNotBeNil)
		})
	})
}

func TestServicefetchOldBindAndLog(t *testing.T) {
	convey.Convey("fetchOldBindAndLog", t, func(ctx convey.C) {
		all, logs, err := s.fetchOldBindAndLog(cntx, 3, net.BindTranType)
		for _, item := range all {
			t.Logf("s.fetchOldBindAndLog all item(%+v)", item)
		}
		t.Logf("logs(%v)", logs)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServicecompareTransitionBind(t *testing.T) {
	var (
		elementID = int64(1)
		binds     = []*net.TokenBindParam{
			{ID: 26, TokenID: "20", ChName: "批量通过"},
			{ID: 24, TokenID: "21,20", ChName: "可以通过"},
		}
	)

	convey.Convey("compareTransitionBind", t, func(ctx convey.C) {
		tx, _ := s.gorm.BeginTx(context.Background())
		diff, _, err, msg := s.compareTranBind(cntx, tx, elementID, binds, true)
		t.Logf("s.compareTranBind diff(%v) msg(%s)", diff, msg)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
