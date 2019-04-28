package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestGetUserCards(t *testing.T) {
	convey.Convey("GetUserCards", t, func() {
		cards, _ := d.GetUserCards(context.TODO(), []int64{27515405, 27515305})
		convey.So(len(cards.Cards), convey.ShouldEqual, 2)
		convey.So(cards.Cards[27515305].Name, convey.ShouldEqual, "Test000044")
	})
}
