package service

import (
	"context"
	"go-common/app/interface/bbq/app-bbq/conf"
	"go-common/app/interface/bbq/app-bbq/dao"
	"go-common/library/log"

	topic "go-common/app/service/bbq/topic/api"
	video_v1 "go-common/app/service/bbq/video/api/grpc/v1"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
)

// Service struct
type Service struct {
	c           *conf.Config
	dao         *dao.Dao
	videoClient video_v1.VideoClient
	topicClient topic.TopicClient
	httpClient  *bm.Client
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		dao:         dao.New(c),
		videoClient: newVideoClient(c.GRPCClient["video"]),
		httpClient:  bm.NewClient(c.HTTPClient.Normal),
	}
	var err error
	if s.topicClient, err = topic.NewClient(nil); err != nil {
		log.Errorw(context.Background(), "log", "get topic client fail")
		panic(err)
	}
	return s
}

// newVideoClient .
func newVideoClient(cfg *conf.GRPCConf) video_v1.VideoClient {
	cc, err := warden.NewClient(cfg.WardenConf).Dial(context.Background(), cfg.Addr)
	if err != nil {
		panic(err)
	}
	return video_v1.NewVideoClient(cc)
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func buvid(device *bm.Device) string {
	// if device.RawMobiApp == "" {
	// 	return device.Buvid3
	// }
	return device.Buvid
}
