package client

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuyAdvance(t *testing.T) {
	var (
		mid  int64 = 27515260
		cid  int64 = 10107292
		mode       = "sp"
	)
	Convey("test buy advance dm", t, func() {
		arg := &model.ArgAdvance{
			Mid:  mid,
			Cid:  cid,
			Mode: mode,
		}
		err := svr.BuyAdvance(context.TODO(), arg)
		fmt.Println(err)
		So(err, ShouldBeNil)
	})
}

func TestAdvanceState(t *testing.T) {
	var (
		mid  int64 = 27515330
		cid  int64 = 10107292
		mode       = "sp"
	)
	Convey("test advance dm state", t, func() {
		arg := &model.ArgAdvance{
			Mid:  mid,
			Cid:  cid,
			Mode: mode,
		}
		res, err := svr.AdvanceState(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestAdvances(t *testing.T) {
	var (
		mid int64 = 27515260
	)
	Convey("test dm advances", t, func() {
		arg := &model.ArgMid{Mid: mid}
		res, err := svr.Advances(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

func TestPassAdvance(t *testing.T) {
	var (
		mid int64 = 7158471
		id  int64 = 2
	)
	Convey("test pass advance dm", t, func() {
		arg := &model.ArgUpAdvance{Mid: mid, ID: id}
		err := svr.PassAdvance(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestDenyAdvance(t *testing.T) {
	var (
		mid int64 = 27515615
		id  int64 = 107
	)
	Convey("test deny advance dm", t, func() {
		arg := &model.ArgUpAdvance{Mid: mid, ID: id}
		err := svr.DenyAdvance(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestCancelAdvance(t *testing.T) {
	var (
		mid int64 = 27515615
		id  int64 = 107
	)
	Convey("test cancel advance dm", t, func() {
		arg := &model.ArgUpAdvance{Mid: mid, ID: id}
		err := svr.CancelAdvance(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}
