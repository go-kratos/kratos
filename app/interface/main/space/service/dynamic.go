package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/main/space/model"
	arcmdl "go-common/app/service/main/archive/api"
	coinmdl "go-common/app/service/main/coin/api"
	thumbup "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

const (
	_dyTypeCoin   = -1
	_dyTypeLike   = -2
	_dyTypeMerge  = -3
	_businessLike = "archive"
	_likeVideoCnt = 100
	_dyListCnt    = 20
	_dyDefaultQn  = 16
	_dyFoldNum    = 3
)

var dyTypeFoldMap = map[int]struct{}{_dyTypeCoin: {}, _dyTypeLike: {}, _dyTypeMerge: {}}

// DynamicList get dynamic list.
func (s Service) DynamicList(c context.Context, arg *model.DyListArg) (dyTotal *model.DyTotal, err error) {
	var (
		list, actList                    []*model.DyItem
		mergeList                        []*model.DyActItem
		dyList                           *model.DyList
		topDy                            *model.DyCard
		dyListTs, lastCoinTs, lastLikeTs int64
		topErr, dyErr                    error
		hasCoin, hasDy, hasLike, top     bool
	)
	fp := arg.Pn == 1
	repeatDyIDs := make(map[int64]int64, 1)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if topDy, topErr = s.topDynamic(errCtx, arg.Vmid, arg.Qn); topErr == nil && topDy != nil {
			top = fp
			repeatDyIDs[topDy.Desc.DynamicID] = topDy.Desc.DynamicID
		}
		return nil
	})
	group.Go(func() error {
		if dyList, dyErr = s.dao.DynamicList(errCtx, arg.Mid, arg.Vmid, arg.DyID, arg.Qn, arg.Pn); dyErr != nil {
			log.Error("s.dao.DynamicList(mid:%d,vmid:%d,dyID:%d,qn:%d,pn:%d) error(%+v)", arg.Mid, arg.Vmid, arg.DyID, arg.Qn, arg.Pn, dyErr)
		}
		return nil
	})
	group.Go(func() error {
		lastCoinTs, lastLikeTs, mergeList = s.actList(errCtx, arg.Mid, arg.Vmid)
		return nil
	})
	if e := group.Wait(); e != nil {
		log.Error("DynamicList group.Wait mid(%d) error(%v)", arg.Vmid, e)
	}
	// rm repeat data
	if dyErr == nil && dyList != nil && len(dyList.Cards) > 0 {
		for _, v := range dyList.Cards {
			if _, ok := repeatDyIDs[v.Desc.DynamicID]; ok {
				continue
			}
			item := new(model.DyResult)
			item.FromCard(v)
			list = append(list, &model.DyItem{Type: v.Desc.Type, Card: item, Ctime: v.Desc.Timestamp})
		}
		hasDy = dyList.HasMore == 1
		dyListTs = dyList.Cards[len(dyList.Cards)-1].Desc.Timestamp
	}
	if len(mergeList) > 0 {
		hasCoin, hasLike, actList = s.filterActList(c, lastCoinTs, lastLikeTs, arg.LastTime, dyListTs, mergeList, fp)
		list = append(list, actList...)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Ctime > list[j].Ctime })
	dyTotal = new(model.DyTotal)
	if top {
		topItem := new(model.DyResult)
		topItem.FromCard(topDy)
		dyTotal.List = append(dyTotal.List, &model.DyItem{Type: topDy.Desc.Type, Top: true, Card: topItem, Ctime: topDy.Desc.Timestamp})
	}
	dyTotal.HasMore = hasDy || hasCoin || hasLike
	dyTotal.List = append(dyTotal.List, list...)
	if s.c.Rule.ActFold {
		dyTotal.List = foldDyActItem(dyTotal.List)
	}
	return
}

func (s *Service) actList(c context.Context, mid, vmid int64) (lastCoinTs, lastLikeTs int64, mergeList []*model.DyActItem) {
	var (
		coinList, likeList, preList []*model.DyActItem
		coinErr, likeErr            error
		coinPcy, likePcy            bool
	)
	group, errCtx := errgroup.WithContext(c)
	privacy := s.privacy(c, vmid)
	if value, ok := privacy[model.PcyCoinVideo]; ok && value != _defaultPrivacy {
		coinPcy = true
	}
	if value, ok := privacy[model.PcyLikeVideo]; ok && value != _defaultPrivacy {
		likePcy = true
	}
	// coin video
	if mid == vmid || !coinPcy {
		group.Go(func() error {
			coinList, coinErr = s.coinVideos(errCtx, vmid, coinPcy)
			return nil
		})
	}
	// like video
	if mid == vmid || !likePcy {
		group.Go(func() error {
			likeList, likeErr = s.likeVideos(errCtx, vmid, likePcy)
			return nil
		})
	}
	group.Wait()
	if coinErr == nil {
		if l := len(coinList); l > 0 {
			preList = append(preList, coinList...)
			lastCoinTs = coinList[l-1].ActionTime
		}
	}
	if likeErr == nil {
		if l := len(likeList); l > 0 {
			preList = append(preList, likeList...)
			lastLikeTs = likeList[l-1].ActionTime
		}
	}
	if len(preList) == 0 {
		return
	}
	sort.Slice(preList, func(i, j int) bool { return preList[i].ActionTime > preList[j].ActionTime })
	if s.c.Rule.Merge {
		mergeList = mergeDyActItem(preList)
	} else {
		mergeList = preList
	}
	return
}

func (s *Service) filterActList(c context.Context, lastCoinTs, lastLikeTs, lastTime, dyListTs int64, mergeList []*model.DyActItem, fp bool) (hasCoin, hasLike bool, list []*model.DyItem) {
	var (
		actList        []*model.DyActItem
		actAids        []int64
		coinTs, likeTs int64
	)
	for _, v := range mergeList {
		if dyListTs == 0 && len(actList) >= _dyListCnt {
			lastActTs := actList[len(actList)-1].ActionTime
			penultActTs := actList[len(actList)-2].ActionTime
			y1, m1, d1 := time.Unix(lastActTs, 0).Date()
			y2, m2, d2 := time.Unix(penultActTs, 0).Date()
			if d1 != d2 || m1 != m2 || y1 != y2 {
				actList = actList[:len(actList)-1]
				break
			}
		}
		if fp {
			if dyListTs > 0 {
				if v.ActionTime >= dyListTs {
					actList = append(actList, v)
					actAids = append(actAids, v.Aid)
				}
			} else {
				actList = append(actList, v)
				actAids = append(actAids, v.Aid)
			}
		} else {
			if dyListTs > 0 {
				if v.ActionTime >= dyListTs && v.ActionTime < lastTime {
					actList = append(actList, v)
					actAids = append(actAids, v.Aid)
				}
			} else {
				if v.ActionTime < lastTime {
					actList = append(actList, v)
					actAids = append(actAids, v.Aid)
				}
			}
		}
		switch v.Type {
		case _dyTypeCoin:
			coinTs = v.ActionTime
		case _dyTypeLike:
			likeTs = v.ActionTime
		}
	}
	if coinTs > lastCoinTs {
		hasCoin = true
	}
	if likeTs > lastLikeTs {
		hasLike = true
	}
	if arcsReply, err := s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: actAids}); err != nil {
		log.Error("DynamicList s.arcClient.Arcs(%v) error(%v)", actAids, err)
	} else {
		for _, v := range actList {
			if arc, ok := arcsReply.Arcs[v.Aid]; ok && arc != nil && arc.IsNormal() {
				video := new(model.VideoItem)
				video.FromArchive(arc)
				video.ActionTime = v.ActionTime
				list = append(list, &model.DyItem{Type: v.Type, Archive: video, Ctime: v.ActionTime, Privacy: v.Privacy})
			}
		}
	}
	return
}

func mergeDyActItem(preList []*model.DyActItem) (mergeList []*model.DyActItem) {
	type privacy struct {
		Num  int
		Coin bool
		Like bool
	}
	aidNumMap := make(map[int64]*privacy, len(preList))
	aidExist := make(map[int64]struct{}, len(preList))
	for _, v := range preList {
		if _, exist := aidNumMap[v.Aid]; !exist {
			aidNumMap[v.Aid] = new(privacy)
		}
		aidNumMap[v.Aid].Num++
		switch v.Type {
		case _dyTypeCoin:
			aidNumMap[v.Aid].Coin = v.Privacy
		case _dyTypeLike:
			aidNumMap[v.Aid].Like = v.Privacy
		}
	}
	for _, v := range preList {
		num := aidNumMap[v.Aid].Num
		if num > 1 {
			if _, ok := aidExist[v.Aid]; !ok {
				v.Type = _dyTypeMerge
				v.Privacy = aidNumMap[v.Aid].Coin && aidNumMap[v.Aid].Like
				mergeList = append(mergeList, v)
			}
			aidExist[v.Aid] = struct{}{}
		} else {
			mergeList = append(mergeList, v)
		}
	}
	return
}

func foldDyActItem(list []*model.DyItem) (foldList []*model.DyItem) {
	l := len(list)
	if l == 0 {
		foldList = make([]*model.DyItem, 0)
		return
	}
	if l < _dyFoldNum {
		foldList = list
		return
	}
	var preCk bool
	for index, v := range list {
		if index == 0 {
			foldList = append(foldList, v)
			continue
		}
		last := index == l-1
		if index >= _dyFoldNum-1 {
			_, tpCheck := dyTypeFoldMap[v.Type]
			y1, m1, d1 := time.Unix(v.Ctime, 0).Date()
			_, preTpCheck := dyTypeFoldMap[list[index-1].Type]
			y2, m2, d2 := time.Unix(list[index-1].Ctime, 0).Date()
			_, check := dyTypeFoldMap[list[index-2].Type]
			y3, m3, d3 := time.Unix(list[index-2].Ctime, 0).Date()
			ck := tpCheck && preTpCheck && check && (y1 == y2 && m1 == m2 && d1 == d2) && (y1 == y3 && m1 == m3 && d1 == d3)
			// append pre item to fold if ck or preCk
			if ck || preCk {
				foldList[len(foldList)-1].Fold = append(foldList[len(foldList)-1].Fold, list[index-1])
				if last {
					foldList[len(foldList)-1].Fold = append(foldList[len(foldList)-1].Fold, v)
				}
			} else {
				foldList = append(foldList, list[index-1])
				if last {
					foldList = append(foldList, v)
				}
			}
			preCk = ck
		}
	}
	return foldList
}

func (s *Service) coinVideos(c context.Context, vmid int64, pcy bool) (list []*model.DyActItem, err error) {
	var (
		coinReply *coinmdl.ListReply
		aids      []int64
	)
	if coinReply, err = s.coinClient.List(c, &coinmdl.ListReq{Mid: vmid, Business: _businessCoin, Ts: time.Now().Unix()}); err != nil {
		log.Error("s.coinClient.List(%d) error(%v)", vmid, err)
		return
	}
	existArcs := make(map[int64]*coinmdl.ModelList, len(coinReply.List))
	for _, v := range coinReply.List {
		if len(aids) > _coinVideoLimit {
			break
		}
		if _, ok := existArcs[v.Aid]; ok {
			continue
		}
		if v.Aid > 0 {
			list = append(list, &model.DyActItem{Aid: v.Aid, Type: _dyTypeCoin, ActionTime: v.Ts, Privacy: pcy})
			existArcs[v.Aid] = v
		}
	}
	return
}

func (s *Service) likeVideos(c context.Context, mid int64, pcy bool) (list []*model.DyActItem, err error) {
	var (
		likes *thumbup.UserTotalLike
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	arg := &thumbup.ArgUserLikes{Mid: mid, Business: _businessLike, Pn: 1, Ps: _likeVideoCnt, RealIP: ip}
	if likes, err = s.thumbup.UserTotalLike(c, arg); err != nil {
		log.Error("s.thumbup.UserTotalLike(%d) error(%v)", mid, err)
		return
	}
	if likes != nil {
		for _, v := range likes.List {
			if v.MessageID > 0 {
				list = append(list, &model.DyActItem{Aid: v.MessageID, Type: _dyTypeLike, ActionTime: int64(v.Time), Privacy: pcy})
			}
		}
	}
	return
}

// topDynamic get top dynamic.
func (s *Service) topDynamic(c context.Context, mid int64, qn int) (res *model.DyCard, err error) {
	var (
		dyID int64
	)
	if dyID, err = s.dao.TopDynamic(c, mid); err != nil {
		return
	}
	if dyID == 0 {
		err = ecode.NothingFound
		return
	}
	if res, err = s.dao.Dynamic(c, mid, dyID, qn); err != nil || res == nil {
		log.Error("Dynamic s.dao.Dynamic mid(%d) dyID(%d) error(%v)", mid, dyID, err)
		err = ecode.NothingFound
	}
	return
}

// SetTopDynamic set top dynamic.
func (s *Service) SetTopDynamic(c context.Context, mid, dynamicID int64) (err error) {
	var (
		dynamic *model.DyCard
		preDyID int64
	)
	if dynamic, err = s.dao.Dynamic(c, mid, dynamicID, _dyDefaultQn); err != nil || dynamic == nil {
		log.Error("SetTopDynamic s.dao.Dynamic(%d) error(%v)", dynamicID, err)
		return
	}
	if dynamic.Desc.UID != mid {
		err = ecode.RequestErr
		return
	}
	if preDyID, err = s.dao.TopDynamic(c, mid); err != nil {
		return
	}
	if preDyID == dynamicID {
		err = ecode.NotModified
		return
	}
	if err = s.dao.AddTopDynamic(c, mid, dynamicID); err == nil {
		s.dao.AddCacheTopDynamic(c, mid, dynamicID)
	}
	return
}

// CancelTopDynamic cancel top dynamic.
func (s *Service) CancelTopDynamic(c context.Context, mid int64, now time.Time) (err error) {
	var dyID int64
	if dyID, err = s.dao.TopDynamic(c, mid); err != nil {
		return
	}
	if dyID == 0 {
		err = ecode.RequestErr
		return
	}
	if err = s.dao.DelTopDynamic(c, mid, now); err == nil {
		s.dao.AddCacheTopDynamic(c, mid, -1)
	}
	return
}
