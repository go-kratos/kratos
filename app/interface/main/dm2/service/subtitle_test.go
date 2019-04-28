package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-common/app/interface/main/dm2/model"
	account "go-common/app/service/main/account/api"
	rpcmodel "go-common/app/service/main/member/model/block"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSaveSubtitleDraft(t *testing.T) {
	Convey("", t, func() {
		var aid int64 = 10098493
		var oid int64 = 401
		var tp int32 = 1
		mid := int64(88888929)
		lan := "zh-CN"
		body := &model.SubtitleBody{
			FontColor:       "#FFFFFF",
			FontSize:        0.4,
			BackgroundAlpha: 0.5,
			BackgroundColor: "#9C27B0",
			Stroke:          "none",
		}
		items := make([]*model.SubtitleItem, 0, 10)
		for i := 0; i < 10; i++ {
			items = append(items, &model.SubtitleItem{
				From:     float64(i * 10),
				To:       float64(i*10 + 5),
				Location: uint8(8),
				Content:  fmt.Sprintf("test_1133331seg_%d", i+1),
			})
		}
		body.Bodys = items
		bs, err := json.Marshal(&body)
		if err != nil {
			return
		}
		_, err = svr.SaveSubtitleDraft(context.Background(), aid, oid, tp, mid, lan, true, true, 0, bs)
		time.Sleep(time.Second)
		t.Logf("err:%v", err)
		So(err, ShouldBeNil)
	})
}

func TestDelSubtitle(t *testing.T) {
	Convey("", t, func() {
		var oid int64 = 101
		var subtitleID = int64(6)
		mid := int64(5)

		err := svr.DelSubtitle(context.Background(), oid, subtitleID, mid)
		So(err, ShouldBeNil)
	})
}

func TestAuditSubtitle(t *testing.T) {
	Convey("", t, func() {
		var oid int64 = 101
		mid := int64(5)

		err := svr.AuditSubtitle(context.Background(), oid, 1, mid, true, "")
		So(err, ShouldBeNil)
	})
}

func TestSubtitleShow(t *testing.T) {
	Convey("", t, func() {
		var oid int64 = 101
		var aid int64 = 10098493

		subtitleID := int64(4)
		start := time.Now()
		res, err := svr.SubtitleShow(context.Background(), aid, oid, subtitleID)
		t.Logf("costing:%v", time.Since(start))
		So(err, ShouldBeNil)
		t.Logf("%+v:", res)
	})
}

func TestSubtitleLock(t *testing.T) {
	Convey("", t, func() {
		var oid int64 = 101
		var tp int32 = 1

		mid := int64(5)
		subtitleID := int64(5)
		start := time.Now()
		err := svr.SubtitleLock(context.Background(), oid, tp, mid, subtitleID, true)
		t.Logf("costing:%v", time.Since(start))
		So(err, ShouldBeNil)
	})
}

func TestSearchAuthor(t *testing.T) {
	Convey("", t, func() {
		var mid int64 = 27515615
		var status int32
		res, err := svr.SearchAuthor(context.Background(), mid, status, 1, 10)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("page:%+v", res)
		t.Logf("err:%v", err)
		// for _, rs := range res.Subtitles {
		// 	t.Logf("rs:%+v", rs)
		// }
	})
}

func TestSearchAssist(t *testing.T) {
	Convey("search assist", t, func() {
		var mid int64 = 27515615
		var status int32
		res, err := svr.SearchAssist(context.Background(), 0, 0, 1, mid, status, 1, 50)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("page:%+v", res)
		t.Logf("err:%v", err)
		for _, rs := range res.Subtitles {
			t.Logf("rs:%+v", rs)
		}
	})
}

func TestBlack(t *testing.T) {
	Convey("", t, func() {
		res, err := svr.accountRPC.Blacks3(context.Background(), &account.MidReq{
			Mid: 27515266,
		})
		t.Logf("err:%v", err)
		t.Logf("blockINfo:%+v", res)

	})
}

func TestBlockInfo(t *testing.T) {
	Convey("blokc ", t, func() {
		res, err := svr.memberRPC.BlockInfo(context.Background(), &rpcmodel.RPCArgInfo{
			MID: 27515256,
		})
		t.Logf("err:%v", err)
		t.Logf("blockINfo:%+v", res)

	})
}
