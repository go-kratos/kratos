package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go-common/app/interface/main/reply/dao/reply"
	model "go-common/app/interface/main/reply/model/reply"
	xmodel "go-common/app/interface/main/reply/model/xreply"
	accmdl "go-common/app/service/main/account/api"
	assmdl "go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
	"sort"

	"go-common/library/sync/errgroup.v2"
)

const (
	defaultChildrenSize = 5
)

// NewCursorByReplyID NewCursorByReplyID
func (s *Service) NewCursorByReplyID(ctx context.Context, oid int64,
	otyp int8, replyID int64, size int, cmp model.Comp) (*model.Cursor, error) {

	rs, err := s.GetReplyByIDs(ctx, oid, otyp, []int64{replyID})
	if err != nil {
		return nil, err
	}
	if r, ok := rs[replyID]; ok {
		return model.NewCursor(int64(r.Floor), 0, size, cmp)
	}
	return nil, ecode.ReplyNotExist
}

// NewSubCursorByReplyID NewSubCursorByReplyID
func (s *Service) NewSubCursorByReplyID(ctx context.Context, oid int64, otyp int8, replyID int64, size int, cmp model.Comp) (rootID int64, cursor *model.Cursor, err error) {
	rs, err := s.GetReplyByIDs(ctx, oid, otyp, []int64{replyID})
	if err != nil {
		return 0, nil, err
	}
	if r, ok := rs[replyID]; ok {
		if r.IsRoot() {
			rootID = r.RpID
			cursor, err = model.NewCursor(0, 1, size, cmp)
			return
		}
		// 不足一页面时，展示够一页
		floor := r.Floor
		if floor < size {
			floor = size
		}
		rootID = r.Root
		cursor, err = model.NewCursor(int64(floor), 0, size, cmp)
		return
	}
	return 0, nil, ecode.ReplyNotExist
}

// GetRootReplyListHeader GetRootReplyListHeader
func (s *Service) GetRootReplyListHeader(ctx context.Context, sub *model.Subject, params *model.CursorParams) (*model.RootReplyListHeader, error) {
	var hotIDs []int64
	res, err := s.replyHotFeed(ctx, params.Mid, sub.Oid, int(sub.Type), 1, params.HotSize+2)
	if err == nil && res != nil && len(res.RpIDs) > 0 {
		log.Info("reply-feed(test): reply abtest mid(%d) oid(%d) type(%d) test name(%s) rpIDs(%v)", params.Mid, sub.Oid, sub.Type, res.Name, res.RpIDs)
		hotIDs = res.RpIDs
	} else {
		if err != nil {
			log.Error("reply-feed error(%v)", err)
			err = nil
		} else {
			log.Info("reply-feed(origin): reply abtest mid(%d) oid(%d) type(%d) test name(%s) rpIDs(%v)", params.Mid, sub.Oid, sub.Type, res.Name, res.RpIDs)
		}
		if hotIDs, err = s.GetRootReplyIDs(ctx, sub.Oid, sub.Type, model.SortByLike, 0, int64(params.HotSize+2)); err != nil {
			log.Error("%v", err)
			return nil, err
		}
	}
	var parentIDs []int64
	parentIDs = append(parentIDs, hotIDs...)

	var adminTopReply, upperTopReply *model.Reply

	if sub.AttrVal(model.SubAttrAdminTop) == model.AttrYes {
		adminTopReply, err = s.GetTopReply(ctx, params.Oid, params.OTyp, model.SubAttrAdminTop)
		if err != nil {
			return nil, err
		}
		if adminTopReply != nil {
			parentIDs = append(parentIDs, adminTopReply.RpID)
		}
	}
	if sub.AttrVal(model.SubAttrUpperTop) == model.AttrYes {
		upperTopReply, err = s.GetTopReply(ctx, params.Oid, params.OTyp, model.SubAttrUpperTop)
		if err != nil {
			return nil, err
		}
		if upperTopReply != nil {
			if !upperTopReply.IsNormal() && sub.Mid != params.Mid {
				upperTopReply = nil
			} else {
				parentIDs = append(parentIDs, upperTopReply.RpID)
			}
		}
	}

	parentChildrenIDRelation, err := s.ParentChildrenReplyIDRelation(ctx, sub, parentIDs)
	if err != nil {
		return nil, err
	}

	idReplyMap, err := s.IDReplyMap(ctx, sub, parentChildrenIDRelation)
	if err != nil {
		return nil, err
	}

	rootIDReplyMap := assemble(idReplyMap, parentChildrenIDRelation)
	if adminTopReply != nil {
		if r, ok := rootIDReplyMap[adminTopReply.RpID]; ok {
			adminTopReply = r
		}
		// For historic reasons, TopReply and HotReply may be overlapped
		hotIDs = Remove(hotIDs, adminTopReply.RpID)
	}
	if upperTopReply != nil {
		if r, ok := rootIDReplyMap[upperTopReply.RpID]; ok {
			upperTopReply = r
		}
		// For historic reasons, TopReply and HotReply may be overlapped
		hotIDs = Remove(hotIDs, upperTopReply.RpID)
	}
	return &model.RootReplyListHeader{
		TopAdmin: adminTopReply,
		TopUpper: upperTopReply,
		Hots:     filterHot(Fetch(rootIDReplyMap, hotIDs), params.HotSize),
	}, nil
}

func filterHot(rs []*model.Reply, maxSize int) (hots []*model.Reply) {
	for _, r := range rs {
		if r.Like >= 3 {
			hots = append(hots, r)
		}
	}
	if hots == nil {
		hots = _emptyReplies
	} else if len(hots) > maxSize {
		hots = hots[:maxSize]
	}
	return hots
}

func needHeader(cursor *model.Cursor, rootLen int) bool {
	return cursor.Latest() ||
		(cursor.Increase() && rootLen < int(cursor.Len()))
}

// NeedInsertPendingReply NeedInsertPendingReply
func NeedInsertPendingReply(params *model.CursorParams, sub *model.Subject) bool {
	return params.Mid > 0 &&
		params.Sort == model.SortByFloor &&
		sub.AttrVal(model.SubAttrAudit) == model.AttrYes
}

func collect(r *model.Reply, allIDs []int64, allMIDs []int64,
	allReply []*model.Reply) ([]int64, []int64, []*model.Reply) {

	if r == nil {
		return nil, nil, nil
	}
	allIDs = append(allIDs, r.RpID)
	allReply = append(allReply, r)
	allMIDs = append(allMIDs, r.Mid)
	if r.Content != nil {
		for _, mid := range r.Content.Ats {
			allMIDs = append(allMIDs, mid)
		}
	}
	return allIDs, allMIDs, allReply
}

// IDReplyMap IDReplyMap
func (s *Service) IDReplyMap(ctx context.Context, sub *model.Subject,
	parentChildrenIDRelation map[int64][]int64) (map[int64]*model.Reply, error) {

	var allIDs []int64
	for parentID, childrenIDs := range parentChildrenIDRelation {
		allIDs = append(allIDs, childrenIDs...)
		allIDs = append(allIDs, parentID)
	}
	// WARNING: GetReplyByIDs should not contains subReplies, but currently there
	// exists a bug, which makes `idReplyMap` may contains sub_reply
	idReplyMap, err := s.GetReplyByIDs(ctx, sub.Oid, sub.Type, allIDs)
	if err != nil {
		return nil, err
	}
	// temporary solution :(, remove all children replies
	for _, reply := range idReplyMap {
		if reply.Replies != nil {
			reply.Replies = reply.Replies[:0]
		}
	}
	return idReplyMap, nil
}

// RootReplyListByCursor RootReplyListByCursor
func (s *Service) RootReplyListByCursor(ctx context.Context, sub *model.Subject, params *model.CursorParams) ([]*model.Reply, error) {
	var parentIDs []int64
	if params.Cursor.Latest() {
		// 忽略错误，这个请求只为了增加统计数据
		s.replyFeed(ctx, params.Mid, 1, 20)
	} else {
		s.replyFeed(ctx, params.Mid, 2, 20)
	}
	rootIDs, err := s.GetRootReplyIDsByCursor(ctx, sub, params.Sort, params.Cursor)
	if err != nil {
		return nil, err
	}
	// 老版本折叠评论的逻辑
	if params.ShowFolded && sub.HasFolded() {
		foldedrpIDs, _ := s.foldedRepliesCursor(ctx, sub, 0, params.Cursor)
		if len(foldedrpIDs) > 0 {
			rootIDs = append(rootIDs, foldedrpIDs...)
			sort.Slice(rootIDs, func(x, y int) bool { return rootIDs[x] > rootIDs[y] })
			length := len(rootIDs)
			if length > params.Cursor.Len() {
				if params.Cursor.Increase() {
					// 对于根评论列表，往楼层大的方向翻页是向上翻，需要从后往前截断
					rootIDs = rootIDs[length-params.Cursor.Len():]
				} else {
					rootIDs = rootIDs[:params.Cursor.Len()]
				}
			}
		}
	}
	parentIDs = append(parentIDs, rootIDs...)
	parentChildrenIDRelation, err := s.ParentChildrenReplyIDRelation(ctx, sub, parentIDs)
	if err != nil {
		return nil, err
	}

	idReplyMap, err := s.IDReplyMap(ctx, sub, parentChildrenIDRelation)
	if err != nil {
		return nil, err
	}
	if NeedInsertPendingReply(params, sub) {
		// WARNING: here we assume that pending replies have no children
		// otherwise, we need to change logic here
		pendingIDReplyMap, err := s.GetPendingReply(ctx, params.Mid, sub.Oid, sub.Type)
		if err != nil {
			return nil, err
		}
		for id, r := range pendingIDReplyMap {
			if r.IsRoot() && !r.IsTop() {
				// insert pending reply into root reply list
				if params.Cursor.Latest() &&
					((len(rootIDs) > 0 && id > rootIDs[0]) || len(rootIDs) == 0) {
					// when fetch latest reply list, and root reply list's length < the default size
					// and the pending reply ID > the max rootID
					// then just append the pending reply
					rootIDs = append([]int64{id}, rootIDs...)
					if len(rootIDs) > int(params.Cursor.Len()) {
						rootIDs = rootIDs[:params.Cursor.Len()]
					}
				} else {
					// otherwise, we need an algorithm to insert pending replyID into
					// rootIDs
					rootIDs = InsertInto(rootIDs, id, int(params.Cursor.Len()), model.OrderDESC)
				}
				parentChildrenIDRelation[id] = []int64{}
			} else if _, ok := idReplyMap[r.Root]; ok {
				// insert pending reply into sub reply list
				parentChildrenIDRelation[r.Root] = InsertInto(parentChildrenIDRelation[r.Root], id, defaultChildrenSize, model.OrderASC)
			} else {
				continue
			}
			sub.ACount++
			idReplyMap[id] = r
		}
	}
	return Fetch(assemble(idReplyMap, parentChildrenIDRelation), rootIDs), nil
}

// Remove Remove
func Remove(arr []int64, k int64) []int64 {
	b := arr[:0]
	for _, a := range arr {
		if a != k {
			b = append(b, a)
		}
	}
	return b
}

// Unique Unique
func Unique(arr []int64) []int64 {
	m := make(map[int64]struct{})
	for _, a := range arr {
		m[a] = struct{}{}
	}
	res := make([]int64, 0)
	for a := range m {
		res = append(res, a)
	}
	return res
}

// GetTopReply GetTopReply
func (s *Service) GetTopReply(ctx context.Context, oid int64, otyp int8, topType uint32) (*model.Reply, error) {
	r, err := s.dao.Mc.GetTop(ctx, oid, otyp, topType)
	if err != nil {
		return nil, err
	}
	if r == nil {
		s.dao.Databus.AddTop(ctx, oid, otyp, topType)
		return nil, nil
	}
	return r, nil
}

// GetReplyFromDBByIDs GetReplyFromDBByIDs
func (s *Service) GetReplyFromDBByIDs(ctx context.Context, oid int64, otyp int8, ids []int64) ([]*model.Reply, error) {
	rs := make([]*model.Reply, 0)
	if len(ids) == 0 {
		return rs, nil
	}

	idReplyMap, err := s.dao.Reply.GetByIds(ctx, oid, otyp, ids)
	if err != nil {
		return nil, err
	}

	idReplyContentMap, err := s.dao.Content.GetByIds(ctx, oid, ids)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		if r, ok := idReplyMap[id]; ok {
			if r == nil {
				rs = append(rs, nil)
				continue
			}
			if content, ok := idReplyContentMap[id]; ok {
				r.Content = content
			}
			rs = append(rs, r)
		}
	}
	return rs, nil
}

// GetReplyByIDs GetReplyByIDs
func (s *Service) GetReplyByIDs(ctx context.Context, oid int64, otyp int8, ids []int64) (map[int64]*model.Reply, error) {
	res := make(map[int64]*model.Reply)
	if len(ids) == 0 {
		return res, nil
	}
	cachedReplies, missedIDs, err := s.dao.Mc.GetReplyByIDs(ctx, ids)
	var rs []*model.Reply
	if err != nil {
		rs, err = s.GetReplyFromDBByIDs(ctx, oid, otyp, ids)
		if err != nil {
			return nil, err
		}
		for _, r := range rs {
			res[r.RpID] = r
		}
		return res, nil
	}
	for _, r := range cachedReplies {
		res[r.RpID] = r
	}
	if len(missedIDs) == 0 {
		return res, nil
	}
	missedReplies, err := s.GetReplyFromDBByIDs(ctx, oid, otyp, missedIDs)
	if err != nil {
		return nil, err
	}
	select {
	case s.replyChan <- replyChan{rps: missedReplies}:
	default:
		log.Error("s.replyChan is full")
	}
	for _, r := range missedReplies {
		res[r.RpID] = r.Clone()
	}
	return res, nil
}

// GetChildrenIDsByCursor GetChildrenIDsByCursor
func (s *Service) GetChildrenIDsByCursor(ctx context.Context, sub *model.Subject, rootID int64, sort int8, cursor *model.Cursor) ([]int64, error) {
	k := reply.GenNewChildrenKeyByRootReplyID(rootID)
	cacheExist, err := s.dao.Redis.ExpireCache(ctx, k)
	if err != nil {
		return nil, err
	}
	var ids []int64
	if cacheExist {
		ids, err = s.dao.Redis.RangeChildrenIDByCursorScore(ctx, k, cursor)
		if err != nil {
			return nil, err
		}
		return ids, nil
	}
	s.dao.Databus.RecoverIndexByRoot(ctx, sub.Oid, rootID, sub.Type)
	switch sort {
	case model.SortByFloor:
		ids, err = s.dao.Reply.ChildrenIDSortByFloorCursor(ctx, sub.Oid, sub.Type, rootID, cursor)
	default:
		return nil, ecode.RequestErr
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetRootReplyIDsByCursor GetRootReplyIDsByCursor
func (s *Service) GetRootReplyIDsByCursor(ctx context.Context, sub *model.Subject, sort int8, cursor *model.Cursor) ([]int64, error) {
	var (
		ids   []int64
		isEnd bool
	)
	if sub.RCount == 0 {
		return []int64{}, nil
	}
	k := s.dao.Redis.CacheKeyRootReplyIDs(sub.Oid, sub.Type, sort)
	cacheExist, err := s.dao.Redis.ExpireCache(ctx, k)
	if err != nil {
		return nil, err
	}
	minFloor := cursor.Current() - 20
	if cursor.Latest() {
		minFloor = int64(sub.Count) - 20
	}
	if minFloor <= 0 {
		minFloor = 1
	}
	if cacheExist {
		ids, isEnd, err = s.dao.Redis.RangeRootIDByCursorScore(ctx, k, cursor)
		if err != nil {
			return nil, err
		}
		if sort == model.SortByFloor && len(ids) < cursor.Len() && !cursor.Increase() && !isEnd {
			ids, err = s.dao.Reply.RootIDSortByFloorCursor(ctx, sub.Oid, sub.Type, cursor)
			if err != nil {
				return nil, err
			}
			s.dao.Databus.RecoverFloorIdx(ctx, sub.Oid, sub.Type, int(minFloor), true)
		}
		return ids, nil
	}
	switch sort {
	case model.SortByFloor:
		s.dao.Databus.RecoverFloorIdx(ctx, sub.Oid, sub.Type, int(minFloor), true)
		ids, err = s.dao.Reply.RootIDSortByFloorCursor(ctx, sub.Oid, sub.Type, cursor)
	default:
		return nil, ecode.RequestErr
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetRootReplyIDs GetRootReplyIDs
func (s *Service) GetRootReplyIDs(ctx context.Context, oid int64, otyp int8, sort int8, offset, limit int64) ([]int64, error) {
	var ids []int64
	k := s.dao.Redis.CacheKeyRootReplyIDs(oid, otyp, sort)
	cacheExist, err := s.dao.Redis.ExpireCache(ctx, k)
	if err != nil {
		return nil, err
	}
	if cacheExist {
		ids, err = s.dao.Redis.RangeRootReplyIDs(ctx, k, int(offset), int(offset+limit-1))
		if err != nil {
			return nil, err
		}
		return ids, nil
	}
	s.dao.Databus.RecoverIndex(ctx, oid, otyp, sort)

	switch sort {
	case model.SortByFloor:
		ids, err = s.dao.Reply.GetIdsSortFloor(ctx, oid, otyp, int(offset), int(limit))
	case model.SortByCount:
		ids, err = s.dao.Reply.GetIdsSortCount(ctx, oid, otyp, int(offset), int(limit))
	case model.SortByLike:
		ids, err = s.dao.Reply.GetIdsSortLike(ctx, oid, otyp, int(offset), int(limit))
	default:
		log.Error("unsupported sort:%d", sort)
		return nil, ecode.RequestErr
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetSubject GetSubject
func (s *Service) GetSubject(ctx context.Context, oid int64, tp int8) (*model.Subject, error) {
	subject, err := s.getSubject(ctx, oid, tp)
	if err != nil {
		return nil, err
	}
	if subject.State == model.SubStateForbid {
		return nil, ecode.ReplyForbidReply
	}
	return subject, nil
}

func elementOf(k int64, arr []int64) bool {
	for _, i := range arr {
		if i == k {
			return true
		}
	}
	return false
}

// InsertInto insert `id` into sorted list(order by `cmp`) `ids`
// after insertion if the total length > `size`, truncate the extra element
func InsertInto(ids []int64, id int64, size int, cmp model.Comp) []int64 {
	if elementOf(id, ids) {
		return ids
	}
	if len(ids) < size {
		return model.SortArr(append(ids, id), cmp)
	}

	ids = model.SortArr(ids, cmp)
	if !withInRange(id, ids[0], ids[len(ids)-1]) {
		return ids
	}

	for i := 0; i < len(ids); i++ {
		if cmp(id, ids[i]) {
			ids = append(ids[:i], append([]int64{id}, ids[i:]...)...)
			break
		}
	}
	return ids[:size]
}

func withInRange(i, begin, end int64) bool {
	return (begin > i && end < i) || (begin < i && end > i)
}

// FillRootReplies FillRootReplies
func (s *Service) FillRootReplies(ctx context.Context,
	rs []*model.Reply,
	mid int64,
	ip string,
	htmlEscape bool,
	sub *model.Subject) {
	var (
		allReply        []*model.Reply
		allIDs, allMIDs []int64
	)
	if mid > 0 {
		allMIDs = append(allMIDs, mid)
	}
	for _, r := range rs {
		allIDs, allMIDs, allReply = collect(r, allIDs, allMIDs, allReply)
		for _, rr := range r.Replies {
			allIDs, allMIDs, allReply = collect(rr, allIDs, allMIDs, allReply)
		}
	}
	s.fillReplies(ctx, sub, allIDs, allReply, Unique(allMIDs), mid, ip, htmlEscape)
}

func (s *Service) fillReplies(ctx context.Context,
	sub *model.Subject,
	allReplyIDs []int64,
	rs []*model.Reply,
	mids []int64,
	reqMid int64,
	ip string,
	htmlEscape bool) {
	var (
		actionMap   map[int64]int8
		blackedMap  map[int64]bool
		relationMap map[int64]*accmdl.RelationReply
		assistMap   map[int64]int
		fansMap     map[int64]*model.FansDetail
		accountMap  map[int64]*accmdl.Card
	)
	g := errgroup.WithContext(ctx)
	if reqMid > 0 {
		g.Go(func(ctx context.Context) error {
			actionMap, _ = s.actions(ctx, reqMid, sub.Oid, allReplyIDs)
			return nil
		})
		g.Go(func(ctx context.Context) error {
			relationMap, _ = s.GetRelationMap(ctx, reqMid, mids, ip)
			return nil
		})
		g.Go(func(ctx context.Context) error {
			blackedMap, _ = s.GetBlacklistMap(ctx, reqMid, ip)
			return nil
		})
	}
	g.Go(func(ctx context.Context) error {
		accountMap, _ = s.GetAccountInfoMap(ctx, mids, ip)
		return nil
	})
	if !(s.IsWhiteAid(sub.Oid, sub.Type)) {
		g.Go(func(ctx context.Context) error {
			assistMap, _ = s.GetAssistMap(ctx, sub.Mid, ip)
			return nil
		})
		g.Go(func(ctx context.Context) error {
			fansMap, _ = s.GetFansMap(ctx, mids, sub.Mid, ip)
			return nil
		})
	}
	g.Wait()
	for _, r := range rs {
		s.fillReply(r,
			htmlEscape,
			accountMap,
			actionMap,
			fansMap,
			blackedMap,
			assistMap,
			relationMap)
	}
}

func (s *Service) fillReply(r *model.Reply,
	escape bool,
	accountMap map[int64]*accmdl.Card,
	actionMap map[int64]int8,
	fansMap map[int64]*model.FansDetail,
	blackedMap map[int64]bool,
	assistMap map[int64]int,
	relationMap map[int64]*accmdl.RelationReply) {

	if r == nil {
		return
	}
	r.FillFolder()
	r.FillStr(escape)

	if r.Content != nil {
		r.Content.FillAts(accountMap)
	}
	r.Action = actionMap[r.RpID]
	r.Member = new(model.Member)
	var (
		ok      bool
		blacked bool
		card    *accmdl.Card
	)
	if card, ok = accountMap[r.Mid]; ok {
		r.Member.Info = new(model.Info)
		r.Member.Info.FromCard(card)
	} else {
		r.Member.Info = new(model.Info)
		*r.Member.Info = *s.defMember
		r.Member.Info.Mid = strconv.FormatInt(r.Mid, 10)
	}
	if r.Member.FansDetail, ok = fansMap[r.Mid]; ok {
		r.FansGrade = r.Member.FansDetail.Status
	}
	if blacked, ok = blackedMap[r.Mid]; ok && blacked {
		r.State = model.ReplyStateBlacklist
	}
	if r.Replies == nil {
		r.Replies = []*model.Reply{}
	}
	if _, ok = assistMap[r.Mid]; ok {
		r.Assist = 1
	}
	if attetion, ok := relationMap[r.Mid]; ok {
		if attetion.Following {
			r.Member.Following = 1
		}
	}
	if r.RCount < 0 {
		r.RCount = 0
	}
}

// Fetch Fetch
func Fetch(idReplyMap map[int64]*model.Reply, ids []int64) []*model.Reply {
	res := make([]*model.Reply, 0, len(ids))
	for _, pid := range ids {
		if p, ok := idReplyMap[pid]; ok && p != nil {
			res = append(res, p)
		}
	}
	return res
}

// assemble insert children replies into their corresponding parents
func assemble(idReplyMap map[int64]*model.Reply, parentChildrenMap map[int64][]int64) map[int64]*model.Reply {
	parentIDs := make([]int64, 0)
	for pid := range parentChildrenMap {
		parentIDs = append(parentIDs, pid)
	}

	res := make(map[int64]*model.Reply)
	for _, pid := range parentIDs {
		if p, ok := idReplyMap[pid]; ok {
			if childrenIDs, ok := parentChildrenMap[pid]; ok {
				for _, childID := range childrenIDs {
					if r, ok := idReplyMap[childID]; ok {
						p.Replies = append(p.Replies, r)
					}
				}
			}
			res[pid] = p
		}
	}
	return res
}

// ParentChildrenReplyIDRelation ParentChildrenReplyIDRelation
func (s *Service) ParentChildrenReplyIDRelation(ctx context.Context, sub *model.Subject, parentIDs []int64) (map[int64][]int64, error) {
	idReplyMap, err := s.GetReplyByIDs(ctx, sub.Oid, sub.Type, parentIDs)
	if err != nil {
		return nil, err
	}

	var parentWithChildren, parentWithoutChildren []int64
	for id, reply := range idReplyMap {
		if reply.RCount > 0 {
			parentWithChildren = append(parentWithChildren, id)
		} else {
			parentWithoutChildren = append(parentWithoutChildren, id)
		}
	}
	parentChildrenIDRelation, err := s.parentChildrenReplyIDRelation(ctx, sub.Oid, sub.Type, parentWithChildren)
	if err != nil {
		return nil, err
	}
	for _, pid := range parentWithoutChildren {
		parentChildrenIDRelation[pid] = []int64{}
	}
	return parentChildrenIDRelation, nil
}

func (s *Service) parentChildrenReplyIDRelation(ctx context.Context, oid int64,
	tp int8, parentIDs []int64) (map[int64][]int64, error) {
	parentChildrenIDRelation, missedIDs, err := s.dao.Redis.ParentChildrenReplyIDMap(ctx, parentIDs, 0, 4)
	if err != nil {
		return nil, err
	}
	if len(missedIDs) > 0 {
		for _, rootID := range missedIDs {
			childrenIDs, err := s.dao.Reply.ChildrenIDsOfRootReply(ctx, oid, rootID, tp, 0, defaultChildrenSize)
			if err != nil {
				return nil, err
			}
			parentChildrenIDRelation[rootID] = childrenIDs
			s.dao.Databus.RecoverIndexByRoot(ctx, oid, rootID, tp)
		}
	}
	return parentChildrenIDRelation, nil
}

// GetAccountInfoMap fn
func (s *Service) GetAccountInfoMap(ctx context.Context, mids []int64, ip string) (map[int64]*accmdl.Card, error) {
	if len(mids) == 0 {
		return _emptyCards, nil
	}
	args := &accmdl.MidsReq{Mids: mids}
	res, err := s.acc.Cards3(ctx, args)
	if err != nil {
		log.Error("s.acc.Infos2(%v), error(%v)", args, err)
		return nil, err
	}
	return res.Cards, nil
}

// GetFansMap fn
func (s *Service) GetFansMap(ctx context.Context, uids []int64, mid int64, ip string) (map[int64]*model.FansDetail, error) {
	fans, err := s.fans.Fetch(ctx, uids, mid, time.Now())
	if err != nil {
		return nil, err
	}
	return fans, nil
}

// GetAssistMap fn
func (s *Service) GetAssistMap(ctx context.Context, mid int64, ip string) (assistMap map[int64]int, err error) {
	arg := &assmdl.ArgAssists{
		Mid:    mid,
		RealIP: ip,
	}
	assistMap = make(map[int64]int)
	ids, err := s.assist.AssistIDs(ctx, arg)
	if err != nil {
		log.Error("s.assist.AssistIDs(%v), error(%v)", arg, err)
		return
	}
	for _, id := range ids {
		assistMap[id] = 1
	}
	return
}

// GetRelationMap GetRelationMap
func (s *Service) GetRelationMap(ctx context.Context, mid int64, targetMids []int64, ip string) (map[int64]*accmdl.RelationReply, error) {
	if len(targetMids) == 0 {
		return _emptyRelations, nil
	}
	relations, err := s.acc.Relations3(ctx, &accmdl.RelationsReq{Mid: mid, Owners: targetMids, RealIp: ip})
	if err != nil {
		log.Error("s.acc.Relations2(%v, %v) error(%v)", mid, targetMids, err)
		return nil, err
	}
	return relations.Relations, nil
}

// GetBlacklistMap GetBlacklistMap
func (s *Service) GetBlacklistMap(ctx context.Context,
	mid int64, ip string) (map[int64]bool, error) {
	if mid == 0 {
		return _emptyBlackList, nil
	}
	args := &accmdl.MidReq{Mid: mid}
	blacklistMap, err := s.acc.Blacks3(ctx, args)
	if err != nil {
		log.Error("s.acc.Blacks(%v) error(%v)", args, err)
		return nil, err
	}
	return blacklistMap.BlackList, nil
}

// GetPendingReply GetPendingReply
func (s *Service) GetPendingReply(ctx context.Context, mid int64, oid int64, typ int8) (map[int64]*model.Reply, error) {
	// WARNING: here we assume that pending replies have no children
	// otherwise, we need to change logic here
	pendingIDs, err := s.dao.Redis.UserAuditReplies(ctx, mid, oid, typ)
	if err != nil {
		return nil, err
	}
	pendingIDReplyMap, err := s.GetReplyByIDs(ctx, oid, typ, pendingIDs)
	if err != nil {
		return nil, err
	}
	return pendingIDReplyMap, nil
}

// GetSubReplyListByCursor GetSubReplyListByCursor
func (s *Service) GetSubReplyListByCursor(ctx context.Context, params *model.CursorParams) (*model.RootReplyList, error) {
	var (
		hasFolded bool
	)
	sub, err := s.Subject(ctx, params.Oid, params.OTyp)
	if err != nil {
		return nil, err
	}
	rp, err := s.ReplyContent(ctx, params.Oid, params.RootID, params.OTyp)
	if err != nil {
		return nil, err
	}
	if rp.IsRoot() && rp.HasFolded() {
		hasFolded = true
	}
	if rp.Root != 0 {
		params.RootID = rp.Root
		root, _ := s.reply(ctx, 0, params.Oid, rp.Root, params.OTyp)
		if root != nil && rp.IsRoot() && rp.HasFolded() {
			hasFolded = true
		}
	}
	childrenIDs, err := s.GetChildrenIDsByCursor(ctx, sub, params.RootID, params.Sort, params.Cursor)
	if err != nil {
		return nil, err
	}
	// 这里是处理被折叠的评论的逻辑
	if params.ShowFolded && hasFolded {
		foldedRpIDs, _ := s.foldedRepliesCursor(ctx, sub, params.RootID, params.Cursor)
		if len(foldedRpIDs) > 0 {
			childrenIDs = append(childrenIDs, foldedRpIDs...)
			sort.Slice(childrenIDs, func(x, y int) bool { return childrenIDs[x] < childrenIDs[y] })
			length := len(childrenIDs)
			if length > params.Cursor.Len() {
				if params.Cursor.Descrease() {
					// 往楼层小的地方翻页， 对于子评论就是往上翻页，这个时候要从后往前截断
					childrenIDs = childrenIDs[length-params.Cursor.Len():]
				} else {
					childrenIDs = childrenIDs[:params.Cursor.Len()]
				}
			}
		}
	}
	parentChildrenIDRelation := map[int64][]int64{params.RootID: childrenIDs}
	idReplyMap, err := s.IDReplyMap(ctx, sub, parentChildrenIDRelation)
	if err != nil {
		return nil, err
	}
	if NeedInsertPendingReply(params, sub) {
		var pendingIDReplyMap map[int64]*model.Reply
		pendingIDReplyMap, err = s.GetPendingReply(ctx, params.Mid, sub.Oid, sub.Type)
		if err != nil {
			return nil, err
		}
		for id, r := range pendingIDReplyMap {
			if _, ok := idReplyMap[r.Root]; ok {
				parentChildrenIDRelation[r.Root] = InsertInto(parentChildrenIDRelation[r.Root], id, defaultChildrenSize, model.OrderASC)
				sub.ACount++
				idReplyMap[id] = r
			}
		}
	}
	rootReply := assemble(idReplyMap, parentChildrenIDRelation)[params.RootID]
	if rootReply == nil || rootReply.IsDeleted() {
		return nil, ecode.ReplyNotExist
	}
	max, min, err := cursorRange(rootReply.Replies, params.Sort)
	if err != nil {
		return nil, err
	}
	s.FillRootReplies(ctx, []*model.Reply{rootReply}, params.Mid, params.IP, params.HTMLEscape, sub)
	return &model.RootReplyList{
		Subject:        sub,
		Roots:          []*model.Reply{rootReply},
		CursorRangeMax: max,
		CursorRangeMin: min,
	}, nil
}

func cursorRange(rs []*model.Reply, sort int8) (max, min int64, err error) {
	if len(rs) > 0 {
		switch sort {
		case model.SortByFloor:
			// NOTE("这里是为了12月13号给bishi搞零时置顶子评论做的ios兼容逻辑")
			var head int64
			if rs[0].RpID != 1237270231 {
				head = int64(rs[0].Floor)
			} else {
				if len(rs) > 1 {
					head = int64(rs[1].Floor)
				} else {
					head = int64(1)
				}
			}
			tail := int64(rs[len(rs)-1].Floor)
			if model.OrderDESC(head, tail) {
				max, min = head, tail
			} else {
				max, min = tail, head
			}
			return
		default:
			err = errors.New("unsupported cursor type")
			log.Error("%v", err)
			return 0, 0, err
		}
	}
	return
}

// GetRootReplyListByCursor GetRootReplyListByCursor
func (s *Service) GetRootReplyListByCursor(ctx context.Context, params *model.CursorParams) (*model.RootReplyList, error) {
	params.HotSize = s.hotNum(params.Oid, params.OTyp)
	sub, err := s.Subject(ctx, params.Oid, params.OTyp)
	if err != nil {
		return nil, err
	}
	roots, err := s.RootReplyListByCursor(ctx, sub, params)
	if err != nil {
		return nil, err
	}
	max, min, err := cursorRange(roots, params.Sort)
	if err != nil {
		return nil, err
	}
	// WARN: rootIDs AND hotIDs may be overlapped
	var allRootReply []*model.Reply
	allRootReply = append(allRootReply, roots...)
	var header *model.RootReplyListHeader
	if needHeader(params.Cursor, len(roots)) {
		header, err = s.GetRootReplyListHeader(ctx, sub, params)
		if err != nil {
			return nil, err
		}
		allRootReply = append(allRootReply, header.Hots...)
		if header.TopAdmin != nil {
			allRootReply = append(allRootReply, header.TopAdmin)
		}
		if header.TopUpper != nil {
			allRootReply = append(allRootReply, header.TopUpper)
		}
	}
	s.FillRootReplies(ctx, allRootReply, params.Mid, params.IP, params.HTMLEscape, sub)
	return &model.RootReplyList{
		Subject:        sub,
		Roots:          roots,
		Header:         header,
		CursorRangeMax: max,
		CursorRangeMin: min,
	}, nil
}

// DialogMaxMinFloor return max and min floor in dialog
func (s *Service) DialogMaxMinFloor(c context.Context, oid int64, tp int8, root, dialog int64) (maxFloor, minFloor int, err error) {
	var (
		ok bool
	)
	if ok, err = s.dao.Redis.ExpireDialogIndex(c, dialog); err != nil {
		log.Error("s.dao.Redis.ExpireDialogIndex error (%v)", err)
		return
	}
	if ok {
		minFloor, maxFloor, err = s.dao.Redis.DialogMinMaxFloor(c, dialog)
	} else {
		minFloor, maxFloor, err = s.dao.Reply.GetDialogMinMaxFloor(c, oid, tp, root, dialog)
	}
	return
}

// DialogByCursor ...
func (s *Service) DialogByCursor(c context.Context, mid, oid int64, tp int8, root, dialog int64, cursor *model.Cursor) (rps []*model.Reply, dialogCursor *model.DialogCursor, dialogMeta *model.DialogMeta, err error) {
	var (
		ok    bool
		rpIDs []int64
		rpMap map[int64]*model.Reply
	)
	dialogCursor = new(model.DialogCursor)
	dialogMeta = new(model.DialogMeta)
	dialogMeta.MaxFloor, dialogMeta.MinFloor, err = s.DialogMaxMinFloor(c, oid, tp, root, dialog)
	if err != nil {
		log.Error("get max and min floor for dialog from redis or db error", err)
		return
	}
	if (cursor.Max() != 0 && cursor.Max() > int64(dialogMeta.MaxFloor)) || (cursor.Min() != 0 && cursor.Min() < int64(dialogMeta.MinFloor)) {
		log.Warn("cursor max %d min %d, dialogmeta max %d min %d", cursor.Max(), cursor.Min(), dialogMeta.MinFloor, dialogMeta.MinFloor)
		err = ecode.RequestErr
		return
	}
	if ok, err = s.dao.Redis.ExpireDialogIndex(c, dialog); err != nil {
		log.Error("s.dao.Redis.ExpireDialogIndex error (%v)", err)
		return
	}
	if ok {
		rpIDs, err = s.dao.Redis.DialogByCursor(c, dialog, cursor)
	} else {
		s.dao.Databus.RecoverDialogIdx(c, oid, tp, root, dialog)
		if cursor.Latest() {
			rpIDs, err = s.dao.Reply.GetIDsByDialogAsc(c, oid, tp, root, dialog, int64(dialogMeta.MinFloor), cursor.Len())
		} else if cursor.Descrease() {
			rpIDs, err = s.dao.Reply.GetIDsByDialogDesc(c, oid, tp, root, dialog, cursor.Current(), cursor.Len())
		} else if cursor.Increase() {
			rpIDs, err = s.dao.Reply.GetIDsByDialogAsc(c, oid, tp, root, dialog, cursor.Current(), cursor.Len())
		} else {
			err = ecode.RequestErr
		}
	}
	if err != nil {
		log.Error("dialog by cursor from redis or db error (%v)", err)
		return
	}
	rpMap, err = s.repliesMap(c, oid, tp, rpIDs)
	if err != nil {
		return
	}
	for _, rpid := range rpIDs {
		if r, ok := rpMap[rpid]; ok {
			rps = append(rps, r)
		}
	}
	if !sort.SliceIsSorted(rps, func(i, j int) bool { return rps[i].Floor < rps[j].Floor }) {
		sort.Slice(rps, func(i, j int) bool { return rps[i].Floor < rps[j].Floor })
	}
	sub, err := s.Subject(c, oid, tp)
	if err != nil {
		log.Error("s.dao.Subject.Get(%d, %d) error(%v)", oid, tp, err)
		return
	}
	if err = s.buildReply(c, sub, rps, mid, false); err != nil {
		return
	}
	dialogCursor.Size = len(rps)
	if dialogCursor.Size == 0 {
		return
	}
	dialogCursor.MinFloor = rps[0].Floor
	dialogCursor.MaxFloor = rps[dialogCursor.Size-1].Floor
	return
}

//  ...
func (s *Service) foldedRepliesCursor(c context.Context, sub *model.Subject, root int64, cursor *model.Cursor) (foldedRpIDs []int64, err error) {
	var (
		xcursor = new(xmodel.Cursor)
		max     = int(cursor.Max())
		min     = int(cursor.Min())
	)
	xcursor.Ps = cursor.Len()
	// 针对子评论的情况
	if cursor.Increase() {
		xcursor.Prev = min
	} else if cursor.Descrease() {
		xcursor.Next = max
	}
	return s.foldedReplies(c, sub, root, xcursor)
}
