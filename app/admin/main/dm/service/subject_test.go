package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchiveList(t *testing.T) {
	req := &model.ArchiveListReq{
		ID:     34261,
		IDType: "oid",
		Sort:   "desc",
		Order:  "mtime",
		Pn:     1,
		Ps:     50,
		Page:   int64(model.CondIntNil),
		State:  int64(model.CondIntNil),
	}
	convey.Convey("test last archive list", t, func() {
		res, err := svr.ArchiveList(context.TODO(), req)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeEmpty)
		t.Logf("===%+v", res.Page)
		for _, v := range res.ArcLists {
			t.Logf("===%+v", v)
		}
	})
}

func TestServiceUptSubjectsState(t *testing.T) {
	convey.Convey("UptSubjectsState", t, func() {
		err := svr.UptSubjectsState(context.TODO(), 1, 111, "test", []int64{1221}, 1, "aaaaa")
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestServiceUpSubjectMaxLimit(t *testing.T) {
	convey.Convey("UpSubjectMaxLimit", t, func() {
		err := svr.UpSubjectMaxLimit(context.TODO(), 1, 10131812, 333)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestServiceSubjectLog(t *testing.T) {
	convey.Convey("SubjectLog", t, func() {
		data, err := svr.SubjectLog(context.TODO(), 1, 1221)
		convey.So(err, convey.ShouldBeNil)
		convey.So(data, convey.ShouldNotBeNil)
		for _, v := range data {
			t.Logf("====%+v", v)
		}
	})
}
