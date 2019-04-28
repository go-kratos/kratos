package like

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/activity/model/like"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_SubjectInitialize(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		err := svf.SubjectInitialize(context.Background(), 10292)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}

func TestService_SubjectMaxIDInitialize(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		err := svf.SubjectMaxIDInitialize(context.Background())
		So(err, ShouldBeNil)
	}))
}

func TestService_SubjectLikeListInitialize(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		err := svf.SubjectLikeListInitialize(context.Background(), 10256)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}

func TestService_LikeActCountInitialize(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		err := svf.LikeActCountInitialize(context.Background(), 10256)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}

func TestService_SubjectUp(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		err := svf.SubjectUp(context.Background(), 10256)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
	}))
}

func TestService_ActSubject(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		data, err := svf.ActSubject(context.Background(), 10340)
		time.Sleep(time.Second)
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}

func TestService_ActProtocol(t *testing.T) {
	Convey("should return without err", t, WithService(func(svf *Service) {
		data, _ := svf.ActProtocol(context.Background(), &like.ArgActProtocol{Sid: 10274})
		time.Sleep(time.Second)
		Printf("%+v %+v", data, data.Protocol)
	}))
}
