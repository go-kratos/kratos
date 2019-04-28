package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConnect(t *testing.T) {
	var (
		server    = "test_server"
		serverKey = "test_server_key"
		token     = []byte(`{"device_id":"test_server_key","business":"dm", "room_id":"test://test_room", "accepts":[1,2,3]}`)
		c         = context.Background()
	)
	Convey("connect", t, WithService(func(s *Service) {
		// connect
		mid, key, roomID, _, accepts, err := s.Connect(c, server, serverKey, "", token)
		So(err, ShouldBeNil)
		So(key, ShouldEqual, serverKey)
		So(roomID, ShouldEqual, "test://test_room")
		So(len(accepts), ShouldEqual, 3)
		t.Log(mid, key, roomID, accepts, err)
		// heartbeat
		err = s.Heartbeat(c, mid, key, server)
		So(err, ShouldBeNil)
		// disconnect
		has, err := s.Disconnect(c, mid, key, server)
		So(err, ShouldBeNil)
		So(has, ShouldEqual, true)
	}))
}
