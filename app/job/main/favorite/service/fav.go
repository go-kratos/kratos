package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/job/main/favorite/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const gap uint64 = 1 << 47
const maxSquence uint64 = 1 + gap*64000
const initSquence uint64 = 1 + gap*2000

func (s *Service) upResource(c context.Context, msg *favmdl.Message) (err error) {
	if msg.Otype == 0 {
		//为了上线的时候兼容老的service逻辑
		msg.Otype = msg.Type
	}
	if msg.Type == favmdl.TypeMusicNew {
		msg.Type = favmdl.TypeVideo
	}
	switch msg.Action {
	case favmdl.ActionAdd:
		if err = s.addFav(c, msg.Otype, msg.Mid, msg.Fid, msg.Oid, msg.FTime, msg.Type); err != nil {
			return
		}
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.upFavStat(c, msg.Otype, msg.Mid, msg.Oid, msg.FTime, true)
	case favmdl.ActionDel:
		if err = s.delFav(c, msg.Otype, msg.Mid, msg.Fid, msg.Oid, msg.FTime); err != nil {
			return
		}
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.upFavStat(c, msg.Otype, msg.Mid, msg.Oid, msg.FTime, false)
	case favmdl.ActionMultiAdd:
		s.batchUpdateFavSeqsByAadd(c, msg.Mid, msg.Fid, msg.Type, msg.Oids)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.upFavStats(c, msg.Otype, msg.Mid, msg.Oids, msg.FTime, true)
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
	case favmdl.ActionMultiDel:
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.upFavStats(c, msg.Otype, msg.Mid, msg.Oids, msg.FTime, false)
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.initFolderRelations(c, msg.Type, msg.Mid, msg.Fid)
		s.initAllRelations(c, msg.Mid, msg.Fid)
	case model.ActionCopy:
		s.batchUpdateFavSeqsByAadd(c, msg.Mid, msg.NewFid, msg.Type, msg.Oids)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.NewFid, msg.FTime)
		s.upFavStats(c, msg.Otype, msg.Mid, msg.Oids, msg.FTime, true)
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.initFolderRelations(c, msg.Type, msg.Mid, msg.NewFid)
		s.initAllRelations(c, msg.Mid, msg.NewFid)
	case model.ActionMove:
		s.batchUpdateFavSeqsByAadd(c, msg.Mid, msg.NewFid, msg.Type, msg.Oids)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.NewFid, msg.FTime)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.OldFid, msg.FTime)
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.initFolderRelations(c, msg.Type, msg.Mid, msg.NewFid)
		s.initAllRelations(c, msg.Mid, msg.NewFid)
	case model.ActionClean:
		s.cleanInvalidFavs(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.upFolderCnt(c, msg.Type, msg.Mid, msg.Fid, msg.FTime)
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
		s.initFolderRelations(c, msg.Type, msg.Mid, msg.Fid)
		s.initAllRelations(c, msg.Mid, msg.Fid)
	case favmdl.ActionFolderDel:
		if err = s.delRelationsByFid(c, msg.Type, msg.Mid, msg.Fid, msg.FTime); err != nil {
			log.Error("s.delRelationsByFid(%d,%d,%d) error(%v)", msg.Type, msg.Mid, msg.Fid, err)
		}
		if err = s.favDao.DelRelationsCache(c, msg.Mid, msg.Fid); err != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", msg.Mid, msg.Fid, err)
		}
		if err = s.favDao.DelAllRelationsCache(c, msg.Mid, msg.Fid); err != nil {
			log.Error("s.favDao.DelAllRelationsCache(%d,%d) error(%v)", msg.Mid, msg.Fid, err)
		}
		s.delRecentOidsMc(c, msg.Type, msg.Mid)
	case favmdl.ActionInitFolderRelations:
		s.initFolderRelations(c, msg.Type, msg.Mid, msg.Fid)
	case favmdl.ActionInitAllFolderRelations:
		s.initAllRelations(c, msg.Mid, msg.Fid)
	case favmdl.ActionInitRelationFids:
		var ok bool
		if ok, err = s.favDao.ExpireRelationOids(c, msg.Otype, msg.Mid); err != nil || ok {
			return
		}
		var rfmap map[int64][]int64
		if rfmap, err = s.favDao.RelationFids(c, msg.Otype, msg.Mid); err != nil {
			log.Error("favDao.FavedFids(%d,%d) error(%v)", msg.Otype, msg.Mid, err)
			return
		}
		if len(rfmap) == 0 {
			if err = s.favDao.SetUnFavedBit(c, msg.Otype, msg.Mid); err != nil {
				log.Error("s,favDao.SetUnFavedBit(type:%d,mid:%d) error(%v)", msg.Otype, msg.Mid, err)
			}
		} else {
			var oids []int64
			for oid := range rfmap {
				oids = append(oids, oid)
			}
			if err = s.favDao.SetRelationOidsCache(c, msg.Otype, msg.Mid, oids); err != nil {
				log.Error("s.favDao.SetRelationOidsCache(%d,%d,%v) error(%v)", msg.Otype, msg.Mid, oids, err)
			}
		}
	case favmdl.ActionSortFavs:
		s.sortFavs(c, msg)
	}
	return
}

func (s *Service) batchUpdateFavSeqsByAadd(c context.Context, mid, fid int64, typ int8, oids []int64) {
	log.Info("begin batchUpdateFavSeqsByAadd(%d,%d,%d,%+v)", mid, fid, typ, oids)
	favs, err := s.favDao.RelationsByOids(c, typ, mid, fid, oids)
	if err != nil {
		log.Error("s.RelationsByOids(%d,%d,%d,%v) failed!err:=%v", typ, mid, fid, oids, err)
		return
	}
	next, err := s.nextSquence(c, mid, fid)
	if err != nil {
		log.Error("s.nextSquence(%d,%d) failed!err:=%v", mid, fid, err)
		return
	}
	for _, fav := range favs {
		fav.Sequence = next
		if next >= maxSquence {
			next++
		} else {
			next += gap
		}
	}
	if err = s.batchUpdateSeqs(c, mid, fid, favs); err != nil {
		log.Error("s.batchUpdateSeqs (%d,%+v) err:=%v", mid, favs, err)
	}
	if err = s.favDao.DelAllRelationsCache(c, mid, fid); err != nil {
		err = nil
		log.Error("s.favDao.DelAllRelationsCache(%d,%d) err(%v)", mid, fid, err)
	}
}

func (s *Service) sortFavs(c context.Context, msg *favmdl.Message) {
	favs, err := s.favDao.AllRelations(c, msg.Mid, msg.Fid, 0, 1024)
	if err != nil {
		log.Error("s.favDao.AllRelations(%d,%d,%d,%d) error(%v)", msg.Mid, msg.Fid, 0, 2000, err)
		return
	}
	if len(favs) > 1000 || len(favs) <= 1 {
		log.Warn("sortFavs invalid fav(%d,%d) length:%d", msg.Mid, msg.Fid, len(favs))
		return
	}
	favsM := make(map[int64]*favmdl.Favorite)
	for i := range favs {
		favsM[favs[i].ResourceID()] = favs[i]
	}
	changed := make([]*favmdl.Favorite, 0)
	var reCount = 0
	sortFavsDesc(favs)
	if favs[len(favs)-1].Sequence == 0 {
		reSequence(favs)
		reCount++
	}
	for _, req := range msg.SortFavs {
		if req.Insert == nil {
			return
		}
		target, ok := favsM[req.Insert.ResourceID()]
		if !ok {
			return
		}
		if req.Pre == nil {
			max := favs[0].Sequence
			if max >= maxSquence {
				reCount++
				reSequence(favs)
				target.Sequence = favs[len(favs)-1].Sequence + gap
			} else if max >= maxSquence-gap {
				target.Sequence = (max + maxSquence + 1) / 2
			} else {
				target.Sequence = max + gap
			}
		} else {
			pre, ok := favsM[req.Pre.ResourceID()]
			if !ok {
				return
			}
			idx := searchIdx(favs, pre.Sequence)
			if idx == -1 {
				return
			}
			if idx < len(favs)-1 {
				next := favs[idx+1]
				if next.Oid == req.Insert.Oid && next.Type == int8(req.Insert.Typ) {
					// already sorted
					continue
				}
				if next.Sequence-pre.Sequence <= 1 {
					reCount++
					reSequence(favs)
				}
				target.Sequence = (pre.Sequence + next.Sequence) / 2
			} else {
				min := pre.Sequence
				if min <= 1 {
					// no space , need to reidx
					reCount++
					reSequence(favs)
					target.Sequence = favs[0].Sequence - gap
				} else if min <= 1+gap {
					// insert into the gap
					target.Sequence = min / 2
				} else {
					target.Sequence = min - gap
				}
			}
		}
		changed = append(changed, target)
		sortFavsDesc(favs)
	}
	if reCount > 0 {
		s.batchUpdateSeqs(c, msg.Mid, msg.Fid, favs)
	} else {
		s.batchUpdateSeqs(c, msg.Mid, msg.Fid, changed)
	}
	for i := range favs {
		favs[len(favs)-1-i].Sequence = uint64(i)
	}
	if err = s.favDao.AddAllRelationsCache(c, msg.Mid, msg.Fid, favs); err != nil {
		err = nil
		log.Error("s.favDao.AddAllRelationsCache(%d,%d,%v) err(%v)", msg.Mid, msg.Fid, favs, err)
	}
}

func (s *Service) batchUpdateSeqs(c context.Context, mid int64, fid int64, favs []*favmdl.Favorite) (err error) {
	for i := 0; i < len(favs); i += 10000 {
		end := i + 10000
		if end > len(favs) {
			end = len(favs)
		}
		_, err = s.favDao.BatchUpdateSeq(c, mid, favs[i:end])
		if err != nil {
			errStr := err.Error()
			if len(errStr) > 200 {
				errStr = errStr[:200]
			}
			log.Error("s.favDao.BatchUpdateSeq(%d,%v) err(%v)", mid, favs, errStr)
			_, err = s.favDao.BatchUpdateSeq(c, mid, favs[i:end])
			if err != nil {
				errStr := err.Error()
				if len(errStr) > 200 {
					errStr = errStr[:200]
				}
				log.Error("s.favDao.BatchUpdateSeq(%d,%v) err(%v)", mid, favs, errStr)
				return err
			}
		}
	}
	return
}

func searchIdx(favs []*favmdl.Favorite, sequence uint64) int {
	i := sort.Search(len(favs), func(i int) bool { return favs[i].Sequence <= sequence })
	if i < len(favs) && favs[i].Sequence == sequence {
		return i
	}
	return -1
}

func sortFavsDesc(favs []*favmdl.Favorite) {
	sort.Slice(favs, func(i, j int) bool {
		if favs[i].Sequence == favs[j].Sequence {
			return favs[i].MTime > favs[j].MTime
		}
		return favs[i].Sequence > favs[j].Sequence
	})
}

// 重新计算所有数据的Sequence
func reSequence(favs []*favmdl.Favorite) {
	seq := initSquence
	last := len(favs) - 1
	for i := range favs {
		favs[last-i].Sequence = seq
		seq += gap
	}
}

func (s *Service) nextSquence(c context.Context, mid, fid int64) (uint64, error) {
	max, err := s.favDao.MaxRelation(c, mid, fid)
	if err != nil {
		log.Error("s.favDao.MaxRelation(%d,%d) error(%v)", mid, fid, err)
		return 0, err
	}
	var seq uint64
	if max == nil {
		seq = initSquence
	} else if max.Sequence == 0 {
		var cnt int
		cnt, err = s.favDao.RelationCnt(c, mid, fid)
		if err != nil {
			log.Error("s.favDao.RelationCnt(%d,%d) error(%v)", mid, fid, err)
			return 0, err
		}
		if cnt <= 50000 {
			seq = initSquence + uint64(cnt+10)*gap
		} else {
			log.Error("nextSquence: can't add res over 50000")
			err = ecode.FavMaxVideoCount
			return 0, err
		}
	} else if max.Sequence+gap <= maxSquence {
		seq = max.Sequence + gap
	} else {
		seq = max.Sequence + 1
	}
	return seq, nil
}

func (s *Service) addFav(c context.Context, typ int8, mid, fid, oid, ftime int64, ftype int8) (err error) {
	v := &favmdl.Favorite{
		Type:  typ,
		Mid:   mid,
		Fid:   fid,
		Oid:   oid,
		CTime: xtime.Time(ftime),
		MTime: xtime.Time(ftime),
	}
	v.Sequence, err = s.nextSquence(c, mid, fid)
	if err != nil {
		return
	}
	rows, err := s.favDao.AddRelation(c, v)
	if err != nil {
		log.Error("s.favDao.AddRelation(%v) error(%v)", v, err)
		return
	}
	if rows < 1 {
		log.Warn("type(%d) oid(%d) already favoured", typ, oid)
		err = model.ErrFavResourceExist
		return
	}
	v.ID = rows
	err = s.cache.Do(c, func(c context.Context) {
		if err = s.favDao.SetFavedBit(c, typ, mid); err != nil {
			log.Error("s.favDao.SetFavedBit(%d, %d) error(%v)", typ, mid, err)
		}
		var ok bool
		if ftype == typ {
			if ok, err = s.favDao.ExpireRelations(c, mid, fid); err != nil {
				log.Error("s.favDao.ExpireRelations(%d,%d) error(%v)", mid, fid, err)
			}
			if ok {
				if err = s.favDao.AddRelationCache(c, v); err != nil {
					log.Error("s.favDao.AddRelationCache(%d, %d,%d) error(%v)", typ, mid, fid, err)
				}
			}
		}
		if ok, err = s.favDao.ExpireAllRelations(c, mid, fid); err != nil {
			log.Error("s.favDao.ExpireAllRelations(%d,%d) error(%v)", mid, fid, err)
		}
		if ok {
			if err = s.favDao.AddAllRelationCache(c, v); err != nil {
				log.Error("s.favDao.AddAllRelationCache(%d, %d,%d) error(%v)", typ, mid, fid, err)
			}
		}
		if ok, err = s.favDao.ExpireRelationOids(c, typ, mid); err != nil {
			log.Error("s.favDao.ExpireRelationOids(%d,%d) error(%v)", typ, mid, err)
		}
		if ok {
			if err = s.favDao.AddRelationOidCache(c, typ, mid, oid); err != nil {
				log.Error("s.favDao.AddRelationOidCache(%d,%d,%d) error(%v)", typ, mid, oid, err)
			}
		}
		if err = s.favDao.DelRelationFidsMc(c, typ, mid, oid); err != nil {
			log.Error("s.favDao.DelRelationFidsMc(%d,%d,%d) error(%v)", typ, mid, oid, err)
		}
	})
	if err != nil {
		log.Error("s.cache.Do error(%v)", err)
	}
	return
}

// DelFav delete a favorite.
func (s *Service) delFav(c context.Context, typ int8, mid, fid, oid, ftime int64) (err error) {
	rows, err := s.favDao.DelRelation(c, typ, mid, fid, oid, xtime.Time(ftime))
	if err != nil {
		log.Error("s.favDao.DelRelation(%d,%d,%d) error(%v)", typ, oid, fid, err)
		return
	}
	if rows < 1 {
		log.Warn("s.favDao.DelRelation(%d,%d,%d,%d) have no del", typ, mid, fid, oid)
		err = ecode.FavResourceAlreadyDel
		return
	}
	err = s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.DelRelationCache(c, mid, fid, oid); err1 != nil {
			log.Error("s.favDao.DelRelationCache(%d,%d,%d) error(%v)", mid, fid, oid, err1)
		}
		if err1 := s.favDao.DelAllRelationCache(c, mid, fid, oid, typ); err1 != nil {
			log.Error("s.favDao.DelAllRelationCache(%d,%d,%d) error(%v)", mid, fid, oid, err1)
		}
		if err1 := s.favDao.DelRelationFidsMc(c, typ, mid, oid); err1 != nil {
			log.Error("s.favDao.DelRelationFidsMc(%d,%d,%d) error(%v)", typ, mid, oid, err1)
		}
	})
	if err != nil {
		log.Error("s.cache.Do error(%v)", err)
	}
	return
}

func (s *Service) setAllRelationCache(c context.Context, mid, fid int64) (err error) {
	var (
		mtime    = xtime.Time(0)
		pageSize = 8000
		favss    []*favmdl.Favorite
	)
	for {
		favs, err1 := s.favDao.AllRelations(c, mid, fid, mtime, pageSize)
		if err1 != nil {
			if err = s.favDao.DelAllRelationsCache(c, mid, fid); err != nil {
				log.Error("s.favDao.DelAllRelationsCache(%d,%d) error(%v)", mid, fid, err)
			}
			return err1
		}
		if len(favs) == 0 {
			break
		}
		mtime = favs[len(favs)-1].MTime
		if mtime == favs[0].MTime {
			mtime++
		}
		favss = append(favss, favs...)
	}
	if len(favss) <= 0 {
		return
	}
	sortFavsDesc(favss)
	var needUpdateSeq bool
	if favss[len(favss)-1].Sequence == 0 {
		needUpdateSeq = true
	}
	for i := range favss {
		favss[len(favss)-1-i].Sequence = uint64(i)
	}
	if err = s.favDao.AddAllRelationsCache(c, mid, fid, favss); err != nil {
		log.Error("s.favDao.AddAllRelationsCache(%d,%d,%d,%v) error(%v)", mid, fid, favss, err)
	}
	if needUpdateSeq {
		reSequence(favss)
		s.batchUpdateSeqs(c, mid, fid, favss)
	}
	return
}

func (s *Service) setRelationCache(c context.Context, tp int8, mid, fid int64) (err error) {
	var (
		mtime    = xtime.Time(0)
		pageSize = 8000
		favss    []*favmdl.Favorite
	)
	for {
		favs, err1 := s.favDao.Relations(c, tp, mid, fid, mtime, pageSize)
		if err1 != nil {
			log.Error("s.favDao.Relations(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, mtime, pageSize, err)
			if err = s.favDao.DelRelationsCache(c, mid, fid); err != nil {
				log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err)
			}
			return err1
		}
		if len(favs) == 0 {
			break
		}
		mtime = favs[len(favs)-1].MTime
		if mtime == favs[0].MTime {
			mtime++
		}
		favss = append(favss, favs...)
	}
	if err = s.favDao.AddRelationsCache(c, tp, mid, fid, favss); err != nil {
		log.Error("s.favDao.AddRelationsCache(%d,%d,%d,%v) error(%v)", tp, mid, fid, favss, err)
	}
	return
}

func (s *Service) upFolderCnt(c context.Context, tp int8, mid, fid, ftime int64) (err error) {
	cnt, err := s.favDao.RelationCnt(c, mid, fid)
	if err != nil {
		log.Error("s.favDao.CntRelations(%d,%d,%d) error(%v)", tp, mid, fid, err)
		return
	}
	folder, err := s.folder(c, tp, mid, fid)
	if err != nil {
		return
	}
	if _, err = s.favDao.UpFolderCnt(c, mid, fid, cnt, xtime.Time(ftime)); err != nil {
		log.Error("s.favDao.UpFolderCnt(%d,%d,%d,%d) error(%v)", mid, fid, ftime, cnt, err)
		return
	}
	folder.Count = cnt
	folder.MTime = xtime.Time(ftime)
	var recent []*favmdl.Resource
	if recent, err = s.favDao.RecentRes(c, mid, fid); err != nil {
		log.Error(" s.favDao.RecentRes(%d,%d) error(%v) or folder is nil", mid, fid, err)
		err = nil
	}
	folder.RecentOids = []int64{}
	folder.RecentRes = []*favmdl.Resource{}
	if len(recent) > 0 {
		folder.RecentRes = recent
		for _, res := range recent {
			if res.Typ == int32(tp) {
				folder.RecentOids = append(folder.RecentOids, res.Oid)
			}
		}
	}
	err = s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.SetFoldersMc(c, folder); err1 != nil {
			log.Error("s.favDao.SetFolderMc(%v) error(%v)", folder, err1)
		}
		if err1 := s.favDao.DelNewCoverCache(c, folder.Mid, folder.ID); err1 != nil {
			log.Error("s.favDao.DelNewCoverCache(%v) error(%v)", folder, err1)
		}
	})
	if err != nil {
		log.Error("s.cache.Do error(%v)", err)
	}
	return
}

func (s *Service) folder(c context.Context, typ int8, mid, fid int64) (folder *favmdl.Folder, err error) {
	if folder, err = s.favDao.FolderMc(c, typ, mid, fid); err != nil {
		log.Error("s.favDao.FolderMc(%d,%d) error(%v)", mid, fid, err)
		err = nil
	}
	if folder == nil {
		if folder, err = s.favDao.Folder(c, typ, mid, fid); err != nil {
			log.Error("favDao.Folder(%d,%d) error(%v) or folder is nil", mid, fid, err)
			return
		}
	}
	if folder == nil {
		err = ecode.FavFolderNotExist
	}
	return
}

// upFavStats update upFavStat.
func (s *Service) upFavStats(c context.Context, typ int8, mid int64, oids []int64, now int64, isAdd bool) error {
	for _, oid := range oids {
		if err := s.upFavStat(c, typ, mid, oid, now, isAdd); err != nil {
			log.Error("s.upFavStat(%d,%d,%d,%d,%d) error(%v)", typ, mid, oid, now, isAdd, err)
			return err
		}
	}
	return nil
}

// upFavStat update update resource fav count.
func (s *Service) upFavStat(c context.Context, tp int8, mid, oid, now int64, isAdd bool) error {
	fids, err := s.favDao.RelationFidsByOid(c, tp, mid, oid)
	if err != nil {
		log.Error("s.favDao.RelationFidsByOid(%d,%d,%d) error(%v)", tp, mid, oid, err)
		return err
	}
	if len(fids) != 0 {
		if err1 := s.favDao.SetRelaitonFidsMc(c, tp, mid, oid, fids); err != nil {
			log.Error("s.favDao.SetRelaitonFidsMc(%d,%d,%d) error(%v)", tp, mid, oid, err1)
		}
	}
	length := len(fids)
	var incr int
	if isAdd && length == 1 {
		incr = 1
		err = s.addFavOperations(c, tp, mid, oid, now)
		if err != nil {
			return err
		}
	} else if !isAdd && length == 0 {
		incr = -1
		err = s.delFavOperations(c, tp, mid, oid, now)
		if err != nil {
			return err
		}
	}
	if incr != 0 {
		cnt, err := s.favDao.StatCnt(c, tp, oid)
		if err != nil {
			return err
		}
		if (cnt + incr) < 0 {
			return nil
		}
		rows, err := s.favDao.UpStatCnt(c, tp, oid, incr, xtime.Time(now))
		if err != nil {
			log.Error("s.favDao.UpStatCnt(%d,%d,%d) error(%v)", tp, oid, incr, err)
			return err
		}
		if rows < 1 {
			log.Warn("s.favDao.UpStatCnt(%d,%d,%d) rows(%d)", tp, oid, incr, rows)
			return nil
		}
		err = s.cache.Do(c, func(c context.Context) {
			if err = s.favDao.SetOidCountMc(c, tp, oid, int64(cnt+incr)); err != nil {
				log.Error("s.favDao.SetOidCountMc(%d,%d,%d) error(%v)", tp, oid, int64(cnt+incr), err)
			}
			if err = s.favDao.DelBatchOidsMc(c, tp, mid); err != nil {
				log.Error("s.favDao.SetOidCountMc(%d,%d) error(%v)", tp, mid, err)
			}
		})
		if err != nil {
			log.Error("s.cache.Do error(%v)", err)
		}
		s.addCoin(c, isAdd, cnt+incr, tp, oid)
		s.pubDao.PubStats(c, tp, oid, int64(cnt+incr))
		// bnj merge stat
		if err1 := s.bnjStatMerge(c, tp, oid, incr); err1 != nil {
			log.Error("s.bnjMergeStat(%d,%d,%d) error(%v)", tp, oid, incr, err1)
		}
	}
	return nil
}

func (s *Service) bnjStatMerge(c context.Context, typ int8, oid int64, incr int) (err error) {
	target := s.mergeTarget(int(typ), oid)
	if target <= 0 {
		return
	}
	cnt, err := s.favDao.StatCnt(c, typ, target)
	if err != nil {
		return
	}
	rows, err := s.favDao.UpStatCnt(c, typ, target, incr, xtime.Time(time.Now().Unix()))
	if err != nil || rows < 1 {
		log.Error("s.favDao.UpStatCnt(%d,%d,%d,%d) error(%v)", typ, target, incr, rows, err)
		return
	}
	s.pubDao.PubStats(c, typ, target, int64(cnt+incr))
	return
}

func (s *Service) addCoin(c context.Context, isAdd bool, count int, tp int8, oid int64) (err error) {
	var (
		mid            int64
		msgAdd, msgDel string
	)
	mod := count % 200
	if mod != 0 && mod != 199 {
		return
	}
	switch tp {
	case favmdl.Article:
		article, err := s.articleRPC(c, oid)
		if err != nil {
			log.Error("s.favDao.ArticleRPC error(%v)", oid, err)
			return err
		}
		meta, ok := article[oid]
		if !ok || meta == nil {
			log.Error("article martmdl.Meta(%v) error(%v)", article, err)
			return err
		}
		mid = meta.Author.Mid
		msgAdd = "专栏CV%d新增200人收藏，总收藏%d"
		msgDel = "专栏CV%d有200人取消收藏，总收藏%d"
	case favmdl.TypeVideo:
		archive, err := s.archiveRPC(c, oid)
		if err != nil {
			log.Error("s.favDao.archiveRPC error(%v)", oid, err)
			return err
		}
		mid = archive.Author.Mid
		msgAdd = "稿件AV%d新增200人收藏，总收藏%d"
		msgDel = "稿件AV%d有200人取消收藏，总收藏%d"
	default:
		log.Warn("this type(%d) need not to add coin", tp)
		return
	}
	// add money to upper
	if isAdd && mod == 0 {
		if err := s.addCoinRPC(c, mid, 1, fmt.Sprintf(msgAdd, oid, count)); err != nil {
			log.Error("s.addCoinRPC(%d,%s) error(%v)", mid, fmt.Sprintf(msgAdd, oid, count), err)
			return err
		}
	}
	if !isAdd && mod == 199 {
		if err := s.addCoinRPC(c, mid, -1, fmt.Sprintf(msgDel, oid, count)); err != nil {
			log.Error("s.addCoinRPC(%d,%s) error(%v)", mid, fmt.Sprintf(msgAdd, oid, count), err)
			return err
		}
	}
	return
}

func (s *Service) delRelationsByFid(c context.Context, typ int8, mid, fid, ftime int64) (err error) {
	var (
		offset int
		count  = s.c.Fav.MaxPageSize
	)
	typs := []int8{typ}
	if typ == 2 {
		// 收藏夹type=2是混合类型的收藏夹，需要删除多个type的稿件关系,现在只有music所以只要append 12
		typs = append(typs, 12)
	}
	for _, tp := range typs {
		for {
			var (
				rows int64
				oids []int64
			)
			if oids, err = s.favDao.OidsByFid(c, tp, mid, fid, offset, count); err != nil {
				log.Error("s.favDao.OidsByFid(%d,%d,%d,%d,%d) error(%v)", tp, mid, fid, offset, count, err)
				time.Sleep(time.Millisecond * 500) // avoid endless loop
				continue
			}
			if len(oids) == 0 {
				break
			}
			if rows, err = s.favDao.DelRelationsByOids(c, tp, mid, fid, oids, xtime.Time(ftime)); err != nil {
				log.Error("s.favDao.DelRelationsByOids(%d,%d,%d,%v) error(%v)", tp, mid, fid, oids, err)
				time.Sleep(time.Millisecond * 500) // avoid endless loop
				continue
			}
			offset += count
			if rows != int64(len(oids)) {
				log.Error("rows!=int64(len(oids)) rows:%d,len(aids):%d", rows, len(oids))
			}
			if rows > 0 {
				s.upFavStats(c, tp, mid, oids, ftime, false)
			}
			time.Sleep(time.Duration(s.c.Fav.SleepTime)) // for binlog cunsumers
		}
	}
	return
}

func (s *Service) addFavOperations(c context.Context, typ int8, mid, oid, now int64) (err error) {
	ok, err := s.favDao.ExpireRelationOids(c, typ, mid)
	if err != nil {
		log.Error("s.favDao.ExpireRelationFids(%d,%d) error(%v)", typ, mid, err)
	} else if ok {
		if err = s.favDao.AddRelationOidCache(c, typ, mid, oid); err != nil {
			log.Error("s.favDao.AddRelationOidCache(%d,%d,%d) error(%v)", typ, mid, oid, err)
		}
	}
	if typ < favmdl.TypeBangumi {
		err = nil
		return
	}
	u := &favmdl.User{
		Type:  typ,
		Oid:   oid,
		Mid:   mid,
		CTime: xtime.Time(now),
		MTime: xtime.Time(now),
	}
	rows, err := s.favDao.AddUser(c, u)
	if err != nil {
		log.Error("s.favDao.AddUser(%+v) error(%v)", u, err)
		return
	}
	if rows == 0 {
		log.Warn("s.favDao.DelUser(%+v) rows(%v)", u, rows)
	}
	return
}

func (s *Service) delFavOperations(c context.Context, typ int8, mid, oid, now int64) (err error) {
	if err = s.favDao.RemRelationOidCache(c, typ, mid, oid); err != nil {
		log.Error("s.favDao.RemRelationOidCache(%d,%d,%d) error(%v)", typ, mid, oid, err)
		err = nil
	}
	if typ < favmdl.TypeBangumi {
		return
	}
	u := &favmdl.User{
		Type:  typ,
		Oid:   oid,
		Mid:   mid,
		State: favmdl.StateIsDel,
		CTime: xtime.Time(now),
		MTime: xtime.Time(now),
	}
	rows, err := s.favDao.DelUser(c, u)
	if err != nil {
		log.Error("s.favDao.DelUser(%+v) error(%v)", u, err)
		return err
	}
	if rows == 0 {
		log.Warn("s.favDao.DelUser(%+v) rows(%v)", u, rows)
	}
	return
}

func (s *Service) delRecentOidsMc(c context.Context, typ int8, mid int64) {
	if err := s.favDao.DelRecentOidsMc(c, typ, mid); err != nil {
		log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err)
	}
	if err := s.favDao.DelRecentResMc(c, favmdl.TypeVideo, mid); err != nil {
		log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
	}
}

func (s *Service) cleanInvalidFavs(c context.Context, typ int8, mid, fid, ftime int64) (err error) {
	if typ != favmdl.TypeVideo {
		return
	}
	var (
		mtime      = xtime.Time(0)
		pageSize   = 8000
		batchCount = s.c.Fav.MaxPageSize
	)
	var oids = make(map[int64]struct{})
	var musicIds = make(map[int64]struct{})
	var batchOids []int64
	for {
		favs, err := s.favDao.AllRelations(c, mid, fid, mtime, pageSize)
		if err != nil {
			return err
		}
		if len(favs) == 0 {
			break
		}
		mtime = favs[len(favs)-1].MTime
		if mtime == favs[0].MTime {
			mtime++
		}
		for _, fav := range favs {
			if fav.Type == favmdl.TypeVideo {
				oids[fav.Oid] = struct{}{}
			} else if fav.Type == favmdl.TypeMusicNew {
				musicIds[fav.Oid] = struct{}{}
			}
		}
	}
	for oid := range oids {
		if len(batchOids) >= batchCount {
			s.cleanVideoFavs(c, mid, fid, ftime, batchOids)
			batchOids = batchOids[:0]
		}
		batchOids = append(batchOids, oid)
	}
	if len(batchOids) > 0 {
		s.cleanVideoFavs(c, mid, fid, ftime, batchOids)
	}
	batchOids = batchOids[:0]
	for oid := range musicIds {
		if len(batchOids) >= batchCount {
			s.cleanMuiscFavs(c, mid, fid, ftime, batchOids)
			batchOids = batchOids[:0]
		}
		batchOids = append(batchOids, oid)
	}
	if len(batchOids) > 0 {
		s.cleanMuiscFavs(c, mid, fid, ftime, batchOids)
	}
	batchOids = batchOids[:0]
	err = s.favDao.SetCleanedCache(c, typ, mid, fid, ftime, s.cleanCDTime)
	s.cache.Do(c, func(c context.Context) {
		if err1 := s.favDao.DelRecentOidsMc(c, typ, mid); err1 != nil {
			log.Error("s.favDao.DelRecentOidsMc(%d,%d) error(%v)", typ, mid, err1)
		}
		if err := s.favDao.DelRecentResMc(c, favmdl.TypeVideo, mid); err != nil {
			log.Error("s.favDao.DelRecentResMc(%d,%d) error(%v)", typ, mid, err)
		}
		if err1 := s.favDao.DelRelationOidsCache(c, typ, mid); err1 != nil {
			log.Error("s.favDao.DelRelationOidsCache(%d,%d) error(%v)", typ, mid, err1)
		}
		if err1 := s.favDao.DelRelationsCache(c, mid, fid); err1 != nil {
			log.Error("s.favDao.DelRelationsCache(%d,%d) error(%v)", mid, fid, err1)
		}
		if err1 := s.favDao.DelAllRelationsCache(c, mid, fid); err1 != nil {
			log.Error("s.favDao.DelAllRelationsCache(%d,%d) error(%v)", mid, fid, err1)
		}
	})
	return
}

func (s *Service) cleanMuiscFavs(c context.Context, mid, fid, ftime int64, oids []int64) (err error) {
	var delOids []int64
	musics, err := s.musicDao.MusicMap(c, oids)
	if err != nil {
		log.Error("s.ArcsRPC(%v) error(%v)", oids, err)
		return
	}
	for _, oid := range oids {
		if _, ok := musics[oid]; !ok {
			delOids = append(delOids, oid)
		}
	}
	if len(delOids) > 0 {
		var rows int64
		if rows, err = s.favDao.DelRelationsByOids(c, favmdl.TypeMusicNew, mid, fid, delOids, xtime.Time(ftime)); err != nil {
			log.Error("s.favDao.DelRelationsByOids(%d,%d,%v) error(%v)", favmdl.TypeMusicNew, mid, fid, delOids, err)
		}
		if rows > 0 {
			s.upFavStats(c, favmdl.TypeMusicNew, mid, delOids, ftime, false)
		}
	}
	return
}

func (s *Service) cleanVideoFavs(c context.Context, mid, fid, ftime int64, oids []int64) (err error) {
	var delOids []int64
	arcs, err := s.ArcsRPC(c, oids)
	if err != nil {
		log.Error("s.ArcsRPC(%v) error(%v)", oids, err)
		return
	}
	for aid, arc := range arcs {
		if arc.IsNormal() {
			continue
		}
		delOids = append(delOids, aid)
	}
	if len(delOids) > 0 {
		var rows int64
		if rows, err = s.favDao.DelRelationsByOids(c, favmdl.TypeVideo, mid, fid, delOids, xtime.Time(ftime)); err != nil {
			log.Error("s.favDao.DelRelationsByOids(%d,%d,%v) error(%v)", favmdl.TypeVideo, mid, fid, delOids, err)
		}
		if rows > 0 {
			s.upFavStats(c, favmdl.TypeVideo, mid, delOids, ftime, false)
		}
	}
	return
}

func (s *Service) initFolderRelations(c context.Context, typ int8, mid, fid int64) (err error) {
	if fid <= 0 {
		log.Warn("folderID must not be zero!%d %d", mid, fid)
		return
	}
	var ok bool
	if ok, err = s.favDao.ExpireRelations(c, mid, fid); err != nil || ok {
		return
	}
	// 顺带更新folder的count
	s.setRelationCache(c, typ, mid, fid)
	return
}

func (s *Service) initAllRelations(c context.Context, mid, fid int64) (err error) {
	if fid <= 0 {
		log.Warn("folderID must not be zero!%d %d", mid, fid)
		return
	}
	var ok bool
	if ok, err = s.favDao.ExpireAllRelations(c, mid, fid); err != nil || ok {
		return
	}
	// 顺带更新folder的count
	s.setAllRelationCache(c, mid, fid)
	return
}
