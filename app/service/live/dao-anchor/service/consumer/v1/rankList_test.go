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

func TestLiveRankList(t *testing.T) {

	flag.Set("conf", "../../../cmd/test.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)

	Convey("Test consume ranklist", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		// 测试小时榜
		err := s.internalLiveRankList(ctx, &databus.Message{
			Value: []byte(`	
			{
				"rank_id":1,
				"rank_type":"RANK_HOUR",
				"rank_list":{"460654":1,"460785":2},
				"expire_time":1549542517
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    4,
			AttrSubId: 1,
			RoomIds:   []int64{460654, 460785},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 2)
		So(resp.Attrs[460654], ShouldNotBeNil)
		So(resp.Attrs[460654].AttrId, ShouldEqual, 4)
		So(resp.Attrs[460654].AttrSubId, ShouldEqual, 1)
		So(resp.Attrs[460654].AttrValue, ShouldEqual, 1)
		So(resp.Attrs[460785], ShouldNotBeNil)
		So(resp.Attrs[460785].AttrId, ShouldEqual, 4)
		So(resp.Attrs[460785].AttrSubId, ShouldEqual, 1)
		So(resp.Attrs[460785].AttrValue, ShouldEqual, 2)
	})

	Convey("Test update ranklist", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		// 测试小时榜
		err := s.internalLiveRankList(ctx, &databus.Message{
			Value: []byte(`	
			{
				"rank_id":1,
				"rank_type":"RANK_HOUR",
				"rank_list":{"460767":1,"460796":2},
				"expire_time":1549542517
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    4,
			AttrSubId: 1,
			RoomIds:   []int64{460654, 460785, 460767, 460796},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 2)
		So(resp.Attrs[460767], ShouldNotBeNil)
		So(resp.Attrs[460767].AttrId, ShouldEqual, 4)
		So(resp.Attrs[460767].AttrSubId, ShouldEqual, 1)
		So(resp.Attrs[460767].AttrValue, ShouldEqual, 1)
		So(resp.Attrs[460796], ShouldNotBeNil)
		So(resp.Attrs[460796].AttrId, ShouldEqual, 4)
		So(resp.Attrs[460796].AttrSubId, ShouldEqual, 1)
		So(resp.Attrs[460796].AttrValue, ShouldEqual, 2)
	})

	Convey("Test update ranklist", t, func(c C) {
		ctx := context.TODO()
		s := &ConsumerService{
			conf: conf.Conf,
			dao:  dao.New(conf.Conf),
		}

		So(s, ShouldNotBeNil)

		// 测试小时榜
		err := s.internalLiveRankList(ctx, &databus.Message{
			Value: []byte(`	
			{
				"rank_id":1,
				"rank_type":"RANK_HOUR",
				"rank_list":{"460796":1},
				"expire_time":1549542517
			}`),
		})
		So(err, ShouldBeNil)

		req := &v1pb.FetchAttrByIDsReq{
			AttrId:    4,
			AttrSubId: 1,
			RoomIds:   []int64{460654, 460785, 460767, 460796},
		}
		// 检查插入数据
		resp, err := dao.New(conf.Conf).FetchAttrByIDs(ctx, req)
		So(err, ShouldBeNil)
		So(len(resp.Attrs), ShouldEqual, 1)
		So(resp.Attrs[460796], ShouldNotBeNil)
		So(resp.Attrs[460796].AttrId, ShouldEqual, 4)
		So(resp.Attrs[460796].AttrSubId, ShouldEqual, 1)
		So(resp.Attrs[460796].AttrValue, ShouldEqual, 1)
	})
}
