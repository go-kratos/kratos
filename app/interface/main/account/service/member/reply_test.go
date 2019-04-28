package member

import (
	"context"
	"testing"

	"go-common/app/interface/main/account/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestService_ReplyHistoryList(t *testing.T) {
	convey.Convey("ReplyHistoryList", t, func() {
		var (
			err error
			//ip              = ""
			mid       int64 = 88889069
			stime           = "1500805318"
			etime           = "1511237870"
			order           = "like"
			sort            = "desc"
			pn        int64 = 1
			ps        int64 = 100
			accessKey       = ""
			cookie          = ""
			rhl       *model.ReplyHistory
		)
		if rhl, err = s.ReplyHistoryList(context.TODO(), mid, stime, etime, order, sort, pn, ps, accessKey, cookie); err != nil {
			convey.So(err, convey.ShouldBeNil)
			t.Logf("err: %v", err)
		}
		for _, v := range rhl.Records {
			t.Logf("title:%s, url:%s, id:%d, type: %d", v.Title, v.URL, v.ID, v.Type)
		}
	})
}
