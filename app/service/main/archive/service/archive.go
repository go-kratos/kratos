package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/dgryski/go-farm"
)

const _qn1080 = 80

// MaxAID get max aid
func (s *Service) MaxAID(c context.Context) (id int64, err error) {
	id, err = s.arc.MaxAID(c)
	return
}

// ArchivesWithPlayer with player
func (s *Service) ArchivesWithPlayer(c context.Context, arg *archive.ArgPlayer, showPGCPlayurl bool) (ap map[int64]*archive.ArchiveWithPlayer, err error) {
	if arg == nil || len(arg.Aids) == 0 {
		err = ecode.RequestErr
		return
	}
	var (
		aids      = arg.Aids
		qn        = arg.Qn
		platform  = arg.Platform
		ip        = arg.RealIP
		fnver     = arg.Fnver
		fnval     = arg.Fnval
		session   = arg.Session
		forceHost = arg.ForceHost
		build     = arg.Build
	)
	ap = make(map[int64]*archive.ArchiveWithPlayer, len(aids))
	var as map[int64]*api.Arc
	if as, err = s.Archives3(c, aids); err != nil {
		return
	}
	for aid, arc := range as {
		ap[aid] = new(archive.ArchiveWithPlayer)
		ap[aid].Archive3 = archive.BuildArchive3(arc)
	}
	isVipQn := s.arc.VipQn(c, arg.Qn)
	if isVipQn {
		qn = _qn1080
	}
	allow := s.arc.ValidateQn(c, qn)
	if platform == "" || platform == "unknown" || ip == "" || ip == "0.0.0.0" {
		allow = false
	}
	if !allow {
		return
	}
	if !s.c.PlayerSwitch {
		return
	}
	if int64(farm.Hash32([]byte(ip)))%100 < s.c.PlayerNum {
		return
	}
	var (
		cids  []int64
		paids []int64
	)
	for _, a := range ap {
		if _, ok := s.bnjList[a.Aid]; ok && (platform == "iphone" && build == 8290) {
			log.Warn("bnj2019 aid(%d) mobi_app(%s) build(%d) in list, no playurl", a.Aid, platform, build)
			continue
		}
		if showPGCPlayurl && a.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes {
			paids = append(paids, a.Aid)
		}
		if a.FirstCid == 0 ||
			a.Access == archive.AccessMember ||
			a.AttrVal(archive.AttrBitIsPGC) == archive.AttrYes ||
			a.AttrVal(archive.AttrBitAllowBp) == archive.AttrYes ||
			a.AttrVal(archive.AttrBitBadgepay) == archive.AttrYes ||
			a.AttrVal(archive.AttrBitUGCPay) == archive.AttrYes ||
			a.AttrVal(archive.AttrBitOverseaLock) == archive.AttrYes ||
			a.AttrVal(archive.AttrBitLimitArea) == archive.AttrYes {
			log.Warn("player aid(%d) cid(%d) pgc(%d) oversea(%d) allowBp(%d) badgepay(%d) limitArea(%d) ugcpay(%d) can not with playurl",
				a.Aid, a.FirstCid, a.AttrVal(archive.AttrBitIsPGC), a.AttrVal(archive.AttrBitOverseaLock), a.AttrVal(archive.AttrBitAllowBp), a.AttrVal(archive.AttrBitBadgepay), a.AttrVal(archive.AttrBitLimitArea), a.AttrVal(archive.AttrBitUGCPay))
			continue
		}
		cids = append(cids, a.FirstCid)
	}
	if len(cids) == 0 && len(paids) == 0 {
		return
	}
	var (
		eg   errgroup.Group
		pm   map[uint32]*archive.BvcVideoItem
		pgcm map[int64]*archive.PlayerInfo
	)
	if len(cids) > 0 {
		eg.Go(func() (err error) {
			// player
			if pm, err = s.arc.PlayerInfos(context.Background(), cids, qn, platform, ip, fnver, fnval, forceHost); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	if len(paids) > 0 {
		eg.Go(func() (err error) {
			// pgc player
			if pgcm, err = s.arc.PGCPlayerInfos(context.Background(), paids, platform, ip, session, fnval, fnver); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return
		})
	}
	eg.Wait()
	for _, arc := range ap {
		if pi, ok := pm[uint32(arc.FirstCid)]; ok {
			arc.PlayerInfo = new(archive.PlayerInfo)
			arc.PlayerInfo.Cid = pi.Cid
			arc.PlayerInfo.ExpireTime = pi.ExpireTime
			arc.PlayerInfo.FileInfo = make(map[int][]*archive.PlayerFileInfo)
			for qn, files := range pi.FileInfo {
				for _, f := range files.Infos {
					arc.PlayerInfo.FileInfo[int(qn)] = append(arc.PlayerInfo.FileInfo[int(qn)], &archive.PlayerFileInfo{
						FileSize:   f.Filesize,
						TimeLength: f.Timelength,
					})
				}
			}
			arc.PlayerInfo.SupportQuality = pi.SupportQuality
			arc.PlayerInfo.SupportFormats = pi.SupportFormats
			arc.PlayerInfo.SupportDescription = pi.SupportDescription
			arc.PlayerInfo.Quality = pi.Quality
			arc.PlayerInfo.URL = pi.Url
			arc.PlayerInfo.VideoCodecid = pi.VideoCodecid
			arc.PlayerInfo.VideoProject = pi.VideoProject
			arc.PlayerInfo.Fnver = pi.Fnver
			arc.PlayerInfo.Fnval = pi.Fnval
			arc.PlayerInfo.Dash = pi.Dash
		}
		if pgci, ok := pgcm[arc.FirstCid]; ok {
			arc.PlayerInfo = pgci
			arc.Rights.Autoplay = 1
		}
	}
	return
}
