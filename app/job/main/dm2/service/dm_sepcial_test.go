package service

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// func TestSpecialDmRemove(t *testing.T) {
// 	Convey("", t, func() {
// 		dmid := int64(920249977)
// 		dm := &model.DM{
// 			ID:       dmid,
// 			Type:     1,
// 			Oid:      19,
// 			Mid:      1,
// 			State:    1,
// 			Pool:     2,
// 			Progress: 10,
// 		}
// 		_, err := testSvc.dao.UpdateDM(context.TODO(), dm)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		So(err, ShouldBeNil)
// 	})
// }

func TestSpecialLocationUpdate(t *testing.T) {
	Convey("", t, func() {

		err := svr.specialLocationUpdate(context.TODO(), 1, 19)
		if err != nil {
			fmt.Println(err)
		}
		So(err, ShouldBeNil)
	})
}
