package service

import (
	"context"

	"go-common/app/interface/main/player/model"
	accmdl "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	ugcmdl "go-common/app/service/main/ugcpay/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_playurlURI     = "/v2/playurl"
	_playurlURIV3   = "/v3/playurl"
	_h5PlayURI      = "/playurl"
	_highQaURI      = "/v2/playurlproj"
	_ugcPayOtypeArc = "archive"
	_relationPaid   = "paid"
)

// Playurl get playurl data.
func (s *Service) Playurl(c context.Context, mid int64, arg *model.PlayurlArg) (data *model.PlayurlRes, err error) {
	var (
		token, playurl string
		isUGCPayArc    bool
		viewReply      *arcmdl.ViewReply
	)
	if arg.HTML5 > 0 {
		if arg.HighQuality > 0 {
			playurl = s.highQaURL
		} else {
			playurl = s.h5PlayURL
		}
	} else {
		if viewReply, err = s.view(c, arg.Aid); err != nil {
			log.Error("Playurl s.arcClient.Arc aid(%d) error(%v)", arg.Aid, err)
			return
		}
		arc := viewReply.Arc
		if !arc.IsNormal() || !hasCid(viewReply.Pages, arg.Cid) {
			err = ecode.NothingFound
			log.Warn("Playurl verifyArchive aid(%d) can not play or no cid(%d)", arg.Aid, arg.Cid)
			return
		}
		if arc.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes || arc.AttrVal(archive.AttrBitBadgepay) == archive.AttrYes {
			err = ecode.NothingFound
			log.Warn("Playurl verifyArchive aid(%d) cid(%d) is pgc", arg.Aid, arg.Cid)
			return
		}
		if arc.AttrVal(archive.AttrBitUGCPay) == archive.AttrYes {
			if mid <= 0 {
				err = ecode.PlayURLNotLogin
				return
			} else if arc.Author.Mid != mid {
				var relation *ugcmdl.AssetRelationResp
				if relation, err = s.ugcPayClient.AssetRelation(c, &ugcmdl.AssetRelationReq{Mid: mid, Oid: arg.Aid, Otype: _ugcPayOtypeArc}); err != nil {
					log.Error("Playurl AssetRelation mid:%d aid:%d error(%+v)", mid, arg.Aid, err)
					err = ecode.PlayURLNotPay
					return
				} else if relation.State != _relationPaid {
					log.Warn("Playurl not pay aid(%d) mid(%d) state(%s)", arg.Aid, mid, relation.State)
					err = ecode.PlayURLNotPay
					return
				}
			}
			isUGCPayArc = true
		}
		if isUGCPayArc || arg.Aid%10 < s.c.Rule.PlayurlGray {
			playurl = s.playURLV3
			if mid > 0 {
				if arg.Qn == 0 {
					arg.Qn = s.c.Rule.AutoQn
				}
				if _, isVipQn := s.vipQn[arg.Qn]; isVipQn {
					if arc.Author.Mid != mid {
						var card *accmdl.CardReply
						if card, err = s.accClient.Card3(c, &accmdl.MidReq{Mid: mid}); err != nil {
							log.Error("Playurl s.accClient.Card3(%d) error(%+v)", mid, err)
							err = nil
							arg.Qn = s.c.Rule.MaxFreeQn
						} else if card.Card.Vip.Status != 1 || card.Card.Vip.Type <= 0 {
							arg.Qn = s.c.Rule.MaxFreeQn
						}
					}
				}
			} else {
				if arg.Qn > s.c.Rule.LoginQn {
					arg.Qn = s.c.Rule.LoginQn
				}
			}
		} else {
			playurl = s.playURL
			if mid > 0 {
				if arg.Qn == 0 {
					arg.Qn = s.c.Rule.AutoQn
				}
				if _, isVipQn := s.vipQn[arg.Qn]; isVipQn {
					if playurlToken, e := s.PlayURLToken(c, mid, arg.Aid, arg.Cid); e != nil {
						log.Warn("Playurl token arg(%+v) error(%v)", arg, e)
					} else if playurlToken != nil {
						token = playurlToken.Token
					}
				}
			} else {
				if arg.Qn > s.c.Rule.LoginQn {
					arg.Qn = s.c.Rule.LoginQn
				}
			}
		}
	}
	if data, err = s.dao.Playurl(c, mid, arg, playurl, token); err != nil {
		log.Error("s.dao.Playurl mid(%d) arg(%+v) token(%s) error(%+v)", mid, arg, token, err)
		// h5 high quality backup
		if arg.HTML5 > 0 && arg.HighQuality > 0 {
			err = nil
			playurl = s.h5PlayURL
			arg.HighQuality = 0
			if data, err = s.dao.Playurl(c, mid, arg, playurl, token); err != nil {
				log.Error("s.dao.Playurl h5 backup mid(%d) arg(%+v) token(%s) error(%+v)", mid, arg, token, err)
			}
		}
	}
	return
}

func hasCid(pages []*arcmdl.Page, cid int64) bool {
	for _, v := range pages {
		if cid == v.Cid {
			return true
		}
	}
	return false
}
