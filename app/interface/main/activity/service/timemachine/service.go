package timemachine

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/like"
	"go-common/app/interface/main/activity/dao/timemachine"
	model "go-common/app/interface/main/activity/model/timemachine"
	tagclient "go-common/app/interface/main/tag/api"
	artrpc "go-common/app/interface/openplatform/article/rpc/client"
	accclient "go-common/app/service/main/account/api"
	arcclient "go-common/app/service/main/archive/api"
	"go-common/library/log"
)

const (
	_typeArticle  = 1
	_typeArchive  = 2
	_totalViewLen = 3
)

// Service time machine service.
type Service struct {
	c                *conf.Config
	likeDao          *like.Dao
	dao              *timemachine.Dao
	article          *artrpc.Service
	tagClient        tagclient.TagRPCClient
	arcClient        arcclient.ArchiveClient
	accClient        accclient.AccountClient
	tmMidMap         map[int64]struct{}
	tagDescMap       map[int64]*model.TagDesc
	tagRegionDescMap map[int64]*model.TagRegionDesc
	regionDescMap    map[int64]*model.RegionDesc
}

// New init.
func New(c *conf.Config) *Service {
	s := &Service{
		c:       c,
		dao:     timemachine.New(c),
		likeDao: like.New(c),
		article: artrpc.New(c.RPCClient2.Article),
	}
	var err error
	if s.tagClient, err = tagclient.NewClient(c.TagClient); err != nil {
		panic(err)
	}
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	if s.accClient, err = accclient.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	go s.loadAdminMids()
	go s.tmDescproc()
	return s
}

func (s *Service) loadAdminMids() {
	for {
		time.Sleep(time.Duration(s.c.Interval.TmInternal))
		tmp := make(map[int64]struct{}, len(s.c.Rule.TmMids))
		for _, mid := range s.c.Rule.TmMids {
			tmp[mid] = struct{}{}
		}
		s.tmMidMap = tmp
	}
}

func (s *Service) tmDescproc() {
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Interval.TmInternal))
			if s.c.Timemachine.TagDescID == 0 {
				continue
			}
			if data, err := s.likeDao.SourceItem(context.Background(), s.c.Timemachine.TagDescID); err != nil {
				log.Error("tmDescproc s.likeDao.SourceItem(%d) error(%v)", s.c.Timemachine.TagDescID, err)
				continue
			} else {
				var descs struct {
					List []*struct {
						Data *model.TagDesc `json:"data"`
					} `json:"list"`
				}
				if err := json.Unmarshal(data, &descs); err != nil {
					log.Error("tmDescproc tag s.likeDao.SourceItem(%s) json error(%v)", string(data), err)
					continue
				}
				tmp := make(map[int64]*model.TagDesc, len(descs.List))
				for _, v := range descs.List {
					tmp[v.Data.TagID] = v.Data
				}
				s.tagDescMap = tmp
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Interval.TmInternal))
			if s.c.Timemachine.TagRegionDescID == 0 {
				continue
			}
			if data, err := s.likeDao.SourceItem(context.Background(), s.c.Timemachine.TagRegionDescID); err != nil {
				log.Error("tmDescproc tagRegion s.likeDao.SourceItem(%d) error(%v)", s.c.Timemachine.TagRegionDescID, err)
				continue
			} else {
				var descs struct {
					List []*struct {
						Data *model.TagRegionDesc `json:"data"`
					} `json:"list"`
				}
				if err := json.Unmarshal(data, &descs); err != nil {
					log.Error("tmDescproc tagRegion s.likeDao.SourceItem(%s) json error(%v)", string(data), err)
					continue
				}
				tmp := make(map[int64]*model.TagRegionDesc, len(descs.List))
				for _, v := range descs.List {
					tmp[v.Data.RID] = v.Data
				}
				s.tagRegionDescMap = tmp
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(s.c.Interval.TmInternal))
			if s.c.Timemachine.RegionDescID == 0 {
				continue
			}
			if data, err := s.likeDao.SourceItem(context.Background(), s.c.Timemachine.RegionDescID); err != nil {
				log.Error("tmDescproc region s.likeDao.SourceItem(%d) error(%v)", s.c.Timemachine.RegionDescID, err)
				continue
			} else {
				var descs struct {
					List []*struct {
						Data *model.RegionDesc `json:"data"`
					} `json:"list"`
				}
				if err := json.Unmarshal(data, &descs); err != nil {
					log.Error("tmDescproc region s.likeDao.SourceItem(%s) json error(%v)", string(data), err)
					continue
				}
				tmp := make(map[int64]*model.RegionDesc, len(descs.List))
				for _, v := range descs.List {
					tmp[v.Data.RID] = v.Data
				}
				s.regionDescMap = tmp
			}
		}
	}()

}
