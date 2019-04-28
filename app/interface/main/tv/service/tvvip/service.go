package tvvip

import (
	"go-common/app/interface/main/tv/conf"
	"go-common/app/service/main/tv/api"
	"go-common/library/log"
)

// Service .
type Service struct {
	conf        *conf.Config
	tvVipClient api.TVServiceClient
}

// New .
func New(c *conf.Config) *Service {
	tvVipClient, err := api.NewClient(c.TvVipClient)
	if err != nil {
		log.Error("client.Dial err(%v)", err)
		panic(err)
	}
	srv := &Service{
		conf:        c,
		tvVipClient: tvVipClient,
	}
	return srv
}
