package dao

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//  go test  -test.v -test.run TestDaoNotGrantActOrders
func TestDaoNotGrantActOrders(t *testing.T) {
	Convey("TestDaoNotGrantActOrders salary coupon", t, func() {
		res, err := d.NotGrantActOrders(context.Background(), "ele", 100)
		for _, v := range res {
			fmt.Println("res:", v)
		}
		So(err, ShouldBeNil)
	})
}
