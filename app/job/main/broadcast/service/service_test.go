package service

import (
	"testing"

	"go-common/app/job/main/broadcast/conf"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewService(t *testing.T) {
	Convey("TestNewService", t, func() {
		s := New(conf.Conf)
		So(s, ShouldNotBeNil)
		Convey("TestPush", func() {
			err := s.pushMsg([]byte("test"))
			So(err, ShouldBeNil)
		})
		Convey("TestPushComet", func() {
			s.pushKeys(1, "test", nil, []byte("test"), 0)
		})
		Convey("TestBroadCast", func() {
			s.broadcast(1, []byte("test"), 100, "", 0)
		})
		Convey("TestBroadCast", func() {
			s.broadcastRoomRawBytes("1", []byte("test"))
		})
		Convey("TestNewRoom", func() {
			r := NewRoom(s, "test", RoomOptions{})
			So(r, ShouldNotBeNil)
			Convey("TestNewPush", func() {
				err := r.Push(7, []byte("test"), 0)
				So(err, ShouldBeNil)
			})
		})

	})
}
