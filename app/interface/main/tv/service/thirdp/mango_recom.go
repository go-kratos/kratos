package thirdp

import (
	"context"
	"time"

	"go-common/app/interface/main/tv/model"
	tpMdl "go-common/app/interface/main/tv/model/thirdp"
	arcwar "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_rtypePGC = 1
	_rtypeUGC = 2
)

func (s *Service) mangoR() (err error) {
	var (
		ctx     = context.Background()
		rids    []int64
		recoms  []*tpMdl.MangoRecom
		params  []*tpMdl.MangoParams
		catInfo *arcwar.Tp
	)
	if rids, err = s.dao.MangoOrder(ctx); err != nil { // pick mango recoms' order
		log.Error("mango MangoOrder Error %v", err)
		return
	}
	if len(rids) == 0 {
		log.Error("mango MangoOrder Empty")
		return
	}
	if recoms, err = s.dao.MangoRecom(ctx, rids); err != nil { // pick mango recom data
		log.Error("mango MangoRecom Rids [%v], Err %v", rids, err)
		return
	}
	for _, recom := range recoms {
		if recom.Rtype == _rtypePGC {
			var sn *model.SeasonCMS
			if sn, err = s.cmsDao.LoadSnCMS(context.Background(), recom.RID); err != nil {
				return err
			}
			param := recom.ToParam()
			param.Category = tpMdl.PgcCat(recom.Category)
			param.Role = sn.Role
			param.PlayTime = sn.Playtime.Time().Format("2006-01-02")
			params = append(params, param)
		} else if recom.Rtype == _rtypeUGC {
			var arc *model.ArcCMS
			if arc, err = s.cmsDao.LoadArcMeta(context.Background(), recom.RID); err != nil {
				return err
			}
			param := recom.ToParam()
			if catInfo, err = s.arcDao.TypeInfo(int32(recom.Category)); err != nil { // pick ugc category name
				log.Warn("MangoRecom Recom RID %d, Cat %d", recom.RID, recom.Category)
			} else {
				param.Category = catInfo.Name
			}
			param.PlayTime = arc.Pubtime.Time().Format("2006-01-02")
			params = append(params, param)
		} else {
			return ecode.TvDangbeiWrongType
		}
	}
	if len(params) > 0 {
		s.mangoRecom = params
	}
	return
}

func (s *Service) mangorproc() {
	for {
		time.Sleep(time.Duration(s.conf.Cfg.PageReload))
		if err := s.mangoR(); err != nil {
			log.Error("mango Error %v", err)
		}
	}
}

// MangoRecom returns the mango recom data
func (s *Service) MangoRecom() (data []*tpMdl.MangoParams) {
	if len(s.mangoRecom) == 0 {
		data = make([]*tpMdl.MangoParams, 0)
		return
	}
	data = s.mangoRecom
	return
}
