package player

import (
	"context"
	"strings"
	"time"

	"go-common/app/interface/main/app-intl/conf"
	accdao "go-common/app/interface/main/app-intl/dao/account"
	arcdao "go-common/app/interface/main/app-intl/dao/archive"
	playerdao "go-common/app/interface/main/app-intl/dao/player"
	resdao "go-common/app/interface/main/app-intl/dao/resource"
	ugcpaydao "go-common/app/interface/main/app-intl/dao/ugcpay"
	"go-common/app/interface/main/app-intl/model"
	"go-common/app/interface/main/app-intl/model/player"
	"go-common/app/interface/main/app-intl/model/player/archive"
	accmdl "go-common/app/service/main/account/model"
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
)

var (
	vipQn = []int64{116, 112, 74}
)

// Service is space service
type Service struct {
	c         *conf.Config
	playerDao *playerdao.Dao
	arcDao    *arcdao.Dao
	accDao    *accdao.Dao
	ugcpayDao *ugcpaydao.Dao
	resDao    *resdao.Dao
	//vip qn
	vipQnMap map[int64]struct{}
	// paster
	pasterCache map[int64]int64
}

// New new space
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		playerDao:   playerdao.New(c),
		arcDao:      arcdao.New(c),
		accDao:      accdao.New(c),
		ugcpayDao:   ugcpaydao.New(c),
		resDao:      resdao.New(c),
		vipQnMap:    make(map[int64]struct{}),
		pasterCache: make(map[int64]int64),
	}

	// type cache
	for _, pn := range vipQn {
		s.vipQnMap[pn] = struct{}{}
	}
	s.loadPasterCID()
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
func (s *Service) Playurl(c context.Context, mid int64, params *player.Param, plat int8, buvid, fp string) (playurl *player.Playurl, err error) {
	var (
		code    int
		reqPath = _playURLV2
		aid     = params.AID
		cid     = params.CID
		qn      = params.Qn
	)
	_, ok := s.pasterCache[cid]
	if aid > 0 && !ok {
		reqPath = _playURLV3
		var (
			arc      *archive.Info
			relation *ugcpaymdl.AssetRelationResp
		)
		if arc, err = s.arcDao.ArchiveCache(c, aid); err != nil {
			log.Error("verifyArchive %+v", err)
			return
		}
		if !arc.IsNormal() || !arc.HasCid(cid) {
			err = ecode.NothingFound
			log.Warn("verifyArchive aid(%d) can not play or no cid(%d)", aid, cid)
			return
		}
		// TODO 历史老坑 纪录片之类的会请求UGC的playurl，先打日志记一记
		if arc.IsPGC() && arc.AttrVal(archive.AttrBitBadgepay) == archive.AttrYes {
			log.Warn("verifyArchive aid(%d) pay(%d) cid(%d) is pgc!!!!", aid, arc.AttrVal(archive.AttrBitBadgepay), cid)
			err = ecode.NothingFound
			return
		}
		if arc.AttrVal(archive.AttrBitUGCPay) == 1 {
			if mid <= 0 {
				err = ecode.PlayURLNotLogin
				return
			} else if arc.Mid != mid {
				if relation, err = s.ugcpayDao.AssetRelation(c, aid, mid); err != nil {
					log.Error("verifyArchive %+v", err)
					err = ecode.PlayURLNotPay
					return
				} else if relation.State != _relationPaid {
					log.Warn("verifyArchive not pay aid(%d) mid(%d) state(%s)", aid, mid, relation.State)
					err = ecode.PlayURLNotPay
					return
				}
			}
		}
		if mid <= 0 && qn > _qn480 {
			qn = _qn480
		}
		_, isVipQn := s.vipQnMap[qn]
		if isVipQn && arc.Mid != mid {
			var card *accmdl.Card
			if card, err = s.accDao.Card3(c, mid); err != nil {
				log.Error("verifyArchive %+v", err)
				err = nil
				qn = _qn1080
			} else if card.Vip.Status != 1 || card.Vip.Type <= 0 {
				qn = _qn1080
			}
		}
	}
	reqURL := s.c.Host.Playurl + reqPath
	playurl, code, err = s.playerDao.Playurl(c, mid, aid, cid, qn, params.Npcybs, params.Fnver, params.Fnval, params.ForceHost, params.Otype, params.MobiApp, buvid, fp, params.Session, reqURL)
	if err != nil {
		log.Error("%+v", err)
		reqURL = s.c.Host.PlayurlBk + _playURLV2
		playurl, code, err = s.playerDao.Playurl(c, mid, aid, cid, qn, params.Npcybs, params.Fnver, params.Fnval, params.ForceHost, params.Otype, params.MobiApp, buvid, fp, params.Session, reqURL)
		if err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if code != ecode.OK.Code() {
		log.Error("playurl aid(%d) cid(%d) code(%d)", aid, cid, code)
		err = ecode.NothingFound
		playurl = nil
		return
	}
	if playurl == nil {
		err = ecode.NothingFound
		return
	}
	// 版本过滤
	if plat == model.PlatAndroidI && params.Build < 2002000 {
		qualitys := make([]int64, 0, len(playurl.AcceptQuality))
		descs := make([]string, 0, len(playurl.AcceptDescription))
		formats := make([]string, 0, len(playurl.AcceptQuality))
		acceptFormats := strings.Split(playurl.AcceptFormat, ",")
		for index, quality := range playurl.AcceptQuality {
			if _, ok := s.vipQnMap[quality]; !ok {
				qualitys = append(qualitys, quality)
				if index < len(playurl.AcceptDescription) {
					descs = append(descs, playurl.AcceptDescription[index])
				}
				if index < len(acceptFormats) {
					formats = append(formats, acceptFormats[index])
				}
			}
		}
		playurl.AcceptQuality = qualitys
		playurl.AcceptDescription = descs
		playurl.AcceptFormat = strings.Join(formats, ",")
	}
	return
}
