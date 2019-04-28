package service

import (
	"testing"

	"go-common/library/queue/databus"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceupConsumer(t *testing.T) {
	convey.Convey("upConsumer", t, func(convCtx convey.C) {
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.upConsumer()
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}

func TestServiceShardingQueueIndex(t *testing.T) {
	convey.Convey("ShardingQueueIndex", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			i := s.ShardingQueueIndex(mid)
			convCtx.Convey("Then i should not be nil.", func(convCtx convey.C) {
				convCtx.So(i, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceStart(t *testing.T) {
	convey.Convey("Start", t, func(convCtx convey.C) {
		var (
			c = make(chan *databus.Message)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.Start(c)
			convCtx.Convey("No return values", func(convCtx convey.C) {
			})
		})
	})
}
