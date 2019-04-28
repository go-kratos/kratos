package v1

import (
	"context"
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	v1pb "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/conf"
	"go-common/app/service/live/dao-anchor/dao"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func TestLiveRoomTag(t *testing.T) {

	flag.Set("conf", "../../../cmd/test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)

	Convey("Test consume room tag", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalLiveRoomTag(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"tag_id":3,
				"tag_sub_id":2,
				"tag_value":1,
				"tag_ext":"",
				"expire_time":1549542517
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.RoomByIDsReq{
			RoomIds: []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchRoomByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.RoomDataSet), ShouldEqual, 1)
		So(resp.RoomDataSet[460654], ShouldNotBeNil)
		So(len(resp.RoomDataSet[460654].TagList), ShouldBeGreaterThanOrEqualTo, 1)
		So(resp.RoomDataSet[460654].TagList, ShouldContain, &v1pb.TagData{
			TagId:       3,
			TagSubId:    2,
			TagValue:    1,
			TagExt:      "",
			TagExpireAt: 1549542517,
		})
	})

	Convey("Update", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalLiveRoomTag(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"tag_id":3,
				"tag_sub_id":1,
				"tag_value":2,
				"tag_ext":"asd",
				"expire_time":1549000000
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.RoomByIDsReq{
			RoomIds: []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchRoomByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.RoomDataSet), ShouldEqual, 1)
		So(resp.RoomDataSet[460654], ShouldNotBeNil)
		So(len(resp.RoomDataSet[460654].TagList), ShouldBeGreaterThanOrEqualTo, 1)
		So(resp.RoomDataSet[460654].TagList, ShouldContain, &v1pb.TagData{
			TagId:       3,
			TagSubId:    1,
			TagValue:    2,
			TagExt:      "asd",
			TagExpireAt: 1549000000,
		})
	})
}
