package v1

import (
	"context"
	"go-common/app/common/live/library/mengde"
	v1pb "go-common/app/service/live/zeus/api/v1"
	"go-common/app/service/live/zeus/expr"
	"go-common/app/service/live/zeus/internal/conf"
	"go-common/app/service/live/zeus/internal/service"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ZeusService struct
type ZeusService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	mengde *mengde.MengdeClient
	match  *service.MatchService
}

//NewZeusService init
func NewZeusService(c *conf.Config) (s *ZeusService) {
	mengde, err := mengde.NewMengdeClient(c.Sven.TreeID, c.Sven.Zone, c.Sven.Env, c.Sven.Build, c.Sven.Token)
	if err != nil {
		panic(err)
	}
	matcherConfig := mengde.Config().Config[c.Index.Matcher]
	match, err := service.NewMatchService(matcherConfig)
	if err != nil {
		panic(err)
	}
	s = &ZeusService{
		conf:   c,
		mengde: mengde,
		match:  match,
	}
	go s.configWatcher()
	return s
}

func (s *ZeusService) configWatcher() {
	for config := range s.mengde.ConfigNotify() {
		if err := s.match.Reload(config.Config[s.conf.Index.Matcher]); err != nil {
			log.Error("zeus reload error:%s", err.Error())
		}
	}
}

// Match implementation
// `method:"POST"`
func (s *ZeusService) Match(ctx context.Context, req *v1pb.MatchRequest) (*v1pb.MatchResponse, error) {
	env := expr.Env{
		"$build":    req.Build,
		"$version":  req.Version,
		"$buvid":    req.Buvid,
		"$platform": req.Platform,
		"$uid":      req.Uid,
	}
	isMatch, extend, err := s.match.Match(req.Group, env)
	if err != nil {
		return nil, ecode.Errorf(-1, err.Error())
	}
	resp := &v1pb.MatchResponse{
		IsMatch: isMatch,
		Extend:  extend,
	}
	return resp, nil
}
