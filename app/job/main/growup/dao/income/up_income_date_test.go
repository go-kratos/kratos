package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GetUpIncomeTable(t *testing.T) {
	Convey("GetUpIncomeTable", t, func() {
		_, err := d.GetUpIncomeTable(context.Background(), "up_income_weekly", "2018-06-01", 0, 2000)
		So(err, ShouldBeNil)
	})
}

func Test_InsertUpIncomeTable(t *testing.T) {
	Convey("InsertUpIncomeTable", t, func() {
		_, err := d.InsertUpIncomeTable(context.Background(), "up_income_weekly", "(1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,'2018-06-24',1,2,3,4)")
		So(err, ShouldBeNil)
	})
}
