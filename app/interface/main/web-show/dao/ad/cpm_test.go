package ad

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Cpms(t *testing.T) {
	Convey("test Cpms", t, WithDao(func(d *Dao) {
		mid := int64(5187977)
		ids := []int64{121, 21, 12}
		data, err := d.Cpms(context.TODO(), mid, ids, "", "", "", "", "", "")
		So(err, ShouldBeNil)
		Printf("%+v", data)
	}))
}
