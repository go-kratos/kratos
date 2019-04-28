package service

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"go-common/app/job/main/videoup-report/model/archive"
)

func TestService_Addarchive(t *testing.T) {
	Convey("addarchive", t, func() {
		err := s.addArchive(context.Background(), &archive.VideoupMsg{
			Aid: 10,
		})
		t.Logf("err(%+v)", err)
	})
}

func TestService_Modifyarchive(t *testing.T) {
	Convey("Modifyarchive", t, func() {
		err := s.modifyArchive(context.Background(), &archive.VideoupMsg{
			Aid:       11,
			TagChange: true,
			AddVideos: true,
		})
		t.Logf("err(%+v)", err)
	})
}

func TestService_arcStateChange(t *testing.T) {
	var (
		err   error
		state int64
		aid   = int64(12)
	)
	Convey("arcStateChange", t, func() {
		a, _ := s.arc.ArchiveByAid(context.TODO(), aid)
		nw := &archive.Archive{
			ID:    a.ID,
			State: archive.StateOpen,
		}
		old := &archive.Archive{
			ID:    a.ID,
			State: archive.StateForbidRecicle,
		}
		err = s.dataDao.CloseReply(context.TODO(), a.ID, a.Mid)
		So(err, ShouldBeNil)

		//只允许关的情况下，状态联动从关-》开不起作用
		err = s.arcStateChange(nw, old, false)
		So(err, ShouldBeNil)
		state, err = s.dataDao.CheckReply(context.TODO(), nw.ID)
		So(err, ShouldBeNil)
		So(state, ShouldEqual, archive.ReplyOff)

		//只允许开的情况下，状态联动从关->开起作用
		err = s.arcStateChange(nw, old, true)
		So(err, ShouldBeNil)
		state, err = s.dataDao.CheckReply(context.TODO(), nw.ID)
		So(err, ShouldBeNil)
		So(state, ShouldEqual, archive.ReplyOn)

		//只允许开的情况下，状态联动从开->关不起作用
		err = s.arcStateChange(old, nw, true)
		So(err, ShouldBeNil)
		state, err = s.dataDao.CheckReply(context.TODO(), nw.ID)
		So(err, ShouldBeNil)
		So(state, ShouldEqual, archive.ReplyOn)

		//只允许关的情况下，状态联动从开->关起作用
		err = s.arcStateChange(old, nw, false)
		So(err, ShouldBeNil)
		state, err = s.dataDao.CheckReply(context.TODO(), nw.ID)
		So(err, ShouldBeNil)
		So(state, ShouldEqual, archive.ReplyOff)
	})
}
