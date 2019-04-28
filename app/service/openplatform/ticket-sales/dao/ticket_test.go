package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTicketsByID(t *testing.T) {
	Convey("TestTicketsByID", t, func() {
		tickets, err := d.TicketsByID(ctx, []int64{101, 102})
		So(err, ShouldBeNil)
		So(tickets, ShouldNotBeNil)
	})
}

func TestTicketSend(t *testing.T) {
	Convey("TestTicketSend", t, func() {
		sends, err := d.TicketSend(ctx, []int64{402, 412}, "send")
		So(err, ShouldBeNil)
		So(sends, ShouldNotBeNil)
	})
}
