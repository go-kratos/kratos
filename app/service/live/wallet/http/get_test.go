package http

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"net/url"
	"strconv"
	"testing"
)

type GetRes struct {
	Code int                           `json:"code"`
	Resp *model.MelonseedWithMetalResp `json:"data"`
}

type StatusRes struct {
	Code int              `json:"code"`
	Resp *model.QueryResp `json:"data"`
}

type GetAllRes struct {
	Code int                        `json:"code"`
	Resp *model.DetailWithMetalResp `json:"data"`
}

func queryGet(t *testing.T, uid int64, platform string) *GetRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", uid))
	req, _ := client.NewRequest("GET", _getURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res GetRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func queryStatus(t *testing.T, uid int64, tid string) *StatusRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", uid))
	params.Set("transaction_id", tid)
	req, _ := client.NewRequest("GET", _queryURL, "127.0.0.1", params)
	req.Header.Set("platform", "pc")

	var res StatusRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func queryGetAll(t *testing.T, uid int64, platform string) *GetAllRes {
	params := url.Values{}
	params.Set("uid", fmt.Sprintf("%d", uid))
	req, _ := client.NewRequest("GET", _getAllURL, "127.0.0.1", params)
	req.Header.Set("platform", platform)

	var res GetAllRes

	err := client.Do(context.TODO(), req, &res)
	if err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	return &res
}

func getTestWallet(t *testing.T, uid int64, platform string) *model.MelonseedWithMetalResp {
	res := queryGet(t, uid, platform)
	if res.Code != 0 {
		t.Errorf("get wallet failed uid : %d, code :%d", uid, res.Code)
		t.FailNow()
	}
	return res.Resp

}

/*
useless now
func getTestWalletDetail(t *testing.T, uid int64, platform string) *model.DetailWithMetalResp {
	res := queryGetAll(t, uid, platform)
	if res.Code != 0 {
		t.Errorf("get wallet failed uid : %d, code :%d", uid, res.Code)
		t.FailNow()
	}
	return res.Resp

}*/

func TestGet(t *testing.T) {
	once.Do(startHTTP)
	Convey("get normal", t, func() {
		res := queryGet(t, 1, "pc")
		So(res.Code, ShouldEqual, 0)
		melon := res.Resp
		coin, err := strconv.Atoi(melon.Gold)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)
		coin, err = strconv.Atoi(melon.Silver)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)
	})

	Convey("uid params error", t, func() {
		res := queryGet(t, -1, "pc")
		So(res.Code, ShouldEqual, ecode.RequestErr)
	})

	Convey("platform params error", t, func() {
		res := queryGet(t, 1, "pc1")
		So(res.Code, ShouldEqual, ecode.RequestErr)
	})

}

func TestGetAll(t *testing.T) {
	once.Do(startHTTP)
	Convey("normal", t, func() {
		res := queryGetAll(t, 1, "pc")
		t.Logf("all:%v", res)
		So(res.Code, ShouldEqual, 0)
		melon := res.Resp
		coin, err := strconv.Atoi(melon.Gold)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)
		coin, err = strconv.Atoi(melon.Silver)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)

		coin, err = strconv.Atoi(melon.SilverPayCnt)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)

		coin, err = strconv.Atoi(melon.GoldPayCnt)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)

		coin, err = strconv.Atoi(melon.GoldRechargeCnt)
		So(err, ShouldBeNil)
		So(coin, ShouldBeGreaterThan, -1)

		So(melon.CostBase, ShouldBeGreaterThanOrEqualTo, 0)

	})
}
