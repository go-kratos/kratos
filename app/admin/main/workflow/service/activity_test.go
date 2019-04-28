package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/workflow/model/param"

	"github.com/smartystreets/goconvey/convey"
)

func TestActivityList(t *testing.T) {
	convey.Convey("ActivityList", t, func() {
		eid, err := s.AddEvent(context.Background(), &param.EventParam{
			Cid:         int64(1),
			AdminID:     int64(1),
			Content:     "test.content",
			Attachments: "test.attachments",
			Event:       int8(1),
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(eid, convey.ShouldNotEqual, 0)

		elist, err := s.ListEvent(context.Background(), 1)
		convey.So(err, convey.ShouldBeNil)
		eids := make([]int64, 0)
		for _, e := range elist {
			eids = append(eids, e.Eid)
		}
		convey.So(eid, convey.ShouldBeIn, eids)

		acts, err := s.ActivityList(context.Background(), int8(1), 1)
		convey.So(err, convey.ShouldBeNil)
		convey.So(acts, convey.ShouldNotBeNil)
		convey.So(acts.Events, convey.ShouldNotBeEmpty)
		convey.So(acts.Logs, convey.ShouldNotBeEmpty)

		eids2 := make([]int64, 0)
		for _, e := range acts.Events {
			eids2 = append(eids2, e.Eid)
		}
		convey.So(eid, convey.ShouldBeIn, eids2)
	})
}
