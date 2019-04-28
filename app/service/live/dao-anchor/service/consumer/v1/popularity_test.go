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

func TestPopularityStatistics(t *testing.T) {

	flag.Set("conf", "../../../cmd/test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)

	Convey("Test consume real pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		// 测试实时人气值
		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":3000,
				"cycle":60,
				"type":88888
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
		So(resp.RoomDataSet[460654].PopularityCount, ShouldEqual, 3000)
	})

	Convey("Test update real pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		// 测试实时人气值
		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":6000,
				"cycle":60,
				"type":88888
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
		So(resp.RoomDataSet[460654].PopularityCount, ShouldEqual, 6000)
	})

	Convey("Test consume 7 day pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":5000,
				"cycle":604800,
				"type":88888
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    1,
			AttrSubId: 2,
			RoomIds:   []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 1)
		So(resp.Attrs[460654], ShouldNotBeNil)
		So(resp.Attrs[460654].AttrId, ShouldEqual, 1)
		So(resp.Attrs[460654].AttrSubId, ShouldEqual, 2)
		So(resp.Attrs[460654].AttrValue, ShouldEqual, 5000)
	})

	Convey("Test update 7 day pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":88000,
				"cycle":604800,
				"type":88888
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    1,
			AttrSubId: 2,
			RoomIds:   []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 1)
		So(resp.Attrs[460654], ShouldNotBeNil)
		So(resp.Attrs[460654].AttrId, ShouldEqual, 1)
		So(resp.Attrs[460654].AttrSubId, ShouldEqual, 2)
		So(resp.Attrs[460654].AttrValue, ShouldEqual, 88000)
	})

	Convey("Test consume 30 day pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":5000,
				"cycle":2592000,
				"type":88888
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    1,
			AttrSubId: 3,
			RoomIds:   []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 1)
		So(resp.Attrs[460654], ShouldNotBeNil)
		So(resp.Attrs[460654].AttrId, ShouldEqual, 1)
		So(resp.Attrs[460654].AttrSubId, ShouldEqual, 3)
		So(resp.Attrs[460654].AttrValue, ShouldEqual, 5000)
	})

	Convey("Test update 30 day pc", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		err := s.internalPopularityStatistics(ctx, &databus.Message{
			Value: []byte(`	
			{
				"msg_id":"1111",
				"room_id":460654,
				"value":88000,
				"cycle":2592000,
				"type":88888
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    1,
			AttrSubId: 3,
			RoomIds:   []int64{460654},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 1)
		So(resp.Attrs[460654], ShouldNotBeNil)
		So(resp.Attrs[460654].AttrId, ShouldEqual, 1)
		So(resp.Attrs[460654].AttrSubId, ShouldEqual, 3)
		So(resp.Attrs[460654].AttrValue, ShouldEqual, 88000)
	})
}
