package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/openplatform/ticket-item/model"
	"sync"
	"testing"
)

var (
	once sync.Once
	d    *Dao
	ctx  = context.TODO()
)

// Test_RawItems
func TestRawItems(t *testing.T) {
	Convey("RawItems", t, func() {
		once.Do(startService)
		res, err := d.RawItems(ctx, model.DataIDs)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res, ShouldNotBeEmpty)
	})
	Convey("RawItems", t, func() {
		once.Do(startService)
		res, err := d.RawItems(ctx, model.NoDataIDs)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeEmpty)
	})
}
