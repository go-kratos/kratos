package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/admin/main/member/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOfficials(t *testing.T) {
	convey.Convey("Officials", t, func() {
		o, total, err := s.Officials(context.Background(), &model.ArgOfficial{
			Mid:   123,
			Role:  []int64{1},
			ETime: xtime.Time(time.Now().Unix()),
			Pn:    1,
			Ps:    20,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(total, convey.ShouldBeGreaterThan, 0)
	})
}

func TestOfficialDoc(t *testing.T) {
	convey.Convey("OfficialDoc", t, func() {
		o, logs, block, spy, realname, mids, err := s.OfficialDoc(context.Background(), 123)
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(logs, convey.ShouldNotBeNil)
		convey.So(block, convey.ShouldNotBeNil)
		convey.So(spy, convey.ShouldNotBeNil)
		convey.So(realname, convey.ShouldNotBeNil)
		convey.So(mids, convey.ShouldNotBeNil)
	})
}

func TestOfficialDocs(t *testing.T) {
	convey.Convey("OfficialDocs", t, func() {
		o, total, err := s.OfficialDocs(context.Background(), &model.ArgOfficialDoc{
			Mid:   123,
			Role:  []int64{1},
			State: []int64{1},
			ETime: xtime.Time(time.Now().Unix()),
			Pn:    1,
			Ps:    20,
		})
		convey.So(err, convey.ShouldBeNil)
		convey.So(o, convey.ShouldNotBeNil)
		convey.So(total, convey.ShouldBeGreaterThan, 0)
	})
}

func TestOfficialDocAudit(t *testing.T) {
	convey.Convey("OfficialDocAudit", t, func() {
		err := s.OfficialDocAudit(context.Background(), &model.ArgOfficialAudit{
			Mid:    123,
			State:  1,
			UID:    111,
			Uname:  "guan",
			Reason: "xxx",
		})
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestOfficialDocEdit(t *testing.T) {
	convey.Convey("OfficialDocEdit", t, func() {
		err := s.OfficialDocEdit(context.Background(), &model.ArgOfficialEdit{
			Mid:   123,
			Name:  "guan",
			Role:  1,
			Title: "title",
			Desc:  "desc",
			UID:   111,
			Uname: "guan",
		})
		convey.So(err, convey.ShouldBeNil)
	})
}
