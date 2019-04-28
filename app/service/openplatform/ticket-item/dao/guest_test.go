package dao

import (
	"context"
	"testing"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"

	. "github.com/smartystreets/goconvey/convey"
)

// TestAddGuest
func TestDao_AddGuest(t *testing.T) {
	Convey("AddGuest", t, func() {
		once.Do(startService)
		res, err := d.AddGuest(context.TODO(), &item.GuestInfoRequest{
			Name:        "gotester",
			GuestImg:    "testerimg.jpg",
			Description: "this is a description",
		})
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestUpdateGuest
func TestDao_UpdateGuest(t *testing.T) {
	Convey("UpdateGuest", t, func() {
		once.Do(startService)
		res, err := d.UpdateGuest(context.TODO(), &item.GuestInfoRequest{
			ID:          55,
			Name:        "gotester22",
			GuestImg:    "testerimg.jpg2222",
			Description: "this is a description222",
		})
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestGuestStatus
func TestDao_GuestStatus(t *testing.T) {
	Convey("GuestStatus", t, func() {
		once.Do(startService)
		res, err := d.GuestStatus(context.TODO(), 57, 1)
		So(res, ShouldBeTrue)
		So(err, ShouldBeNil)
	})
}

// TestGetGuests
func TestDao_GetGuests(t *testing.T) {
	Convey("GetGuest", t, func() {
		once.Do(startService)
		res, err := d.GetGuests(context.TODO(), 60)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
