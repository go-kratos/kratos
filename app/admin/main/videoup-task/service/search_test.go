package service

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"

	"go-common/app/admin/main/videoup-task/model"
)

func Test_searchQAVideo(t *testing.T) {
	Init()

	convey.Convey("search qavideo", t, func() {
		pm := &model.ListParams{
			State:     1,
			Ps:        50,
			Pn:        1,
			Sort:      "desc",
			Order:     "id",
			CTimeFrom: "2018-05-08 00:00:00",
			CTimeTo:   "2018-07-20 00:00:00",
			//FTimeFrom: "2018-07-10 00:00:00",
			//FTimeTo:   "2018-07-20 00:00:00",
			FansFrom:    1,
			FansTo:      100,
			UID:         []int64{421, 481},
			UPGroup:     []int64{1, 2},
			AuditStatus: []int{-4, -2, 10000},
			ArcTypeID:   []int64{76},
			TagID:       []int64{14, 6, 23, 15},
			Keyword:     []string{"普通"},
			TaskID:      []int64{8326, 8327},
			//Limit:       2,
		}
		resp, err := s.searchQAVideo(context.TODO(), pm)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("resp(%+v)", resp)
	})
}
