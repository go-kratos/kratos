package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_InsertAvIncome(t *testing.T) {
	Convey("InsertAvIncome", t, func() {
		c := context.Background()
		d.db.Exec(c, "DELETE FROM av_income WHERE date='2018-06-01'")
		_, err := d.InsertAvIncome(c, "(123,2,6,1,'2018-06-01',100,100,100,50,'2018-06-01',100)")
		So(err, ShouldBeNil)
	})
}
