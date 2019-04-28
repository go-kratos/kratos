package v1

import (
	"context"

	v1appconfpb "go-common/app/interface/live/app-interface/api/http/v1"
	"go-common/app/interface/live/app-interface/conf"
	titansSdk "go-common/app/service/live/resource/sdk"
	"go-common/library/ecode"
	"go-common/library/log"
)

//AppConfService struct
type AppConfService struct {
	conf *conf.Config
}

// NewAppConfService init
func NewAppConfService(c *conf.Config) (s *AppConfService) {
	s = &AppConfService{
		conf: c,
	}
	InitTitan()
	return s
}

//GetConf 获取移动端配置
func (s *AppConfService) GetConf(ctx context.Context, req *v1appconfpb.GetConfReq) (resp *v1appconfpb.GetConfResp, err error) {

	value, ok := s.conf.AppConf[req.GetKey()]
	if !ok {
		log.Error("[AppConf] GetConf Key err: %s", req.GetKey())
		return nil, ecode.AppConfKeyErr
	}

	resp = &v1appconfpb.GetConfResp{
		Value: value,
	}
	conf, terr := titansSdk.Get(req.GetKey())
	if terr != nil {
		log.Error("[AppConf] GetConf  titansSdk.Get err: %+v", err)
	}
	if conf != "" {
		resp.Value = conf
	}
	return
}

//InitTitan 初始化kv配置
func InitTitan() {
	conf := &titansSdk.Config{
		TreeId: 61019,
		Expire: 1,
	}
	titansSdk.Init(conf)
}
