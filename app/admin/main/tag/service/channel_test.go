package service

// import (
// 	"context"
// 	"testing"

// 	. "github.com/smartystreets/goconvey/convey"
// )

// func TestServiceChannel(t *testing.T) {
// 	var (
// 		id             = int64(22)
// 		operator       = "unit test"
// 		order          = "ctime"
// 		sort           = "DESC"
// 		tp             = int32(2)
// 		ps       int32 = 1
// 		pn       int32 = 20
// 	)
// 	Convey("ChanneList", func() {
// 		testSvc.ChanneList(context.TODO(), []int64{id}, operator, order, sort, "", "", tp, -1, -1, pn, ps)
// 	})
// 	Convey("ChannelState", func() {
// 		testSvc.ChannelState(context.TODO(), id, 0)
// 	})
// 	Convey("ChannelCategory", func() {
// 		testSvc.ChannelCategory(context.TODO())
// 	})
// 	Convey("CategoryAdd", func() {
// 		testSvc.CategoryAdd(context.TODO(), "unit")
// 	})
// 	Convey("DeleteCategory", func() {
// 		testSvc.DeleteCategory(context.TODO(), id)
// 	})
// 	Convey("StateCategory", func() {
// 		testSvc.DeleteCategory(context.TODO(), id)
// 	})
// 	Convey("DeletALlChanRule", func() {
// 		testSvc.DeletALlChanRule(context.TODO(), id)
// 	})
// }
