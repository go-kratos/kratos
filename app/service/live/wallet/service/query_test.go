package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"testing"
)

const RandSource = "abcdefghijklmnopq12345678"

func testGetRandString(size int) string {
	res := ""

	sourceLen := int64(len(RandSource))

	for i := 0; i < size; i++ {
		index := r.Int63n(sourceLen)
		res = res + RandSource[index:index+1]
	}
	return res
}

func queryQuery(t *testing.T, tid string) (success bool, resp *model.QueryResp, err error) {

	platform := "pc" // 对于query接口platform目前不起作用
	bp := getTestDefaultBasicParam(tid)
	v, err := s.Query(ctx, bp, 0, platform, tid)
	if err == nil {
		success = true
		resp = v.(*model.QueryResp)
	}
	return

}

func queryQueryWithUid(t *testing.T, tid string, uid int64) (success bool, resp *model.QueryResp, err error) {

	platform := "pc" // 对于query接口platform目前不起作用
	bp := getTestDefaultBasicParam(tid)
	v, err := s.Query(ctx, bp, uid, platform, tid)
	if err == nil {
		success = true
		resp = v.(*model.QueryResp)
	}
	return

}

func TestService_Query(t *testing.T) {
	Convey("normal", t, testWith(func() {
		tid := getTestTidForCall(t, 0)

		success, resp, err := queryQueryWithUid(t, tid, 1)
		if !success {
			So(err, ShouldEqual, ecode.NothingFound)
		} else {
			So(resp.Status == TX_STATUS_FAILED || resp.Status == TX_STATUS_SUCC, ShouldBeTrue)
		}

	}))

	Convey("tid invalid", t, testWith(func() {
		tid := testGetRandString(32)

		success, resp, err := queryQuery(t, tid)
		So(success, ShouldBeFalse)
		So(resp, ShouldBeNil)
		So(err, ShouldEqual, ecode.ServerErr)
	}))
}
