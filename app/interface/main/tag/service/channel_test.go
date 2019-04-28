package service

import (
	"context"
	"testing"

	"go-common/app/interface/main/tag/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChannel(t *testing.T) {
	var (
		tid int64 = 10176
		oid int64 = 78
		mid int64 = 14771787
		typ int8  = 1
	)
	arg := &model.ArgChannelResource{
		Tid:        tid,
		Mid:        mid,
		Plat:       1,
		LoginEvent: 2,
		RequestCNT: 10,
		DisplayID:  1,
		From:       model.ChannelFromH5,
		Type:       3,
		Build:      1,
		Name:       "unit",
		Buvid:      "1af4857eada327c03590b484d3db75c6",
		RealIP:     "127.0.0.1",
	}
	Convey("channel service", t, WithService(func(s *Service) {
		testSvc.ChannelCategory(context.Background())
		testSvc.ChanneList(context.Background(), mid, int64(typ), 0)
		testSvc.RecommandChannel(context.Background(), mid, 0)
		testSvc.DiscoveryChannel(context.Background(), mid, 0)
		testSvc.ResChannelCheckBack(context.Background(), []int64{oid}, int32(typ))
		testSvc.ChannelResources(context.Background(), arg)
	}))
}
