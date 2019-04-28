package service

import (
	"context"
	"time"

	"go-common/app/interface/main/app-player/conf"
	"go-common/app/interface/main/app-player/dao"
	accdao "go-common/app/interface/main/app-player/dao/account"
	arcdao "go-common/app/interface/main/app-player/dao/archive"
	resdao "go-common/app/interface/main/app-player/dao/resource"
	ugcpaydao "go-common/app/interface/main/app-player/dao/ugcpay"
	"go-common/app/interface/main/app-player/model"
	arcmdl "go-common/app/interface/main/app-player/model/archive"
	accrpc "go-common/app/service/main/account/api"
	ugcpaymdl "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_qn1080       = 80
	_qn480        = 32
	_relationPaid = "paid"
	_playURLV2    = "/v2/playurl"
	_playURLV3    = "/v3/playurl"
	_androidBuild = 5340000
	_iosBuild     = 8230
	_ipadHDBuild  = 12070
)

// Service is space service
type Service struct {
	c         *conf.Config
	dao       *dao.Dao
	arcDao    *arcdao.Dao
	accDao    *accdao.Dao
	ugcpayDao *ugcpaydao.Dao
	resDao    *resdao.Dao
	bnjArcs   map[int64]*arcmdl.Info
	// paster
	pasterCache map[int64]int64
}

// New new space
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:         c,
		dao:       dao.New(c),
		arcDao:    arcdao.New(c),
		accDao:    accdao.New(c),
		ugcpayDao: ugcpaydao.New(c),
		resDao:    resdao.New(c),
		// paster
		pasterCache: make(map[int64]int64),
	}
	s.loadPasterCID()
	s.loadBnjArc()
	go s.bnjTickproc()
	go s.reloadproc()
	return
}

func (s *Service) loadPasterCID() (err error) {
	var tmpPaster map[int64]int64
	if tmpPaster, err = s.resDao.PasterCID(context.Background()); err != nil {
		log.Error("%v", err)
		return
	}
	s.pasterCache = tmpPaster
	return
}

// reloadproc reload data.
func (s *Service) reloadproc() {
	for {
		time.Sleep(time.Minute * 10)
		s.loadPasterCID()
	}
}

// Playurl is
func (s *Service) Playurl(c context.Context, mid int64, params *model.Param, plat int8, buvid, fp string) (playurl *model.Playurl, err error) {
	var (
		code, isSp int
		reqPath    = _playURLV2
		aid        = params.AID
		cid        = params.CID
		build      = params.Build
		qn         = params.Qn
	)
	_, ok := s.pasterCache[cid]
	if aid > 0 && aid%10 < s.c.AidGray && !ok {
		reqPath = _playURLV3
		var (
			arc      *arcmdl.Info
			relation *ugcpaymdl.AssetRelationResp
		)
		if arc, ok = s.bnjArcs[aid]; !ok {
			if arc, err = s.arcDao.ArchiveCache(c, aid); err != nil {
				log.Error("verifyArchive %+v", err)
				return
			}
		}
		if !arc.IsNormal() || !arc.HasCid(cid) {
			err = ecode.NothingFound
			log.Warn("verifyArchive aid(%d) can not play or no cid(%d)", aid, cid)
			return
		}
		// TODO 历史老坑 纪录片之类的会请求UGC的playurl，先打日志记一记
		if arc.IsPGC() && arc.AttrVal(arcmdl.AttrBitBadgepay) == arcmdl.AttrYes {
			log.Warn("verifyArchive aid(%d) pay(%d) cid(%d) is pgc!!!!", aid, arc.AttrVal(arcmdl.AttrBitBadgepay), cid)
			err = ecode.NothingFound
			return
		}
		if arc.AttrVal(arcmdl.AttrBitUGCPay) == 1 {
			if mid <= 0 {
				aid, cid = s.validBuild(plat, build)
				if aid == 0 {
					err = ecode.PlayURLNotLogin
					return
				}
			} else if arc.Mid != mid {
				if relation, err = s.ugcpayDao.AssetRelation(c, aid, mid); err != nil {
					log.Error("verifyArchive %+v", err)
					err = ecode.PlayURLNotPay
					return
				} else if relation.State != _relationPaid {
					log.Warn("verifyArchive not pay aid(%d) mid(%d) state(%s)", aid, mid, relation.State)
					aid, cid = s.validBuild(plat, build)
					if aid == 0 {
						err = ecode.PlayURLNotPay
						return
					}
				}
			}
		}
		if mid > 0 {
			if arc.Mid == mid {
				isSp = 1
			} else {
				var card *accrpc.CardReply
				if card, err = s.accDao.Card(c, mid); err != nil {
					log.Error("verifyArchive %+v", err)
					err = nil
				} else if card.Card != nil && card.Card.Vip.IsValid() {
					isSp = 1
				}
			}
		} else if qn > _qn480 { //未登录最高清晰度 480
			qn = _qn480
		}
	}
	reqURL := s.c.Host.Playurl + reqPath
	playurl, code, err = s.dao.Playurl(c, mid, aid, cid, qn, params.Npcybs, params.Fnver, params.Fnval, params.ForceHost, isSp, params.Otype, params.MobiApp, buvid, fp, params.Session, reqURL)
	if err != nil {
		log.Error("%+v", err)
		reqURL = s.c.Host.PlayurlBk + _playURLV2
		playurl, code, err = s.dao.Playurl(c, mid, aid, cid, qn, params.Npcybs, params.Fnver, params.Fnval, params.ForceHost, isSp, params.Otype, params.MobiApp, buvid, fp, params.Session, reqURL)
		if err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if code != ecode.OK.Code() {
		log.Error("playurl aid(%d) cid(%d) code(%d)", aid, cid, code)
		err = ecode.NothingFound
		playurl = nil
	}
	return
}

func (s *Service) validBuild(plat int8, build int) (aid, cid int64) {
	if (model.IsIphone(plat) && build < _iosBuild) || (model.IsAndroid(plat) && build < _androidBuild) {
		aid = s.c.PhoneAid
		cid = s.c.PhoneCid
	} else if model.IsIPad(plat) && build <= _iosBuild {
		aid = s.c.PadAid
		cid = s.c.PadCid
	} else if model.IsIPadHD(plat) && build <= _ipadHDBuild {
		aid = s.c.PadHDAid
		cid = s.c.PadHDCid
	}
	return
}
