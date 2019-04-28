package service

import (
	"context"
	"testing"

	pb "go-common/app/service/main/history/api/grpc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestService_AddHistory(t *testing.T) {
	var (
		c   = context.Background()
		arg = &pb.AddHistoryReq{Business: "pgc", Kid: 1, Mid: 3}
	)
	Convey("add data", t, func() {
		_, err := s.AddHistory(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_AddHistories(t *testing.T) {
	var (
		c   = context.Background()
		arg = &pb.AddHistoriesReq{Histories: []*pb.AddHistoryReq{{Business: "pgc", Kid: 1, Mid: 3}}}
	)
	Convey("add data", t, func() {
		_, err := s.AddHistories(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_DelHistories(t *testing.T) {
	Convey("del data", t, func() {
		c := context.Background()
		arg := &pb.DelHistoriesReq{Mid: 1, Records: []*pb.DelHistoriesReq_Record{{Business: "pgc", ID: 1}}}
		_, err := s.DelHistories(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_ClearHistory(t *testing.T) {
	Convey("clear data", t, func() {
		c := context.Background()
		arg := &pb.ClearHistoryReq{Mid: 1}
		_, err := s.ClearHistory(c, arg)
		So(err, ShouldBeNil)
	})
}

func TestService_UserHistories(t *testing.T) {
	Convey("get user histories", t, func() {
		c := context.Background()
		arg := &pb.UserHistoriesReq{Mid: 1, Ps: 10}
		gotReply, err := s.UserHistories(c, arg)
		So(err, ShouldBeNil)
		So(gotReply, ShouldNotBeEmpty)
	})
}

func TestService_Histories(t *testing.T) {
	Convey("get histories", t, func() {
		c := context.Background()
		arg := &pb.HistoriesReq{Mid: 1, Business: "pgc", Kids: []int64{1}}
		gotReply, err := s.Histories(c, arg)
		So(err, ShouldBeNil)
		So(gotReply, ShouldNotBeEmpty)
	})
}

// func TestService_FlushCache(t *testing.T) {
// 	Convey("get histories", t, func() {
// 		c := context.Background()
// 		arg := &pb.FlushCacheReq{Merges: []*model.Merge{{Mid: 1, Bid: 4}}}
// 		_, err := s.FlushCache(c, arg)
// 		So(err, ShouldBeNil)
// 	})
// }
