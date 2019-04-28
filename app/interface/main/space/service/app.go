package service

import (
	"context"
	"time"

	"go-common/app/interface/main/space/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	accwar "go-common/app/service/main/account/api"
	blkmdl "go-common/app/service/main/member/model/block"
	relmdl "go-common/app/service/main/relation/model"
	upmdl "go-common/app/service/main/up/api/v1"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_devicePad   = "pad"
	_platShopH5  = 1
	_shopGoodsOn = 1
)

// AppIndex app index info.
func (s *Service) AppIndex(c context.Context, arg *model.AppIndexArg) (data *model.AppIndex, err error) {
	if env.DeployEnv == env.DeployEnvProd {
		if _, ok := s.BlacklistValue[arg.Vmid]; ok {
			err = ecode.NothingFound
			return
		}
	}
	var appInfo *model.AppAccInfo
	if appInfo, err = s.appAccInfo(c, arg.Mid, arg.Vmid, arg.Platform, arg.Device); err != nil {
		return
	}
	data = new(model.AppIndex)
	data.Info = appInfo
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		data.Tab, _ = s.appTabInfo(errCtx, arg.Mid, arg.Vmid, arg.Device, arg.Platform)
		return nil
	})
	if arg.Device == _devicePad {
		group.Go(func() error {
			data.Archive, _ = s.UpArcs(errCtx, arg.Vmid, _samplePn, arg.Ps)
			return nil
		})
	}
	group.Go(func() error {
		dyListArg := &model.DyListArg{Mid: arg.Mid, Vmid: arg.Vmid, Qn: arg.Qn, Pn: _samplePn}
		data.Dynamic, _ = s.DynamicList(errCtx, dyListArg)
		return nil
	})
	group.Wait()
	if arg.Device == _devicePad {
		if data.Archive != nil && len(data.Archive.List) > 0 {
			data.Tab.Archive = true
		}
	}
	if data.Dynamic != nil && len(data.Dynamic.List) > 0 {
		data.Tab.Dynamic = true
	}
	return
}

// AppAccInfo get app account info.
func (s *Service) appAccInfo(c context.Context, mid, vmid int64, platform, device string) (data *model.AppAccInfo, err error) {
	var (
		profile *accwar.ProfileStatReply
		ip      = metadata.String(c, metadata.RemoteIP)
	)
	data = new(model.AppAccInfo)
	if profile, err = s.accClient.ProfileWithStat3(c, &accwar.MidReq{Mid: vmid}); err != nil {
		if ecode.Cause(err) == ecode.UserNotExist {
			err = ecode.NothingFound
			return
		}
		log.Error("s.accClient.ProfileWithStat3(%d) error(%v)", vmid, err)
		err = nil
		profile = model.DefaultProfileStat
	} else if profile == nil || profile.Profile == nil {
		profile = model.DefaultProfileStat
	}
	data.FromProfile(profile)
	if data.Mid == 0 {
		data.Mid = vmid
	}
	group, errCtx := errgroup.WithContext(c)
	data.Relation = struct{}{}
	data.BeRelation = struct{}{}
	if mid > 0 {
		if mid != vmid {
			group.Go(func() error {
				if relation, err := s.relation.Relation(errCtx, &relmdl.ArgRelation{Mid: mid, Fid: vmid, RealIP: ip}); err != nil {
					log.Error("s.relation.Relation(%d,%d,%s) error %v", mid, vmid, ip, err)
				} else if relation != nil {
					data.Relation = relation
				}
				return nil
			})
			group.Go(func() error {
				if relation, err := s.relation.Relation(errCtx, &relmdl.ArgRelation{Mid: vmid, Fid: mid, RealIP: ip}); err != nil {
					log.Error("s.relation.Relation(%d,%d,%s) error %v", vmid, mid, ip, err)
				} else if relation != nil {
					data.BeRelation = relation
				}
				return nil
			})
		} else {
			data.LevelInfo = profile.LevelInfo
			if data.Silence == _silenceForbid {
				group.Go(func() error {
					if i, err := s.member.BlockInfo(errCtx, &blkmdl.RPCArgInfo{MID: mid}); err != nil {
						log.Error("s.member.BlockInfo mid(%d) error(%v)", mid, err)
						data.Block = &model.AccBlock{Status: _accBlockDefault}
					} else {
						data.Block = &model.AccBlock{
							Status: int(i.BlockStatus),
						}
						if i.BlockStatus == blkmdl.BlockStatusLimit {
							if time.Now().Unix() >= i.EndTime {
								data.Block.IsDue = _accBlockDue
							}
							if status, err := s.dao.IsAnswered(errCtx, mid, i.StartTime); err == nil {
								data.Block.IsAnswered = status
							}
						}
					}
					return nil
				})
			}
		}
	}
	//get top photo
	group.Go(func() error {
		data.TopPhoto, _ = s.dao.TopPhoto(errCtx, mid, vmid, platform, device)
		return nil
	})
	//get live status
	group.Go(func() error {
		if live, err := s.dao.Live(errCtx, vmid, ""); err != nil || live == nil {
			log.Error("s.dao.Live error(%+v) live(%+v)", err, live)
			data.Live = struct{}{}
		} else {
			data.Live = live
		}
		return nil
	})
	//get live metal
	group.Go(func() error {
		if fansBadge, err := s.dao.LiveMetal(errCtx, vmid); err != nil {
			log.Error("s.dao.LiveMetal error(%+v)", err)
		} else {
			data.FansBadge = fansBadge
		}
		return nil
	})
	//get audio card
	group.Go(func() error {
		if card, err := s.dao.AudioCard(errCtx, vmid); err != nil {
			log.Error("s.dao.AudioCard error(%+v)", err)
		} else {
			if v, ok := card[vmid]; ok && v.Type == _audioCardOn && v.Status == 1 {
				data.Audio = 1
			}
		}
		return nil
	})
	//get elec info
	group.Go(func() error {
		if elec, err := s.dao.ElecInfo(errCtx, vmid, mid); err != nil || elec == nil {
			log.Error("appAccInfo s.dao.ElecInfo vmid:%d mid:%d error(%+v) elec(%+v)", vmid, mid, err, elec)
			data.Elec = struct{}{}
		} else {
			elec.Show = true
			data.Elec = elec
		}
		return nil
	})
	//get shop info
	group.Go(func() error {
		if shop, err := s.dao.ShopLink(errCtx, vmid, _platShopH5); err != nil || shop == nil {
			log.Error("s.dao.ShopInfo error(%+v) shop(%+v)", err, shop)
			data.Shop = struct{}{}
		} else {
			data.Shop = &model.ShopInfo{ID: shop.ShopID, Name: shop.Name, URL: shop.JumpURL}
		}
		return nil
	})
	//audio card
	group.Go(func() error {
		if cert, err := s.dao.AudioUpperCert(errCtx, vmid); err != nil {
			log.Error("s.dao.AudioUpperCert error(%+v)", err)
		} else if cert != nil && cert.Cert != nil && cert.Cert.Type != -1 && cert.Cert.Desc != "" {
			if data.OfficialInfo.Type == _officialNoType {
				data.OfficialInfo.Type = cert.Cert.Type
			}
			if data.OfficialInfo.Desc != "" {
				data.OfficialInfo.Desc = data.OfficialInfo.Desc + "、" + cert.Cert.Desc
			} else {
				data.OfficialInfo.Desc = cert.Cert.Desc
			}
		}
		return nil
	})
	//group count
	group.Go(func() error {
		if fansGroup, err := s.dao.GroupsCount(errCtx, mid, vmid); err != nil {
			log.Error("s.dao.GroupsCount mid(%d) vmid(%d) error(%v)", mid, vmid, err)
		} else {
			data.FansGroup = fansGroup
		}
		return nil
	})
	group.Wait()
	return
}

// AppTabInfo get app tab info.
func (s *Service) appTabInfo(c context.Context, mid, vmid int64, device, platform string) (tab *model.AppTab, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	tab = new(model.AppTab)
	privacy := s.privacy(c, vmid)
	group, errCtx := errgroup.WithContext(c)
	// pad tab dy,arc value out this func
	if device != _devicePad {
		group.Go(func() error {
			if dyCnt, err := s.dao.DynamicCnt(errCtx, vmid); err != nil {
				log.Error("s.dao.DynamicCnt error(%+v)", err)
			} else if dyCnt > 0 {
				tab.Dynamic = true
			}
			return nil
		})
		group.Go(func() error {
			if shop, err := s.dao.ShopLink(errCtx, vmid, _platShopH5); err != nil {
				log.Error("s.dao.ShopInfo error(%+v)", err)
			} else if shop != nil && shop.ShowItemsTab == _shopGoodsOn {
				tab.Shop = true
			}
			return nil
		})
		group.Go(func() error {
			if reply, err := s.upClient.UpCount(c, &upmdl.UpCountReq{Mid: vmid}); err != nil {
				log.Error("s.arc.UpCount2 mid(%d) error(%v)", vmid, err)
			} else if reply.Count > 0 {
				tab.Archive = true
			}
			return nil
		})
		group.Go(func() error {
			if article, err := s.art.UpArtMetas(errCtx, &artmdl.ArgUpArts{Mid: vmid, Pn: 1, Ps: 10, RealIP: ip}); err != nil {
				log.Error("s.art.UpArtMetas(%d) error(%v)", vmid, err)
			} else if article != nil && len(article.Articles) > 0 {
				tab.Article = true
			}
			return nil
		})
		group.Go(func() error {
			if audioCnt, err := s.dao.AudioCnt(errCtx, vmid); err != nil {
				log.Error("s.dao.AudioCnt error(%+v)", err)
			} else if audioCnt > 0 {
				tab.Audio = true
			}
			return nil
		})
		group.Go(func() error {
			if albumCnt, err := s.dao.AlbumCount(errCtx, vmid); err == nil && albumCnt > 0 {
				tab.Album = true
			}
			return nil
		})
		if value, ok := privacy[model.PcyGame]; (ok && value == _defaultPrivacy) || mid == vmid {
			group.Go(func() error {
				if _, gameCnt, err := s.dao.AppPlayedGame(errCtx, vmid, platform, _samplePn, _samplePs); err == nil && gameCnt > 0 {
					tab.Game = true
				}
				return nil
			})
		}
	}
	if value, ok := privacy[model.PcyFavVideo]; (ok && value == _defaultPrivacy) || mid == vmid {
		group.Go(func() error {
			if fav, err := s.dao.FavFolder(errCtx, mid, vmid); err != nil {
				log.Error("s.dao.FavFolder error(%+v)", err)
			} else if len(fav) > 0 {
				for _, v := range fav {
					if v.CurCount > 0 {
						tab.Favorite = true
						break
					}
				}
			}
			return nil
		})
	}
	if value, ok := privacy[model.PcyBangumi]; (ok && value == _defaultPrivacy) || mid == vmid {
		group.Go(func() error {
			if _, cnt, err := s.dao.BangumiList(errCtx, vmid, _samplePn, _samplePs); err != nil {
				log.Error("s.dao.BangumiList mid(%d) error(%v)", vmid, err)
			} else if cnt > 0 {
				tab.Bangumi = true
			}
			return nil
		})
	}
	group.Wait()
	return
}

// AppTopPhoto get app top photo.
func (s *Service) AppTopPhoto(c context.Context, mid, vmid int64, platform, device string) (imgURL string) {
	imgURL, _ = s.dao.TopPhoto(c, mid, vmid, platform, device)
	return
}
