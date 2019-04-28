package http

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/apm/model/need"
)

var (
	_nadduri      = "%s/x/admin/apm/need/add"
	_nlisturi     = "%s/x/admin/apm/need/list"
	_nedituri     = "%s/x/admin/apm/need/edit"
	_nverifyuri   = "%s/x/admin/apm/need/verify"
	_nthumsupuri  = "%s/x/admin/apm/need/thumbsup"
	_nvotelisturi = "%s/x/admin/apm/need/vote/list"
)

func TestNeedList(t *testing.T) {
	convey.Convey("", t, func() {
		params := url.Values{}
		params.Set("pn", "1")
		params.Set("ps", "5")
		res := new(struct {
			Code    int             `json:"code"`
			Message string          `json:"message"`
			Data    *need.NListResp `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_nlisturi, _domain), "", params, &res)
		t.Logf("res%+v", res.Data)

	})
}

func TestNeedAdd(t *testing.T) {
	convey.Convey("TestNeedAdd normal", t, func() {
		params := url.Values{}
		params.Set("title", "dsds")
		params.Set("content", "sds")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_nadduri, _domain), "", params, &res)
		t.Logf("res%+v", res)
		convey.So(res.Code, convey.ShouldEqual, 0)
	})

	//convey.Convey("TestNeedAdd params error", t, func() {
	//	//	params := url.Values{}
	//	//	params.Set("title", "提一个小需求阿斯加德卡萨")
	//	//	res := Response{}
	//	//	_ = requests("POST", fmt.Sprintf(_nadduri, _domain), "", params, &res)
	//	//	t.Logf("res%+v", res)
	//	//	convey.So(res.Code, convey.ShouldEqual, -400)
	//	//})
}

func TestNeedEdit(t *testing.T) {
	convey.Convey("TestNeedEdit", t, func() {
		params := url.Values{}
		params.Set("title", "fss")
		params.Set("content", "fss")
		params.Set("id", "26")

		res := Response{}
		_ = requests("POST", fmt.Sprintf(_nedituri, _domain), "", params, &res)
		convey.So(res.Code, convey.ShouldEqual, 0)

	})
}
func TestNeedVerify(t *testing.T) {
	convey.Convey("TestNeedVerify", t, func() {
		params := url.Values{}
		params.Set("status", "2")
		params.Set("id", "28")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_nverifyuri, _domain), "", params, &res)
		convey.So(res.Code, convey.ShouldEqual, 70018)

	})
}

func TestNeedThumbsUp(t *testing.T) {
	convey.Convey("TestNeedThumbsUp", t, func() {
		params := url.Values{}
		params.Set("req_id", "29")
		params.Set("like_type", "1")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_nthumsupuri, _domain), "", params, &res)
		convey.So(res.Code, convey.ShouldEqual, 0)

	})
}

func TestVoteList(t *testing.T) {
	convey.Convey("TestVoteList", t, func() {
		params := url.Values{}
		params.Set("req_id", "11")
		params.Set("like_type", "1")
		res := new(struct {
			Code    int                `json:"code"`
			Message string             `json:"message"`
			Data    *need.VoteListResp `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_nvotelisturi, _domain), "", params, &res)
		t.Logf("res%+v", res)
	})
}
