package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InsertUpIncome(t *testing.T) {
	Convey("InsertUpIncome", t, func() {
		_, err := d.InsertUpIncome(context.Background(), "(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,'2018-06-24',1,2,3,4,5,6)")
		So(err, ShouldBeNil)
	})
}
