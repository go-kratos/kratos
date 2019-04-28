package relation

import (
	"context"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/dao/relation"
	"go-common/app/interface/main/account/model"
	acml "go-common/app/service/main/account/model"
	account "go-common/app/service/main/account/rpc/client"
	archive "go-common/app/service/main/archive/api/gorpc"
	mrl "go-common/app/service/main/relation/model"
	rlrpc "go-common/app/service/main/relation/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

var (
	_emptyFollowings = []*model.Following{}
	_emptyTagInfos   = []*model.Tag{}
	_emptyTags       = make(map[int64]string)
	_allTagsStr      = "all"
	_specialTagsStr  = "special"
	_defaultTagsStr  = "default"
	_listTagsStr     = "list"
	_emptySpList     = []int64{}
)

// Service struct of service.
type Service struct {
	// conf
	c *conf.Config
	// rpc
	relationRPC *rlrpc.Service
	accountRPC  *account.Service3
	archiveRPC  *archive.Service2
	// dao
	dao *relation.Dao
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:           c,
		relationRPC: rlrpc.New(c.RPCClient2.Relation),
		accountRPC:  account.New3(c.RPCClient2.Account),
		archiveRPC:  archive.New2(c.RPCClient2.Archive),
		dao:         relation.New(c),
	}
	return
}

// Modify modify user relation.
func (s *Service) Modify(c context.Context, mid, fid int64, act int8, src uint8, ric map[string]string) (err error) {
	if act < mrl.ActAddFollowing || act > mrl.ActDelFollower {
		err = ecode.RequestErr
		return
	}
	arg := &mrl.ArgFollowing{Mid: mid, Fid: fid, Source: src, Action: act, Infoc: ric}
	if err = s.relationRPC.ModifyRelation(c, arg); err != nil {
		log.Error("s.relationRPC.ModifyRelation(mid:%d,fid:%d,src:%d,act:%d) err(%v)", mid, fid, act, src, err)
	}
	return
}

// BatchModify batch modify user relation.
func (s *Service) BatchModify(c context.Context, mid int64, fids []int64, act int8, src uint8, ric map[string]string) (result *model.BatchModifyResult, err error) {
	if len(fids) > 50 {
		err = ecode.RequestErr
		return
	}
	for _, fid := range fids {
		if fid <= 0 || fid == mid {
			err = ecode.RequestErr
			return
		}
	}
	// luoweiling: 把非加关注的动作全部拒绝掉
	if act != mrl.ActAddFollowing {
		err = ecode.RequestErr
		return
	}
	if act < mrl.ActAddFollowing || act > mrl.ActDelFollower {
		err = ecode.RequestErr
		return
	}

	// zhangsusu: 批量关注里保持悄悄关注的状态
	whispers, err := s.relationRPC.Whispers(c, &mrl.ArgMid{
		Mid:    mid,
		RealIP: "",
	})
	if err != nil {
		log.Error("Failed to get user whispers: mid: %d: %+v", mid, err)
		return
	}
	whispersmap := make(map[int64]struct{}, len(whispers))
	for _, w := range whispers {
		whispersmap[w.Mid] = struct{}{}
	}
	filteredFids := make([]int64, 0, len(fids))
	for _, fid := range fids {
		if _, ok := whispersmap[fid]; ok {
			continue
		}
		filteredFids = append(filteredFids, fid)
	}

	raiseErr := func(in error) error {
		shouldRaise := map[int]struct{}{
			ecode.RelFollowAlreadyBlack.Code():  {},
			ecode.RelFollowReachTelLimit.Code(): {},
			ecode.RelFollowReachMaxLimit.Code(): {},
		}
		ec := ecode.Cause(in)
		if _, ok := shouldRaise[ec.Code()]; ok {
			return ec
		}
		return nil
	}

	lock := sync.Mutex{}
	result = &model.BatchModifyResult{
		FailedFids: []int64{},
	}
	wg := sync.WaitGroup{}
	for _, fid := range filteredFids {
		fid := fid
		wg.Add(1)
		go func() {
			defer wg.Done()

			arg := &mrl.ArgFollowing{Mid: mid, Fid: fid, Source: src, Action: act, Infoc: ric}
			rerr := s.relationRPC.ModifyRelation(c, arg)
			if rerr == nil {
				return
			}
			lock.Lock()
			defer lock.Unlock()
			err = raiseErr(rerr)
			log.Error("s.relationRPC.ModifyRelation(mid:%d,fid:%d,src:%d,act:%d) err(%v)", mid, fid, act, src, rerr)
			result.FailedFids = append(result.FailedFids, fid)
		}()
	}
	wg.Wait()
	return
}

// Relation get user relation.
func (s *Service) Relation(c context.Context, mid, fid int64) (f *mrl.Following, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgRelation{Mid: mid, Fid: fid, RealIP: ip}
	if f, err = s.relationRPC.Relation(c, arg); err != nil {
		log.Error("s.Relation(mid %d,fid %d) err(%v)", mid, fid, err)
	}
	return
}

// Relations get relations between users.
func (s *Service) Relations(c context.Context, mid int64, fids []int64) (f map[int64]*mrl.Following, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}
	if f, err = s.relationRPC.Relations(c, arg); err != nil {
		log.Error("s.Relations(mid %d,fids %d) err(%v)", mid, fids, err)
	}
	return
}

// Blacks get user black list.
func (s *Service) Blacks(c context.Context, mid int64, version uint64, pn, ps int64) (f []*model.Following, crc32v uint32, total int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	fr, err := s.relationRPC.Blacks(c, arg)
	if err != nil {
		log.Error("s.Blacks(mid %d) err(%v)", mid, err)
		return
	}
	total = len(fr)
	stat, err := s.relationRPC.Stat(c, arg)
	if err != nil {
		log.Error("s.Stat(mid %d) err(%v)", mid, err)
		return
	}
	total = int(stat.Black)
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= int64(len(fr)):
		fr = fr[:0]
	case end >= int64(len(fr)):
		fr = fr[start:]
	default:
		fr = fr[start:end]
	}
	if len(fr) == 0 {
		f = _emptyFollowings
		return
	}
	temp := []byte(fmt.Sprintf("%s", fr))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	var (
		mids  []int64
		infos map[int64]*acml.Info
		fi    *mrl.Following
	)
	for _, fi = range fr {
		mids = append(mids, fi.Mid)
	}
	accArg := &acml.ArgMids{Mids: mids}
	if infos, err = s.accountRPC.Infos3(c, accArg); err != nil {
		log.Error("s.accountRPC.Infos3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi = range fr {
		tmp := &model.Following{Following: fi}
		info, ok := infos[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch infos with mid: %d", fi.Mid)
			continue
		}
		tmp.Face = info.Face
		tmp.Uname = info.Name
		tmp.Sign = info.Sign
		f = append(f, tmp)
	}
	return
}

// Whispers get user Whispers.
func (s *Service) Whispers(c context.Context, mid int64, pn, ps int64, version uint64) (f []*model.Following, crc32v uint32, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	fr, err := s.relationRPC.Whispers(c, arg)
	if err != nil {
		log.Error("s.Whispers(mid %d) err(%v)", mid, err)
		return
	}
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= int64(len(fr)):
		fr = fr[:0]
	case end >= int64(len(fr)):
		fr = fr[start:]
	default:
		fr = fr[start:end]
	}
	if len(fr) == 0 {
		f = _emptyFollowings
		return
	}
	temp := []byte(fmt.Sprintf("%s", fr))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	var (
		mids  []int64
		cards map[int64]*acml.Card
		fi    *mrl.Following
	)
	for _, fi = range fr {
		mids = append(mids, fi.Mid)
	}
	accArg := &acml.ArgMids{Mids: mids}
	if cards, err = s.accountRPC.Cards3(c, accArg); err != nil {
		log.Error("s.accountRPC(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi = range fr {
		tmp := &model.Following{Following: fi}
		card, ok := cards[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", mid)
			continue
		}
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign

		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}

		// tmp.Vip = cards[fi.Mid].Vip
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		f = append(f, tmp)
	}
	return
}

// Friends get user friends list: follow eachother.
func (s *Service) Friends(c context.Context, mid int64, version uint64) (f []*model.Following, crc32v uint32, err error) {
	var (
		mids   []int64
		cards  map[int64]*acml.Card
		fi     *mrl.Following
		fo, fs []*mrl.Following
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	if fo, err = s.relationRPC.Followings(c, arg); err != nil {
		log.Error("s.Followings(mid %d) err(%v)", mid, err)
		return
	}
	for _, fi = range fo {
		if mrl.Attr(fi.Attribute) == mrl.AttrFriend {
			fs = append(fs, fi)
		}
	}
	if len(fs) == 0 {
		f = _emptyFollowings
		return
	}
	temp := []byte(fmt.Sprintf("%s", fo))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	for _, fi = range fs {
		mids = append(mids, fi.Mid)
	}
	accArg := &acml.ArgMids{Mids: mids}
	if cards, err = s.accountRPC.Cards3(c, accArg); err != nil {
		log.Error("s.accountRPC.Cards3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi = range fs {
		tmp := &model.Following{Following: fi}
		card, ok := cards[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", fi.Mid)
			continue
		}
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign
		// tmp.OfficialVerify = cards[fi.Mid].Official
		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}

		// tmp.Vip = infos[fi.Mid].Vip
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		f = append(f, tmp)
	}
	return
}

// Followers get user followings.
func (s *Service) Followers(c context.Context, vmid, mid, pn, ps int64, version uint64) (f []*model.Following, crc32v uint32, total int, err error) {
	var (
		mids  []int64
		cards map[int64]*acml.Card
		fi    *mrl.Following
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	arg := &mrl.ArgMid{Mid: vmid, RealIP: ip}
	fr, err := s.relationRPC.Followers(c, arg)
	if err != nil {
		log.Error("s.Followers(mid %d) err(%v)", vmid, err)
		return
	}
	stat, err := s.relationRPC.Stat(c, arg)
	if err != nil {
		log.Error("s.Stat(mid %d) err(%v)", vmid, err)
		return
	}
	total = int(stat.Follower)
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= int64(len(fr)):
		fr = fr[:0]
	case end >= int64(len(fr)):
		fr = fr[start:]
	default:
		fr = fr[start:end]
	}
	if len(fr) == 0 {
		f = _emptyFollowings
		return
	}
	for _, fi = range fr {
		mids = append(mids, fi.Mid)
	}
	// !self, compute !self user and up's followings' attr
	var frs map[int64]*mrl.Following
	if mid != 0 {
		argfrs := &mrl.ArgRelations{Mid: mid, Fids: mids, RealIP: ip}
		frs, err = s.relationRPC.Relations(c, argfrs)
		if err != nil {
			log.Error("s.relationRPC.Relations(c, %v) error(%v)", argfrs, err)
			return
		}
	}
	temp := []byte(fmt.Sprintf("%s", fr))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	accArg := &acml.ArgMids{Mids: mids}
	if cards, err = s.accountRPC.Cards3(c, accArg); err != nil {
		log.Error("s.accountRPC.Cards3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi = range fr {
		tmp := &model.Following{Following: fi}
		card, ok := cards[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", mid)
			continue
		}
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign
		if frst, ok := frs[fi.Mid]; ok {
			tmp.Attribute = frst.Attribute
		} else {
			tmp.Attribute = mrl.AttrNoRelation
		}
		// tmp.OfficialVerify = cards[fi.Mid].Official
		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}

		// tmp.Vip = infos[fi.Mid].Vip
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		f = append(f, tmp)
	}
	return
}

// Followings get user followings list.
func (s *Service) Followings(c context.Context, vmid, mid, pn, ps int64, version uint64, order string) (f []*model.Following, crc32v uint32, total int, err error) {
	var (
		mids  []int64
		cards map[int64]*acml.Card
		fi    *mrl.Following
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	arg := &mrl.ArgMid{Mid: vmid, RealIP: ip}
	fr, err := s.relationRPC.Followings(c, arg)
	if err != nil {
		log.Error("s.Followings(mid %d) err(%v)", vmid, err)
		return
	}
	stat, err := s.relationRPC.Stat(c, arg)
	if err != nil {
		log.Error("s.Stat(mid %d) err(%v)", vmid, err)
		return
	}
	total = int(stat.Following)
	if order == "asc" {
		sort.Sort(ByMTime(fr))
	}
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= int64(len(fr)):
		fr = fr[:0]
	case end >= int64(len(fr)):
		fr = fr[start:]
	default:
		fr = fr[start:end]
	}
	if len(fr) == 0 {
		f = _emptyFollowings
		return
	}
	for _, fi = range fr {
		mids = append(mids, fi.Mid)
	}
	// !self, compute !self user and up's followings' attr
	var frs map[int64]*mrl.Following
	if mid != vmid && mid != 0 {
		argfrs := &mrl.ArgRelations{Mid: mid, Fids: mids, RealIP: ip}
		frs, err = s.relationRPC.Relations(c, argfrs)
		if err != nil {
			log.Error("s.relationRPC.Relations(c, %v) error(%v)", argfrs, err)
			return
		}
	}
	temp := []byte(fmt.Sprintf("%s", fr))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == version {
		err = ecode.NotModified
		return
	}
	accArg := &acml.ArgMids{Mids: mids}
	if cards, err = s.accountRPC.Cards3(c, accArg); err != nil {
		log.Error("s.accountRPC.Cards3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi = range fr {
		tmp := &model.Following{Following: fi}
		card, ok := cards[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", mid)
			continue
		}
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign
		if mid != vmid {
			if frst, ok := frs[fi.Mid]; ok {
				tmp.Attribute = frst.Attribute
			} else {
				tmp.Attribute = mrl.AttrNoRelation
			}
		}
		// tmp.OfficialVerify = cards[fi.Mid].Official
		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}

		// tmp.Vip = infos[fi.Mid].Vip
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		f = append(f, tmp)
	}
	return
}

// Stat get user relation stat.
func (s *Service) Stat(c context.Context, mid int64, self bool) (st *mrl.Stat, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	if st, err = s.relationRPC.Stat(c, arg); err != nil {
		log.Error("s.Stat(mid %d) err(%v)", mid, err)
		return
	}
	if !self {
		st.Whisper = 0
		st.Black = 0
	}
	return
}

// Stats get users relation stat.
func (s *Service) Stats(c context.Context, mids []int64) (st map[int64]*mrl.Stat, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMids{Mids: mids, RealIP: ip}
	return s.relationRPC.Stats(c, arg)
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	return
}

// ByMTime implements sort.Interface for []model.Following based on the MTime field.
type ByMTime []*mrl.Following

func (mt ByMTime) Len() int           { return len(mt) }
func (mt ByMTime) Swap(i, j int)      { mt[i], mt[j] = mt[j], mt[i] }
func (mt ByMTime) Less(i, j int) bool { return mt[i].MTime < mt[j].MTime }

// Tag get tag info by tag.
func (s *Service) Tag(c context.Context, mid int64, tagid int64, pn int64, ps int64) (tagInfo []*model.Tag, err error) {
	var (
		mids  []int64
		cards map[int64]*acml.Card
		ip    = metadata.String(c, metadata.RemoteIP)
	)
	arg := &mrl.ArgTagId{Mid: mid, TagId: tagid, RealIP: ip}
	if mids, err = s.relationRPC.Tag(c, arg); err != nil {
		log.Error("s.relationRPC(%d).Arg(%v) error(%v)", mid, arg, err)
		return
	}
	var tmpMids []int64
	start, end := (pn-1)*ps, pn*ps
	switch {
	case start >= int64(len(mids)):
		tmpMids = mids[:0]
	case end >= int64(len(mids)):
		tmpMids = mids[start:]
	default:
		tmpMids = mids[start:end]
	}
	if len(tmpMids) == 0 {
		tagInfo = _emptyTagInfos
		return
	}
	accArg := &acml.ArgMids{Mids: mids}
	if cards, err = s.accountRPC.Cards3(c, accArg); err != nil {
		log.Error("s.accountRPC.Cards3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, mid = range tmpMids {
		tmp := &model.Tag{Mid: mid}
		card, ok := cards[mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", mid)
			continue
		}
		tmp.Mid = mid
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign
		// tmp.OfficialVerify = cards[mid].Official
		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}
		// tmp.Vip = infos[mid].Vip
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		tagInfo = append(tagInfo, tmp)
	}
	return
}

// Tags is.
func (s *Service) Tags(c context.Context, mid int64) (tagsCount []*mrl.TagCount, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	if tagsCount, err = s.relationRPC.Tags(c, arg); err != nil {
		log.Error("s.relationRPC(%d).Arg(%v) error(%v)", mid, arg, err)
		return
	}
	return
}

// MobileTags is.
func (s *Service) MobileTags(c context.Context, mid int64) (tagsCount map[string][]*mrl.TagCount, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgMid{Mid: mid, RealIP: ip}
	tags, err := s.relationRPC.Tags(c, arg)
	if err != nil {
		log.Error("s.relationRPC(%d).Arg(%v) error(%v)", mid, arg, err)
		return
	}
	var st *mrl.Stat
	argStat := &mrl.ArgMid{Mid: mid, RealIP: ip}
	if st, err = s.relationRPC.Stat(c, argStat); err != nil {
		log.Error("s.Stat(mid %d) err(%v)", mid, err)
		return
	}
	tagsCount = map[string][]*mrl.TagCount{
		_allTagsStr: {{
			Tagid: -1,
			Name:  "公开关注",
			Count: st.Following,
		}},
		_specialTagsStr: {{
			Tagid: -10,
			Name:  "特别关注",
			Count: 0,
		}},
		_listTagsStr:    make([]*mrl.TagCount, 0, len(tags)),
		_defaultTagsStr: make([]*mrl.TagCount, 0, 1),
	}
	for _, v := range tags {
		if v.Tagid == 0 {
			tagsCount[_defaultTagsStr] = append(tagsCount[_defaultTagsStr], v)
		} else if v.Tagid == -10 {
			tagsCount[_specialTagsStr][0].Count = v.Count
		} else {
			tagsCount[_listTagsStr] = append(tagsCount[_listTagsStr], v)
		}
	}
	return
}

// UserTag is.
func (s *Service) UserTag(c context.Context, mid int64, fid int64) (tags map[int64]string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgRelation{Mid: mid, Fid: fid, RealIP: ip}
	if tags, err = s.relationRPC.UserTag(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	if tags == nil {
		tags = _emptyTags
	}
	return
}

// CreateTag is.
func (s *Service) CreateTag(c context.Context, mid int64, tag string) (tagInfo int64, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTag{Mid: mid, Tag: tag, RealIP: ip}
	if tagInfo, err = s.relationRPC.CreateTag(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// UpdateTag is.
func (s *Service) UpdateTag(c context.Context, mid int64, tagID int64, new string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTagUpdate{Mid: mid, TagId: tagID, New: new, RealIP: ip}
	if err = s.relationRPC.UpdateTag(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// DelTag is.
func (s *Service) DelTag(c context.Context, mid int64, tagID int64) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTagDel{Mid: mid, TagId: tagID, RealIP: ip}
	if err = s.relationRPC.DelTag(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// TagsAddUsers is.
func (s *Service) TagsAddUsers(c context.Context, mid int64, tagIds string, fids string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTagsMoveUsers{Mid: mid, BeforeID: 0, AfterTagIds: tagIds, Fids: fids, RealIP: ip}
	if err = s.relationRPC.TagsAddUsers(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// TagsCopyUsers is.
func (s *Service) TagsCopyUsers(c context.Context, mid int64, tagIds string, fids string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTagsMoveUsers{Mid: mid, BeforeID: 0, AfterTagIds: tagIds, Fids: fids, RealIP: ip}
	if err = s.relationRPC.TagsCopyUsers(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// TagsMoveUsers is.
func (s *Service) TagsMoveUsers(c context.Context, mid, beforeid int64, afterTagIdsStr, fidsStr string) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &mrl.ArgTagsMoveUsers{Mid: mid, BeforeID: beforeid, AfterTagIds: afterTagIdsStr, Fids: fidsStr, RealIP: ip}
	if err = s.relationRPC.TagsMoveUsers(c, arg); err != nil {
		log.Error("s.relationRPC.Arg(%v) error(%v)", arg, err)
		return
	}
	return
}

// Prompt report and get prompt status.
func (s *Service) Prompt(c context.Context, arg *mrl.ArgPrompt) (b bool, err error) {
	return s.relationRPC.Prompt(c, arg)
}

// ClosePrompt close prompt.
func (s *Service) ClosePrompt(c context.Context, arg *mrl.ArgPrompt) (err error) {
	return s.relationRPC.ClosePrompt(c, arg)
}

// AddSpecial add fid into special.
func (s *Service) AddSpecial(c context.Context, arg *mrl.ArgFollowing) (err error) {
	return s.relationRPC.AddSpecial(c, arg)
}

// DelSpecial del fid from sepcial.
func (s *Service) DelSpecial(c context.Context, arg *mrl.ArgFollowing) (err error) {
	return s.relationRPC.DelSpecial(c, arg)
}

// Special get user special list.
func (s *Service) Special(c context.Context, mid int64) (l []int64, err error) {
	arg := &mrl.ArgMid{
		Mid: mid,
	}
	l, err = s.relationRPC.Special(c, arg)
	if len(l) == 0 {
		l = _emptySpList
	}
	return
}

// Unread check unread status, for the 'show red point' function.
func (s *Service) Unread(c context.Context, mid int64, disableAutoReset bool) (show bool, err error) {
	arg := &mrl.ArgMid{
		Mid: mid,
	}
	// if !disableAutoReset {
	// 	defer s.ResetUnread(c, mid)
	// }
	return s.relationRPC.FollowersUnread(c, arg)
}

// ResetUnread is
func (s *Service) ResetUnread(c context.Context, mid int64) (err error) {
	arg := &mrl.ArgMid{
		Mid: mid,
	}
	return s.relationRPC.ResetFollowersUnread(c, arg)
}

// UnreadCount unread count.
func (s *Service) UnreadCount(c context.Context, mid int64, disableAutoReset bool) (count int64, err error) {
	arg := &mrl.ArgMid{
		Mid: mid,
	}
	// if !disableAutoReset {
	// 	defer func() {
	// 		s.ResetUnread(c, mid)
	// 		s.ResetUnreadCount(c, mid)
	// 	}()
	// }
	return s.relationRPC.FollowersUnreadCount(c, arg)
}

// ResetUnreadCount is
func (s *Service) ResetUnreadCount(c context.Context, mid int64) (err error) {
	arg := &mrl.ArgMid{
		Mid: mid,
	}
	return s.relationRPC.ResetFollowersUnreadCount(c, arg)
}

// RecommendTagSuggestDetail is
func (s *Service) RecommendTagSuggestDetail(c context.Context, arg *model.ArgTagSuggestRecommend) (*model.TagSuggestRecommendInfo, error) {
	result, err := s.RecommendTagSuggest(c, arg)
	if err != nil {
		return nil, err
	}
	if len(result) <= 0 {
		empty := &model.TagSuggestRecommendInfo{
			TagName:  arg.TagName,
			UpList:   []*model.RecommendInfo{},
			MatchCnt: 0,
		}
		return empty, nil
	}
	detail := result[0]

	upMids := func() []int64 {
		mids := make([]int64, 0, len(detail.UpList))
		for _, up := range detail.UpList {
			mid, perr := strconv.ParseInt(up.Mid, 10, 64)
			if perr != nil {
				log.Warn("Failed to parse mid: %s: %+v", up.Mid, perr)
				continue
			}
			mids = append(mids, mid)
		}
		return mids
	}()

	rels, err := s.relationRPC.Relations(c, &mrl.ArgRelations{
		Mid:    arg.Mid,
		Fids:   upMids,
		RealIP: arg.RemoteIP,
	})
	if err != nil {
		return nil, err
	}
	for _, up := range detail.UpList {
		mid, err := strconv.ParseInt(up.Mid, 10, 64)
		if err != nil {
			log.Warn("Failed to parse mid: %s: %+v", err, up.Mid)
			continue
		}
		r, ok := rels[mid]
		if !ok {
			log.Warn("Failed to get relation between %d and %d", arg.Mid, mid)
			up.Relation = &mrl.Following{Mid: mid} // empty relation
			continue
		}
		up.Relation = r
	}

	return detail, nil
}

// RecommendTagSuggest is
func (s *Service) RecommendTagSuggest(c context.Context, arg *model.ArgTagSuggestRecommend) ([]*model.TagSuggestRecommendInfo, error) {
	resp, err := s.dao.TagSuggestRecommend(c, arg.Mid, arg.ContextID, arg.TagName, arg.Device, arg.PageSize, arg.RemoteIP)
	if err != nil {
		return nil, err
	}

	allrecs := make([]*model.RecommendContent, 0)
	for _, rec := range resp.Data {
		allrecs = append(allrecs, rec.UpList...)
	}
	allrecinfos, err := s.collectionAsRecommendUserInfo(c, allrecs, resp.TrackID, arg.RemoteIP)
	if err != nil {
		return nil, err
	}
	allrecinfomap := make(map[string]*model.RecommendInfo, len(allrecinfos))
	for _, recinfo := range allrecinfos {
		allrecinfomap[recinfo.Mid] = recinfo
	}

	getRecInfos := func(mids ...int64) []*model.RecommendInfo {
		out := make([]*model.RecommendInfo, 0, len(mids))
		for _, mid := range mids {
			smid := strconv.FormatInt(mid, 10)
			recinfo, ok := allrecinfomap[smid]
			if !ok {
				log.Warn("Failed to get user info with mid: %d", mid)
				continue
			}
			out = append(out, recinfo)
		}
		return out
	}

	result := make([]*model.TagSuggestRecommendInfo, 0, len(resp.Data))
	for _, rec := range resp.Data {
		trecinfo := &model.TagSuggestRecommendInfo{
			TagName:  rec.TagName,
			MatchCnt: rec.MatchCnt,
		}
		rinfos := getRecInfos(rec.UpIDs()...)
		trecinfo.UpList = rinfos
		if len(rinfos) != len(rec.UpIDs()) {
			log.Warn("Inconsistent user info and recommend match count: %d, %d", len(rinfos), len(rec.UpIDs()))
			trecinfo.MatchCnt -= int64(len(rec.UpIDs()) - len(rinfos))
		}
		result = append(result, trecinfo)
	}
	return result, nil
}

// RecommendFollowlistEmpty is
func (s *Service) RecommendFollowlistEmpty(c context.Context, arg *model.ArgRecommend) ([]*model.RecommendInfo, error) {
	return s.recommend(c, "followlist_empty", arg)
}

// RecommendAnswerOK is
func (s *Service) RecommendAnswerOK(c context.Context, arg *model.ArgRecommend) ([]*model.RecommendInfo, error) {
	return s.recommend(c, "answer_ok", arg)
}

func (s *Service) collectionAsRecommendUserInfo(c context.Context, recs []*model.RecommendContent, trackID string, ip string) ([]*model.RecommendInfo, error) {
	recmap := make(map[int64]*model.RecommendContent, len(recs))
	mids := make([]int64, 0, len(recs))
	for _, r := range recs {
		mids = append(mids, r.UpID)
		recmap[r.UpID] = r
	}

	cards, err := s.accountRPC.Cards3(c, &acml.ArgMids{Mids: mids})
	if err != nil {
		return nil, err
	}

	stats, err := s.relationRPC.Stats(c, &mrl.ArgMids{Mids: mids, RealIP: ip})
	if err != nil {
		return nil, err
	}

	// TODO: cache types
	types, err := s.archiveRPC.Types2(c)
	if err != nil {
		return nil, err
	}
	typeName := func(tid int16) string {
		t := types[tid]
		if t == nil {
			return ""
		}
		return t.Name
	}

	ris := make([]*model.RecommendInfo, 0, len(recs))
	for _, rec := range recs {
		if rec == nil {
			log.Warn("Invalid recommend content: %+v", rec)
			continue
		}

		card, ok := cards[rec.UpID]
		if !ok {
			log.Warn("Failed to get user card with mid: %d", rec.UpID)
			continue
		}

		stat, ok := stats[rec.UpID]
		if !ok {
			log.Warn("Failed to get stat with mid: %d", rec.UpID)
			continue
		}

		ri := &model.RecommendInfo{}
		ri.FromCard(card)
		ri.RecommendContent = *rec
		ri.TrackID = trackID
		ri.Fans = stat.Follower
		ri.TypeName = typeName(ri.Tid)
		ri.SecondTypeName = typeName(ri.SecondTid)

		// zhangsusu: 拼粉丝数作为推荐理由(后来又下线了)
		// fs := followerString(stat.Follower)
		// if fs != "" {
		// 	parts := []string{}
		// 	if ri.RecReason != "" {
		// 		parts = append(parts, ri.RecReason)
		// 	}
		// 	parts = append(parts, fs)
		// 	ri.RecReason = strings.Join(parts, "，")
		// }

		ris = append(ris, ri)
	}

	return ris, nil
}

func (s *Service) recommend(c context.Context, serviceArea string, arg *model.ArgRecommend) ([]*model.RecommendInfo, error) {
	resp, err := s.dao.Recommend(c, arg.Mid, serviceArea, arg.MainTids, arg.SubTids, arg.Device, arg.PageSize, arg.RemoteIP)
	if err != nil {
		return nil, err
	}
	return s.collectionAsRecommendUserInfo(c, resp.Data, resp.TrackID, arg.RemoteIP)
}

// func followerString(follower int64) string {
// 	if follower <= 0 {
// 		return ""
// 	}
// 	if follower < 10000 {
// 		return fmt.Sprintf("%d粉丝", follower)
// 	}
// 	return fmt.Sprintf("%.1f万粉丝", float64(follower)/float64(10000))
// }

// AchieveGet is
func (s *Service) AchieveGet(c context.Context, arg *model.ArgAchieveGet) (*mrl.AchieveGetReply, error) {
	rpcArg := &mrl.ArgAchieveGet{
		Award: arg.Award,
		Mid:   arg.Mid,
	}
	return s.relationRPC.AchieveGet(c, rpcArg)
}

// Achieve is
func (s *Service) Achieve(c context.Context, arg *model.ArgAchieve) (*model.AchieveReply, error) {
	rpcArg := &mrl.ArgAchieve{
		AwardToken: arg.AwardToken,
	}
	achieve, err := s.relationRPC.Achieve(c, rpcArg)
	if err != nil {
		return nil, err
	}

	reply := &model.AchieveReply{
		Achieve:  *achieve,
		Metadata: make(map[string]interface{}),
	}

	info, err := s.accountRPC.Info3(c, &acml.ArgMid{Mid: achieve.Mid})
	if err != nil {
		return nil, err
	}

	reply.Metadata["mid"] = info.Mid
	reply.Metadata["name"] = info.Name

	return reply, nil
}

// FollowerNotifySetting get new-follower-notification setting
func (s *Service) FollowerNotifySetting(c context.Context, mid int64) (followerNotify *mrl.FollowerNotifySetting, err error) {
	arg := &mrl.ArgMid{
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	return s.relationRPC.FollowerNotifySetting(c, arg)
}

// EnableFollowerNotify enable new-follower-notification
func (s *Service) EnableFollowerNotify(c context.Context, mid int64) (err error) {
	arg := &mrl.ArgMid{
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	return s.relationRPC.EnableFollowerNotify(c, arg)
}

// DisableFollowerNotify enable new-follower-notification
func (s *Service) DisableFollowerNotify(c context.Context, mid int64) (err error) {
	arg := &mrl.ArgMid{
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	return s.relationRPC.DisableFollowerNotify(c, arg)
}

// SameFollowings is
func (s *Service) SameFollowings(c context.Context, arg *model.ArgSameFollowing) (f []*model.Following, crc32v uint32, total int, err error) {
	sfArg := &mrl.ArgSameFollowing{Mid1: arg.VMid, Mid2: arg.Mid}
	fr, err := s.relationRPC.SameFollowings(c, sfArg)
	if err != nil {
		log.Error("s.SameFollowings(%+v) err(%v)", arg, err)
		return
	}
	total = len(fr)
	if arg.Order == "asc" {
		// 直接倒序即可
		for i := len(fr)/2 - 1; i >= 0; i-- {
			opp := len(fr) - 1 - i
			fr[i], fr[opp] = fr[opp], fr[i]
		}
	}
	start, end := (arg.PN-1)*arg.PS, arg.PN*arg.PS
	switch {
	case start >= int64(len(fr)):
		fr = fr[:0]
	case end >= int64(len(fr)):
		fr = fr[start:]
	default:
		fr = fr[start:end]
	}
	if len(fr) == 0 {
		f = _emptyFollowings
		return
	}
	mids := make([]int64, 0, len(fr))
	for _, fi := range fr {
		mids = append(mids, fi.Mid)
	}
	temp := []byte(fmt.Sprintf("%s", fr))
	crc32v = crc32.Checksum(temp, crc32.IEEETable)
	if uint64(crc32v) == arg.ReVersion {
		err = ecode.NotModified
		return
	}
	accArg := &acml.ArgMids{Mids: mids}
	cards, err := s.accountRPC.Cards3(c, accArg)
	if err != nil {
		log.Error("s.accountRPC.Cards3(mid:%v) err(%v)", accArg, err)
		return
	}
	for _, fi := range fr {
		tmp := &model.Following{Following: fi}
		card, ok := cards[fi.Mid]
		if !ok {
			log.Warn("Failed to fetch card with mid: %d", fi.Mid)
			continue
		}
		tmp.Face = card.Face
		tmp.Uname = card.Name
		tmp.Sign = card.Sign
		of := card.Official
		if of.Role == 0 {
			tmp.OfficialVerify.Type = -1
		} else {
			if of.Role <= 2 {
				tmp.OfficialVerify.Type = 0
			} else {
				tmp.OfficialVerify.Type = 1
			}
			tmp.OfficialVerify.Desc = of.Title
		}
		tmp.Vip.Type = int(card.Vip.Type)
		tmp.Vip.VipStatus = int(card.Vip.Status)
		tmp.Vip.DueDate = card.Vip.DueDate
		f = append(f, tmp)
	}
	return
}
