package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiAddUpFilters(t *testing.T) {
	Convey("", t, func() {
		tx, err := testDao.BeginBiliDMTrans(context.TODO())
		So(err, ShouldBeNil)
		So(tx, ShouldNotBeNil)
		Convey("", func() {
			_, err := testDao.UpdateUpFilter(tx, 16299551, 33, 0)
			So(err, ShouldBeNil)
			if err != nil {
				_, err = testDao.UpdateUpFilterCnt(tx, 162, 1, 1, 1000)
				So(err, ShouldBeNil)
			}
			tx.Commit()
		})
	})
}
