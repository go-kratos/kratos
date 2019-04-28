package service

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"

	model "go-common/app/interface/main/reply/model/reply"
	accmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"go-common/library/sync/errgroup.v2"
)

const (
	// hot count web
	_webHotCnt = 3
	// 热评列表中至少点赞数
	_hotLikes = 3
	// 取5条，因为有可能管理员置顶和up置顶同时存在
	_hotFilters = 5
)

func withinFloor(rpIDs []int64, rpID int64, pn, ps int, asc bool) bool {
	if len(rpIDs) == 0 || (pn == 1 && len(rpIDs) < ps) {
		return true
	}
	first := rpIDs[0]
	last := rpIDs[len(rpIDs)-1]
	if asc {
		if (first < rpID && last > rpID) || (len(rpIDs) < ps && last < rpID) {
			return true
		}
	} else {
		if (pn == 1 && first < rpID) || (first > rpID && last < rpID) || (len(rpIDs) < ps && last > rpID) {
			return true
		}
	}
	return false
}

// SecondReplies return second replies.
func (s *Service) SecondReplies(c context.Context, mid, oid, rootID, jumpID int64, tp int8, pn, ps int, escape bool) (seconds []*model.Reply, root *model.Reply, upMid int64, toPn int, err error) {
	var (
		ok   bool
		sub  *model.Subject
		jump *model.Reply
	)
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if sub, err = s.Subject(c, oid, tp); err != nil {
		return
	}
	// jump to the child reply page list
	if jumpID > 0 {
		if jump, err = s.ReplyContent(c, oid, jumpID, tp); err != nil {
			return
		}
		if jump.Root == 0 {
			root = jump
			pn = 1
		} else {
			if root, err = s.ReplyContent(c, oid, jump.Root, tp); err != nil {
				return
			}
			if pos := s.getReplyPosByRoot(c, root, jump); pos > ps {
				pn = (pos-1)/ps + 1
			} else {
				pn = 1
			}
		}
	} else {
		if root, err = s.ReplyContent(c, oid, rootID, tp); err != nil {
			return
		}
		if root.Root != 0 {
			if root, err = s.ReplyContent(c, oid, root.Root, tp); err != nil {
				return
			}
		}
	}
	if root.IsDeleted() {
		err = ecode.ReplyDeleted
		return
	}
	upMid = sub.Mid
	toPn = pn
	// get reply second reply content
	rootMap := make(map[int64]*model.Reply, 1)
	rootMap[root.RpID] = root
	secondMap, _, _ := s.secondReplies(c, sub, rootMap, mid, pn, ps)
	if seconds, ok = secondMap[root.RpID]; !ok {
		seconds = _emptyReplies
	}
	// get reply dependency info
	rs := make([]*model.Reply, 0, len(seconds)+1)
	rs = append(rs, root)
	rs = append(rs, seconds...)
	if err = s.buildReply(c, sub, rs, mid, escape); err != nil {
		return
	}
	return
}

// JumpReplies jump to page by reply id.
func (s *Service) JumpReplies(c context.Context, mid, oid, rpID int64, tp int8, ps, sndPs int, escape bool) (roots, hots []*model.Reply, topAdmin, topUpper *model.Reply, sub *model.Subject, pn, sndPn, total int, err error) {
	var (
		rootPos, sndPos int
		rootRp, rp      *model.Reply
		fixedSeconds    []*model.Reply
	)
	if !model.LegalSubjectType(tp) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if sub, err = s.Subject(c, oid, tp); err != nil {
		return
	}
	if rp, err = s.ReplyContent(c, oid, rpID, tp); err != nil {
		return
	}
	if rp.Root == 0 && rp.Parent == 0 {
		rootPos = s.getReplyPos(c, sub, rp)
	} else {
		if rootRp, err = s.ReplyContent(c, oid, rp.Root, tp); err != nil {
			return
		}
		rootPos = s.getReplyPos(c, sub, rootRp)
		sndPos = s.getReplyPosByRoot(c, rootRp, rp)
	}
	// root page number
	pn = (rootPos-1)/ps + 1
	// second page number
	if sndPos > sndPs {
		sndPn = (sndPos-1)/sndPs + 1
	} else {
		sndPn = 1
	}
	// get reply content
	roots, seconds, total, err := s.rootReplies(c, sub, mid, model.SortByFloor, pn, ps, 1, sndPs)
	if err != nil {
		return
	}
	if rootRp != nil && rootRp.RCount > 0 {
		if fixedSeconds, err = s.repliesByRoot(c, oid, rootRp.RpID, tp, sndPn, sndPs); err != nil {
			return
		}
		for _, rp := range roots {
			if rp.RpID == rootRp.RpID {
				rp.Replies = fixedSeconds
			}
		}
	}
	// top and hots
	topAdmin, topUpper, hots, hseconds, err := s.topAndHots(c, sub, mid, true, true)
	if err != nil {
		log.Error("s.topAndHots(%d,%d,%d) error(%v)", oid, tp, mid, err)
		err = nil // degrade
	}
	rs := make([]*model.Reply, 0, len(roots)+len(seconds)+len(hseconds)+len(fixedSeconds)+2)
	rs = append(rs, roots...)
	rs = append(rs, seconds...)
	rs = append(rs, hseconds...)
	rs = append(rs, hots...)
	rs = append(rs, fixedSeconds...)
	if topAdmin != nil {
		rs = append(rs, topAdmin)
	}
	if topUpper != nil {
		rs = append(rs, topUpper)
	}
	if err = s.buildReply(c, sub, rs, mid, escape); err != nil {
		return
	}
	return
}

// RootReplies return a page list of reply.
func (s *Service) RootReplies(c context.Context, params *model.PageParams) (page *model.PageResult, err error) {
	if !model.LegalSubjectType(params.Type) {
		err = ecode.ReplyIllegalSubType
		return
	}
	sub, err := s.Subject(c, params.Oid, params.Type)
	if err != nil {
		return
	}
	topAdmin, topUpper, hots, hseconds, err := s.topAndHots(c, sub, params.Mid, params.NeedHot, params.NeedSecond)
	if err != nil {
		log.Error("s.topAndHots(%+v) error(%v)", params, err)
		err = nil // degrade
	}
	roots, seconds, total, err := s.rootReplies(c, sub, params.Mid, params.Sort, params.PageNum, params.PageSize, 1, s.sndDefCnt)
	if err != nil {
		return
	}
	rs := make([]*model.Reply, 0, len(roots)+len(hots)+len(hseconds)+len(seconds)+2)
	rs = append(rs, hots...)
	rs = append(rs, roots...)
	rs = append(rs, hseconds...)
	rs = append(rs, seconds...)
	if topAdmin != nil {
		rs = append(rs, topAdmin)
	}
	if topUpper != nil {
		rs = append(rs, topUpper)
	}
	if err = s.buildReply(c, sub, rs, params.Mid, params.Escape); err != nil {
		return
	}
	page = &model.PageResult{
		Subject:  sub,
		TopAdmin: topAdmin,
		TopUpper: topUpper,
		Hots:     hots,
		Roots:    roots,
		Total:    total,
		AllCount: sub.ACount,
	}
	return
}

func (s *Service) topAndHots(c context.Context, sub *model.Subject, mid int64, needNot, needSnd bool) (topAdmin, topUpper *model.Reply, hots, seconds []*model.Reply, err error) {
	var (
		ok        bool
		hotIDs    []int64
		rootMap   map[int64]*model.Reply
		secondMap map[int64][]*model.Reply
	)
	// get hot replies
	if needNot {
		if hotIDs, _, err = s.rootReplyIDs(c, sub, model.SortByLike, 1, s.hotNumWeb(sub.Oid, sub.Type)+2, false); err != nil {
			return
		}
		if rootMap, err = s.repliesMap(c, sub.Oid, sub.Type, hotIDs); err != nil {
			return
		}
	}
	// get top replies
	if topAdmin, err = s.topReply(c, sub, model.SubAttrAdminTop); err != nil {
		return
	}
	if topUpper, err = s.topReply(c, sub, model.SubAttrUpperTop); err != nil {
		return
	}
	if rootMap == nil {
		rootMap = make(map[int64]*model.Reply)
	}
	if topAdmin != nil {
		rootMap[topAdmin.RpID] = topAdmin
	}
	if topUpper != nil {
		if !topUpper.IsNormal() && sub.Mid != mid {
			topUpper = nil
		} else {
			rootMap[topUpper.RpID] = topUpper
		}
	}
	// get second replies
	if needSnd {
		if secondMap, seconds, err = s.secondReplies(c, sub, rootMap, mid, 1, s.sndDefCnt); err != nil {
			return
		}
		if topAdmin != nil {
			if topAdmin.Replies, ok = secondMap[topAdmin.RpID]; !ok {
				topAdmin.Replies = _emptyReplies
			}
		}
		if topUpper != nil {
			if topUpper.Replies, ok = secondMap[topUpper.RpID]; !ok {
				topUpper.Replies = _emptyReplies
			}
		}
	}
	if len(hotIDs) == 0 {
		hots = _emptyReplies
		return
	}
	hotSize := s.hotNumWeb(sub.Oid, sub.Type)
	for _, rootID := range hotIDs {
		if hotSize != _hotSizeWeb && len(hots) >= _hotSizeWeb {
			break
		} else if len(hots) >= hotSize {
			break
		}
		if rp, ok := rootMap[rootID]; ok && rp.Like >= _hotLikes && !rp.IsTop() {
			if rp.Replies, ok = secondMap[rp.RpID]; !ok {
				rp.Replies = _emptyReplies
			}
			hots = append(hots, rp)
		}
	}
	return
}

func (s *Service) rootReplies(c context.Context, sub *model.Subject, mid int64, msort int8, pn, ps, secondPn, secondPs int) (roots, seconds []*model.Reply, total int, err error) {
	var (
		rootMap map[int64]*model.Reply
	)
	// get root replies
	rootIDs, total, err := s.rootReplyIDs(c, sub, msort, pn, ps, true)
	if err != nil {
		return
	}
	if len(rootIDs) > 0 {
		if rootMap, err = s.repliesMap(c, sub.Oid, sub.Type, rootIDs); err != nil {
			return
		}
	}
	// get pending audit replies
	if msort == model.SortByFloor && sub.AttrVal(model.SubAttrAudit) == model.AttrYes {
		var (
			pendingTotal   int
			pendingIDs     []int64
			rootPendingMap map[int64]*model.Reply
		)
		if rootPendingMap, _, pendingTotal, err = s.userAuditReplies(c, mid, sub.Oid, sub.Type); err != nil {
			err = nil // degrade
		}
		if rootMap == nil {
			rootMap = make(map[int64]*model.Reply)
		}
		for _, rp := range rootPendingMap {
			if withinFloor(rootIDs, rp.RpID, pn, ps, false) {
				rootMap[rp.RpID] = rp
				pendingIDs = append(pendingIDs, rp.RpID)
			}
		}
		if len(pendingIDs) > 0 {
			rootIDs = append(rootIDs, pendingIDs...)
			sort.Sort(model.DescFloors(rootIDs))
		}
		sub.ACount += pendingTotal
	}
	if len(rootIDs) == 0 {
		roots = _emptyReplies
		return
	}
	// get second replies
	secondMap, seconds, err := s.secondReplies(c, sub, rootMap, mid, secondPn, secondPs)
	if err != nil {
		return
	}
	for _, rootID := range rootIDs {
		if rp, ok := rootMap[rootID]; ok {
			if rp.Replies, ok = secondMap[rp.RpID]; !ok {
				rp.Replies = _emptyReplies
			}
			if msort != model.SortByFloor {
				//if not sort by floor,can't contain the top comment
				if rp.IsTop() {
					continue
				}
			}
			roots = append(roots, rp)
		}
	}
	if roots == nil {
		roots = _emptyReplies
	}
	return
}

func (s *Service) secondReplies(c context.Context, sub *model.Subject, rootMap map[int64]*model.Reply, mid int64, pn, ps int) (res map[int64][]*model.Reply, rs []*model.Reply, err error) {
	var (
		rootIDs, secondIDs []int64
		secondIdxMap       map[int64][]int64
		secondMap          map[int64]*model.Reply
	)
	for rootID, info := range rootMap {
		if info.RCount > 0 {
			rootIDs = append(rootIDs, rootID)
		}
	}
	if len(rootIDs) > 0 {
		if secondIdxMap, secondIDs, err = s.getIdsByRoots(c, sub.Oid, rootIDs, sub.Type, pn, ps); err != nil {
			return
		}
		if secondMap, err = s.repliesMap(c, sub.Oid, sub.Type, secondIDs); err != nil {
			return
		}
	}
	// get pending audit replies
	if sub.AttrVal(model.SubAttrAudit) == model.AttrYes {
		var secondPendings map[int64][]*model.Reply
		if _, secondPendings, _, err = s.userAuditReplies(c, mid, sub.Oid, sub.Type); err != nil {
			err = nil // degrade
		}
		if secondIdxMap == nil {
			secondIdxMap = make(map[int64][]int64)
		}
		if secondMap == nil {
			secondMap = make(map[int64]*model.Reply)
		}
		for rootID, rs := range secondPendings {
			var pendingIDs []int64
			if r, ok := rootMap[rootID]; ok {
				for _, r := range rs {
					if withinFloor(secondIdxMap[rootID], r.RpID, pn, ps, true) {
						secondMap[r.RpID] = r
						pendingIDs = append(pendingIDs, r.RpID)
					}
				}
				r.RCount += len(rs)
			}
			if len(pendingIDs) > 0 {
				secondIdxMap[rootID] = append(secondIdxMap[rootID], pendingIDs...)
				sort.Sort(model.AscFloors(secondIdxMap[rootID]))
			}
		}
	}
	res = make(map[int64][]*model.Reply, len(secondIdxMap))
	for root, idxs := range secondIdxMap {
		seconds := make([]*model.Reply, 0, len(idxs))
		for _, rpid := range idxs {
			if r, ok := secondMap[rpid]; ok {
				seconds = append(seconds, r)
			}
		}
		res[root] = seconds
		rs = append(rs, seconds...)
	}
	return
}

// FilDelReply delete reply which is deleted
func (s *Service) FilDelReply(rps []*model.Reply) (filtedRps []*model.Reply) {
	for _, rp := range rps {
		if !rp.IsDeleted() {
			var childs []*model.Reply
			for _, child := range rp.Replies {
				if !child.IsDeleted() {
					childs = append(childs, child)
				}
			}
			rp.Replies = childs
			filtedRps = append(filtedRps, rp)
		}
	}
	return
}

// EmojiReplaceI EmojiReplace international
func (s *Service) EmojiReplaceI(mobiAPP string, build int64, roots ...*model.Reply) {
	if mobiAPP == "android_i" && build > 1125000 && build < 2005000 {
		for _, root := range roots {
			if root != nil {
				if root.Content != nil {
					if emoCodes := _emojiCode.FindAllString(root.Content.Message, -1); len(emoCodes) > 0 {
						root.Content.Message = RepressEmotions(root.Content.Message, emoCodes)
					}
				}
				for _, rp := range root.Replies {
					if rp.Content != nil {
						if emoCodes := _emojiCode.FindAllString(rp.Content.Message, -1); len(emoCodes) > 0 {
							rp.Content.Message = RepressEmotions(rp.Content.Message, emoCodes)
						}
					}
				}
			}
		}
	}
}

// EmojiReplace EmojiReplace
func (s *Service) EmojiReplace(plat int8, build int64, roots ...*model.Reply) {
	if (plat == model.PlatIPad || plat == model.PlatPadHd || plat == model.PlatIPhone) && build <= 8170 {
		for _, root := range roots {
			if root != nil {
				for _, rp := range root.Replies {
					if rp.Content != nil {
						if emoCodes := _emojiCode.FindAllString(rp.Content.Message, -1); len(emoCodes) > 0 {
							rp.Content.Message = RepressEmotions(rp.Content.Message, emoCodes)
						}
					}
				}
			}
		}
	}
}

func (s *Service) buildReply(c context.Context, sub *model.Subject, replies []*model.Reply, mid int64, escape bool) (err error) {
	var (
		ok          bool
		assistMap   map[int64]int
		actionMap   map[int64]int8
		blackedMap  map[int64]bool
		attetionMap map[int64]*accmdl.RelationReply
		rpIDs       = make([]int64, 0, len(replies))
		mids        = make([]int64, 0, len(replies))
		uniqMids    = make(map[int64]struct{}, len(replies))
		fansMap     map[int64]*model.FansDetail
		accMap      map[int64]*accmdl.Card
	)
	if len(replies) == 0 {
		return
	}
	for _, rp := range replies {
		if rp.Content != nil {
			for _, mid := range rp.Content.Ats {
				uniqMids[mid] = struct{}{}
			}
		}
		uniqMids[rp.Mid] = struct{}{}
		rpIDs = append(rpIDs, rp.RpID)
	}
	for mid := range uniqMids {
		mids = append(mids, mid)
	}
	g := errgroup.WithContext(c)
	if mid > 0 {
		g.Go(func(c context.Context) error {
			actionMap, _ = s.actions(c, mid, sub.Oid, rpIDs)
			return nil
		})
		g.Go(func(c context.Context) error {
			attetionMap, _ = s.getAttentions(c, mid, mids)
			return nil
		})
		g.Go(func(c context.Context) error {
			blackedMap, _ = s.GetBlacklist(c, mid)
			return nil
		})
	}
	g.Go(func(c context.Context) error {
		accMap, _ = s.getAccInfo(c, mids)
		return nil
	})
	if !s.IsWhiteAid(sub.Oid, sub.Type) {
		g.Go(func(c context.Context) error {
			fansMap, _ = s.FetchFans(c, mids, sub.Mid)
			return nil
		})
		g.Go(func(c context.Context) error {
			assistMap = s.getAssistList(c, sub.Mid)
			return nil
		})
	}
	g.Wait()
	// set reply info
	for _, r := range replies {
		r.FillFolder()
		r.FillStr(escape)
		if r.Content != nil {
			r.Content.FillAts(accMap)
		}
		r.Action = actionMap[r.RpID]
		// member info and degrade
		r.Member = new(model.Member)
		var card *accmdl.Card
		if card, ok = accMap[r.Mid]; ok {
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
		if _, ok = blackedMap[r.Mid]; ok {
			r.State = model.ReplyStateBlacklist
		}
		if _, ok = assistMap[r.Mid]; ok {
			r.Assist = 1
		}
		if attetion, ok := attetionMap[r.Mid]; ok {
			if attetion.Following {
				r.Member.Following = 1
			}
		}
		// temporary fix situation: rcount < 0
		if r.RCount < 0 {
			r.RCount = 0
		}
	}
	return
}

func (s *Service) rootReplyIDs(c context.Context, sub *model.Subject, sort int8, pn, ps int, loadCount bool) (rpIDs []int64, count int, err error) {
	var (
		ok    bool
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	if count = sub.RCount; start >= count {
		return
	}
	if sort == model.SortByLike {
		mid := metadata.Int64(c, metadata.Mid)
		res, err := s.replyHotFeed(c, mid, sub.Oid, int(sub.Type), pn, ps)
		if err == nil && res.RpIDs != nil && len(res.RpIDs) > 0 {
			log.Info("reply-feed(test): reply abtest mid(%d) oid(%d) type (%d) test name(%s) rpIDs(%v)", mid, sub.Oid, sub.Type, res.Name, res.RpIDs)
			rpIDs = res.RpIDs
			if loadCount {
				count = int(res.Count)
			}
			return rpIDs, count, nil
		}
		if err != nil {
			log.Error("reply-feed error(%v)", err)
			err = nil
		} else {
			log.Info("reply-feed(origin): reply abtest mid(%d) oid(%d) type (%d) test name(%s) rpIDs(%v)", mid, sub.Oid, sub.Type, res.Name, res.RpIDs)
		}
	}
	// expire the index cache
	if ok, err = s.dao.Redis.ExpireIndex(c, sub.Oid, sub.Type, sort); err != nil {
		log.Error("s.dao.Redis.ExpireIndex(%d,%d,%d) error(%v)", sub.Oid, sub.Type, sort, err)
		return
	}
	if !ok {
		// here we can get that Serviceub.RCount > 0
		switch sort {
		case model.SortByFloor:
			s.dao.Databus.RecoverFloorIdx(c, sub.Oid, sub.Type, end+1, false)
			rpIDs, err = s.dao.Reply.GetIdsSortFloor(c, sub.Oid, sub.Type, start, ps)
		case model.SortByCount:
			s.dao.Databus.RecoverIndex(c, sub.Oid, sub.Type, sort)
			rpIDs, err = s.dao.Reply.GetIdsSortCount(c, sub.Oid, sub.Type, start, ps)
		case model.SortByLike:
			s.dao.Databus.RecoverIndex(c, sub.Oid, sub.Type, sort)
			if rpIDs, err = s.dao.Reply.GetIdsSortLike(c, sub.Oid, sub.Type, start, ps); err != nil {
				return
			}
			if loadCount {
				count, err = s.dao.Reply.CountLike(c, sub.Oid, sub.Type)
			}
		}
		if err != nil {
			log.Error("s.rootIDs(%d,%d,%d,%d,%d) error(%v)", sub.Oid, sub.Type, sort, start, ps, err)
			return
		}
	} else {
		var isEnd bool
		if rpIDs, isEnd, err = s.dao.Redis.Range(c, sub.Oid, sub.Type, sort, start, end); err != nil {
			log.Error("s.dao.Redis.Range(%d,%d,%d,%d,%d) error(%v)", sub.Oid, sub.Type, sort, start, end, err)
			return
		}
		if (sort == model.SortByLike || sort == model.SortByCount) && loadCount {
			if count, err = s.dao.Redis.CountReplies(c, sub.Oid, sub.Type, sort); err != nil {
				log.Error("s.dao.Redis.CountLike(%d,%d,%d) error(%v)", sub.Oid, sub.Type, sort, err)
			}
		}
		if sort == model.SortByFloor && len(rpIDs) < ps && !isEnd {
			//The addition and deletion of comments may result in the display of duplicate entries
			rpIDs, err = s.dao.Reply.GetIdsSortFloor(c, sub.Oid, sub.Type, start, ps)
			if err != nil {
				log.Error("s.rootIDs(%d,%d,%d,%d,%d) error(%v)", sub.Oid, sub.Type, sort, start, ps, err)
				return
			}
			s.dao.Databus.RecoverFloorIdx(c, sub.Oid, sub.Type, end+1, false)
		}
	}
	return
}

// topReply return top replies from cache.
func (s *Service) topReply(c context.Context, sub *model.Subject, top uint32) (rp *model.Reply, err error) {
	if top != model.ReplyAttrUpperTop && top != model.ReplyAttrAdminTop {
		return
	}
	if sub.AttrVal(top) == model.AttrYes && sub.Meta != "" {
		var meta model.SubjectMeta
		err = json.Unmarshal([]byte(sub.Meta), &meta)
		if err != nil {
			log.Error("s.topReply(%d,%d,%d) unmarshal error(%v)", sub.Oid, sub.Type, top, err)
			return
		}
		var rpid int64
		if top == model.SubAttrAdminTop && meta.AdminTop != 0 {
			rpid = meta.AdminTop

		} else if top == model.SubAttrUpperTop && meta.UpperTop != 0 {
			rpid = meta.UpperTop
		}
		if rpid != 0 {
			rp, err = s.ReplyContent(c, sub.Oid, rpid, sub.Type)
			if err != nil {
				log.Error("s.GetReply(%d,%d,%d) error(%v)", sub.Oid, sub.Type, rpid, err)
				return
			}
			if rp == nil {
				log.Error("s.GetReply(%d,%d,%d) is nil", sub.Oid, sub.Type, rpid)
			}
			return
		}
	}

	if sub.AttrVal(top) == model.AttrYes {
		if rp, err = s.dao.Mc.GetTop(c, sub.Oid, sub.Type, top); err != nil {
			log.Error("s.dao.Mc.GetTop(%d,%d,%d) error(%v)", sub.Oid, sub.Type, top, err)
			return
		}
		if rp == nil {
			s.dao.Databus.AddTop(c, sub.Oid, sub.Type, top)
		}
	}
	return
}

func (s *Service) userAuditReplies(c context.Context, mid, oid int64, tp int8) (rootMap map[int64]*model.Reply, secondMap map[int64][]*model.Reply, total int, err error) {
	rpIDs, err := s.dao.Redis.UserAuditReplies(c, mid, oid, tp)
	if err != nil {
		log.Error("s.dao.Redis.Range(%d,%d,%d) error(%v)", oid, tp, mid, err)
		return
	}
	rpMap, err := s.repliesMap(c, oid, tp, rpIDs)
	if err != nil {
		return
	}
	total = len(rpMap)
	rootMap = make(map[int64]*model.Reply)
	secondMap = make(map[int64][]*model.Reply)
	for _, rp := range rpMap {
		if rp.Root == 0 {
			if !rp.IsTop() {
				rootMap[rp.RpID] = rp
			}
		} else {
			secondMap[rp.Root] = append(secondMap[rp.Root], rp)
		}
	}
	return
}

// repliesMap multi get reply from cache or db when missed and fill content.
func (s *Service) repliesMap(c context.Context, oid int64, tp int8, rpIDs []int64) (res map[int64]*model.Reply, err error) {
	if len(rpIDs) == 0 {
		return
	}
	res, missIDs, err := s.dao.Mc.GetMultiReply(c, rpIDs)
	if err != nil {
		log.Error("s.dao.Mc.GetMultiReply(%d,%d,%d) error(%v)", oid, tp, rpIDs, err)
		err = nil
		res = make(map[int64]*model.Reply, len(rpIDs))
		missIDs = rpIDs
	}
	if len(missIDs) > 0 {
		var (
			mrp map[int64]*model.Reply
			mrc map[int64]*model.Content
		)
		if mrp, err = s.dao.Reply.GetByIds(c, oid, tp, missIDs); err != nil {
			log.Error("s.reply.GetByIds(%d,%d,%d) error(%v)", oid, tp, rpIDs, err)
			return
		}
		if mrc, err = s.dao.Content.GetByIds(c, oid, missIDs); err != nil {
			log.Error("s.content.GetByIds(%d,%d) error(%v)", oid, rpIDs, err)
			return
		}
		rs := make([]*model.Reply, 0, len(missIDs))
		for _, rpID := range missIDs {
			if rp, ok := mrp[rpID]; ok {
				rp.Content = mrc[rpID]
				res[rpID] = rp
				rs = append(rs, rp.Clone())
			}
		}
		// asynchronized add reply cache
		select {
		case s.replyChan <- replyChan{rps: rs}:
		default:
			log.Warn("s.replyChan is full")
		}
	}
	return
}

// ReplyContent get reply and content.
func (s *Service) ReplyContent(c context.Context, oid, rpID int64, tp int8) (r *model.Reply, err error) {
	if r, err = s.dao.Mc.GetReply(c, rpID); err != nil {
		log.Error("replyCacheDao.GetReply(%d, %d, %d) error(%v)", oid, rpID, tp, err)
	}
	if r == nil {
		if r, err = s.dao.Reply.Get(c, oid, rpID); err != nil {
			log.Error("s.reply.GetReply(%d, %d) error(%v)", oid, rpID, err)
			return nil, err
		}
		if r == nil {
			return nil, ecode.ReplyNotExist
		}
		if r.Content, err = s.dao.Content.Get(c, oid, rpID); err != nil {
			return nil, err
		}
		if err = s.dao.Mc.AddReply(c, r); err != nil {
			log.Error("mc.AddReply(%d,%d,%d) error(%v)", oid, rpID, tp, err)
			err = nil
		}
	}
	if r.Oid != oid || r.Type != tp {
		log.Warn("reply dismatches with parameter, oid: %d, rpID: %d, tp: %d, actual: %d, %d, %d", oid, rpID, tp, r.Oid, r.RpID, r.Type)
		return nil, ecode.RequestErr
	}
	return r, nil
}

func (s *Service) repliesByRoot(c context.Context, oid, root int64, tp int8, pn, ps int) (res []*model.Reply, err error) {
	var (
		cache bool
		rpIDs []int64
		start = (pn - 1) * ps
		end   = start + ps - 1
	)
	if cache, err = s.dao.Redis.ExpireIndexByRoot(c, root); err != nil {
		return
	}
	if cache {
		if rpIDs, err = s.dao.Redis.RangeByRoot(c, root, start, end); err != nil {
			log.Error("s.dao.Redis.RangeByRoots() err(%v)", err)
			return
		}
	} else {
		if rpIDs, err = s.dao.Reply.GetIdsByRoot(c, oid, root, tp, start, ps); err != nil {
			log.Error("s.dao.Reply.GetIdsByRoot(oid %d,tp %d,root %d) err(%v)", oid, tp, root, err)
		}
		s.dao.Databus.RecoverIndexByRoot(c, oid, root, tp)
	}
	rs, err := s.repliesMap(c, oid, tp, rpIDs)
	if err != nil {
		return
	}
	for _, rpID := range rpIDs {
		if rp, ok := rs[rpID]; ok {
			res = append(res, rp)
		}
	}
	return
}

// ReplyHots return the hot replies.
func (s *Service) ReplyHots(c context.Context, oid int64, typ int8, pn, ps int) (sub *model.Subject, res []*model.Reply, err error) {
	if !model.LegalSubjectType(typ) {
		err = ecode.ReplyIllegalSubType
		return
	}
	if sub, err = s.Subject(c, oid, typ); err != nil {
		log.Error("s.Subject(%d,%d) error(%v)", oid, typ, err)
		return
	}
	hotIDs, _, err := s.rootReplyIDs(c, sub, model.SortByLike, pn, ps, false)
	if err != nil {
		return
	}
	rootMap, err := s.repliesMap(c, sub.Oid, sub.Type, hotIDs)
	if err != nil {
		return
	}
	for _, rpID := range hotIDs {
		if rp, ok := rootMap[rpID]; ok {
			res = append(res, rp)
		}
	}
	if len(res) == 0 {
		res = _emptyReplies
	}
	return
}

// Dialog ...
func (s *Service) Dialog(c context.Context, mid, oid int64, tp int8, root, dialog int64, pn, ps int, escape bool) (rps []*model.Reply, err error) {
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1
		ok    bool
		rpIDs []int64
	)
	if ok, err = s.dao.Redis.ExpireDialogIndex(c, dialog); err != nil {
		log.Error("s.dao.Redis.ExpireDialogIndex error (%v)", err)
		return
	}
	if ok {
		rpIDs, err = s.dao.Redis.RangeRpsByDialog(c, dialog, start, end)
	} else {
		s.dao.Databus.RecoverDialogIdx(c, oid, tp, root, dialog)
		rpIDs, err = s.dao.Reply.GetIDsByDialog(c, oid, tp, root, dialog, start, ps)
	}
	if err != nil {
		log.Error("range replies by dialog from redis or db error (%v)", err)
		return
	}
	rpMap, err := s.repliesMap(c, oid, tp, rpIDs)
	if err != nil {
		return
	}
	for _, rpID := range rpIDs {
		if r, ok := rpMap[rpID]; ok {
			rps = append(rps, r)
		}
	}

	sub, err := s.Subject(c, oid, tp)
	if err != nil {
		return
	}
	for _, rp := range rps {
		rp.DialogStr = strconv.FormatInt(rp.Dialog, 10)
	}
	if err = s.buildReply(c, sub, rps, mid, escape); err != nil {
		return
	}
	return
}
