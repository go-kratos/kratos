package service

import (
	"context"
	"sort"

	"go-common/app/interface/main/reply/model/reply"
	xmodel "go-common/app/interface/main/reply/model/xreply"

	"go-common/library/sync/errgroup.v2"
)

const (
	_secondRepliesPn = 1
	_secondRepliesPs = 3
	_hotSize         = 5
	_hotSizeWeb      = 3
)

func fillHot(rs []*reply.Reply, sub *reply.Subject, maxSize int) (hots []*reply.Reply) {
	for _, r := range rs {
		if r.Like >= 3 && !r.IsTop() && !r.IsDeleted() {
			hots = append(hots, r)
		}
	}
	if hots == nil {
		hots = _emptyReplies
	} else if len(hots) > maxSize {
		hots = hots[:maxSize]
	}
	if sub.RCount <= 20 && len(hots) > 0 {
		hots = _emptyReplies
	}
	return hots
}

// buildResCursor ...
func buildResCursor(reqCursor *xmodel.Cursor, replies []*reply.Reply, cursorMode int) (resCursor xmodel.CursorRes) {
	length := len(replies)
	switch cursorMode {
	case xmodel.CursorModePage:
		// next和prev只可能一个不为0
		curPage := reqCursor.Next + reqCursor.Prev
		if reqCursor.Latest() {
			resCursor.IsBegin = true
			if length == 0 {
				resCursor.Next = curPage
				resCursor.Prev = curPage
				resCursor.IsEnd = true
			} else {
				// 下一页是第二页
				resCursor.Next = 2
				resCursor.Prev = 1
				if length < reqCursor.Ps {
					resCursor.IsEnd = true
				}
			}
		} else if reqCursor.Forward() {
			resCursor.Next = curPage + 1
			resCursor.Prev = curPage
			if length < reqCursor.Ps {
				resCursor.IsEnd = true
			}
		} else if reqCursor.Backward() {
			resCursor.Next = curPage
			resCursor.Prev = curPage - 1
			if length < reqCursor.Ps {
				resCursor.IsBegin = true
			}
		}
	case xmodel.CursorModeCursor:
		if reqCursor.Latest() {
			resCursor.IsBegin = true
			// Latest, 点进来默认访问的情况
			if length == 0 {
				resCursor.Next = reqCursor.Next
				resCursor.Prev = reqCursor.Prev
				resCursor.IsEnd = true
			} else {
				resCursor.Next = replies[length-1].Floor
				resCursor.Prev = replies[0].Floor
				if length < reqCursor.Ps {
					resCursor.IsEnd = true
				}
			}
		} else if reqCursor.Forward() {
			if length == 0 {
				resCursor.Next = reqCursor.Next
				resCursor.Prev = reqCursor.Prev
				resCursor.IsEnd = true
			} else {
				resCursor.Next = replies[length-1].Floor
				resCursor.Prev = replies[0].Floor
				if length < reqCursor.Ps {
					resCursor.IsEnd = true
				}
			}
		} else if reqCursor.Backward() {
			if length == 0 {
				resCursor.Next = reqCursor.Next
				resCursor.Prev = reqCursor.Prev
				resCursor.IsBegin = true
			} else {
				resCursor.Next = replies[length-1].Floor
				resCursor.Prev = replies[0].Floor
				if length < reqCursor.Ps {
					resCursor.IsBegin = true
				}
			}
		}
	}
	return
}

// replyMisc ...
func (s *Service) replyCommonRes(ctx context.Context, mid, oid int64, tp int8, sub *reply.Subject) (res xmodel.CommonRes) {
	if sub == nil {
		return
	}
	g := errgroup.WithContext(ctx)
	g.Go(func(ctx context.Context) error {
		if s.RelationBlocked(ctx, sub.Mid, mid) {
			res.Blacklist = 1
		}
		return nil
	})
	g.Go(func(ctx context.Context) error {
		if ok, _ := s.CheckAssist(ctx, sub.Mid, mid); ok {
			res.Assist = 1
		}
		return nil
	})
	// 默认都是展示
	res.Config.ShowAdmin, res.Config.ShowEntry, res.Config.ShowFloor = 1, 1, 1
	g.Go(func(ctx context.Context) error {
		if cfg, _ := s.GetReplyLogConfig(ctx, sub, 1); cfg != nil {
			res.Config.ShowAdmin, res.Config.ShowEntry = cfg.ShowAdmin, cfg.ShowEntry
		}
		// 特殊稿件不显示楼层
		if !s.ShowFloor(oid, tp) {
			res.Config.ShowFloor = 0
		}
		return nil
	})
	g.Wait()
	res.Upper.Mid = sub.Mid
	return
}

// fillXreplyRes ...
func (s *Service) fillXreplyRes(ctx context.Context, sub *reply.Subject, req *xmodel.ReplyReq, res *xmodel.ReplyRes) {
	res.CommonRes = s.replyCommonRes(ctx, req.Mid, req.Oid, req.Type, sub)
	res.Notice = s.RplyNotice(ctx, req.Plat, req.Build, req.MobiApp, req.Buvid)
	// 过滤掉可能的脏数据
	res.Replies = s.FilDelReply(res.Replies)
	res.Hots = s.FilDelReply(res.Hots)
	// 过滤掉赞数小于3的
	res.Hots = fillHot(res.Hots, sub, s.hotNum(req.Oid, req.Type))
	res.Cursor.AllCount = sub.ACount
	res.Folder = sub.Folder()
}

// Xreply ...
func (s *Service) Xreply(ctx context.Context, req *xmodel.ReplyReq) (res *xmodel.ReplyRes, err error) {
	var (
		sub       *reply.Subject
		cursor    *reply.Cursor
		cursorRes *reply.RootReplyList
		pageRes   *reply.PageResult
	)
	res = new(xmodel.ReplyRes)
	mode, supportMode := req.ModeInfo(s.sortByHot, s.sortByTime)
	switch mode {
	case xmodel.ModeOrigin, xmodel.ModeTime:
		// 这里原来用的闭区间，所以这里有坑
		if req.Cursor.Forward() {
			if req.Cursor.Next > 1 {
				req.Cursor.Next--
			}
		} else if req.Cursor.Backward() {
			req.Cursor.Prev++
		}
		if cursor, err = reply.NewCursor(int64(req.Cursor.Next), int64(req.Cursor.Prev), req.Cursor.Ps, reply.OrderDESC); err != nil {
			return
		}
		if req.Cursor.Backward() {
			if cursor, err = reply.NewCursor(int64(req.Cursor.Next), int64(req.Cursor.Prev), req.Cursor.Ps, reply.OrderASC); err != nil {
				return
			}
		}
		cursorParams := &reply.CursorParams{
			IP:         req.IP,
			Oid:        req.Oid,
			OTyp:       req.Type,
			Sort:       reply.SortByFloor,
			HTMLEscape: false,
			Cursor:     cursor,
			HotSize:    s.hotNum(req.Oid, req.Type),
			Mid:        req.Mid,
		}
		if cursorRes, err = s.GetRootReplyListByCursor(ctx, cursorParams); err != nil {
			return
		}
		sub = cursorRes.Subject
		res.Replies = cursorRes.Roots
		if cursorRes.Header != nil {
			res.Top.Admin = cursorRes.Header.TopAdmin
			res.Top.Upper = cursorRes.Header.TopUpper
			res.Hots = cursorRes.Header.Hots
			// 纯按楼层排序 去掉热评
			if mode == xmodel.ModeTime {
				res.Hots = _emptyReplies
			}
		}
		if !sort.SliceIsSorted(res.Replies, func(i, j int) bool { return res.Replies[i].Floor > res.Replies[j].Floor }) {
			sort.Slice(res.Replies, func(i, j int) bool { return res.Replies[i].Floor > res.Replies[j].Floor })
		}

		// 这里原来用的闭区间，所以这里有坑
		if req.Cursor.Next == 1 {
			res.Replies = _emptyReplies
		}
		// 由服务端来控制翻页逻辑
		res.Cursor = buildResCursor(&xmodel.Cursor{Prev: req.Cursor.Prev, Next: req.Cursor.Next, Ps: req.Cursor.Ps}, res.Replies, xmodel.CursorModeCursor)
	case xmodel.ModeHot:
		var curPage = req.Cursor.Prev + req.Cursor.Next
		if req.Cursor.Latest() {
			curPage = 1
		}
		pageParams := &reply.PageParams{
			Mid:        req.Mid,
			Oid:        req.Oid,
			Type:       req.Type,
			Sort:       reply.SortByLike,
			PageNum:    curPage,
			PageSize:   req.Cursor.Ps,
			NeedSecond: true,
			Escape:     false,
			NeedHot:    false,
		}
		if pageRes, err = s.RootReplies(ctx, pageParams); err != nil {
			return
		}
		sub = pageRes.Subject
		res.Replies = pageRes.Roots
		res.Top.Admin = pageRes.TopAdmin
		res.Top.Upper = pageRes.TopUpper
		// 按页码翻页控制返回页码
		res.Cursor = buildResCursor(&xmodel.Cursor{Prev: req.Cursor.Prev, Next: req.Cursor.Next, Ps: req.Cursor.Ps}, res.Replies, xmodel.CursorModePage)
	}
	res.Cursor.Mode = mode
	res.Cursor.SupportMode = supportMode
	s.fillXreplyRes(ctx, sub, req, res)
	return
}

// SubFoldedReply ...
func (s *Service) SubFoldedReply(ctx context.Context, req *xmodel.SubFolderReq) (res *xmodel.SubFolderRes, err error) {
	var (
		rpIDs       []int64
		rootMap     map[int64]*reply.Reply
		childrenMap map[int64][]*reply.Reply
		rootRps     []*reply.Reply
		childrenRps []*reply.Reply
		sub         *reply.Subject
	)
	res = new(xmodel.SubFolderRes)
	if req.Cursor.Backward() {
		return
	}
	cursor := &xmodel.Cursor{
		Ps:   req.Cursor.Ps,
		Next: req.Cursor.Next,
	}
	if sub, err = s.Subject(ctx, req.Oid, req.Type); err != nil {
		return
	}
	if rpIDs, err = s.foldedReplies(ctx, sub, 0, cursor); err != nil {
		return
	}
	if rootMap, err = s.repliesMap(ctx, req.Oid, req.Type, rpIDs); err != nil {
		return
	}
	if childrenMap, childrenRps, err = s.secondReplies(ctx, sub, rootMap, req.Mid, _secondRepliesPn, _secondRepliesPs); err != nil {
		return
	}
	for _, rpID := range rpIDs {
		if r, ok := rootMap[rpID]; ok {
			if children, hasChild := childrenMap[rpID]; hasChild {
				r.Replies = children
				childrenRps = append(childrenRps, children...)
			} else {
				r.Replies = _emptyReplies
			}
			rootRps = append(rootRps, r)
		}
	}
	if rootRps != nil {
		res.Replies = rootRps
	} else {
		res.Replies = _emptyReplies
	}
	var rps []*reply.Reply
	rps = append(rps, rootRps...)
	rps = append(rps, childrenRps...)
	if err = s.buildReply(ctx, sub, rps, req.Mid, false); err != nil {
		return
	}
	res.Cursor = buildResCursor(&xmodel.Cursor{Prev: req.Cursor.Prev, Next: req.Cursor.Next, Ps: req.Cursor.Ps}, res.Replies, xmodel.CursorModeCursor)
	res.CommonRes = s.replyCommonRes(ctx, req.Mid, req.Oid, req.Type, sub)
	return
}

// RootFoldedReply ...
func (s *Service) RootFoldedReply(ctx context.Context, req *xmodel.RootFolderReq) (res *xmodel.RootFolderRes, err error) {
	var (
		rpIDs       []int64
		childrenMap map[int64]*reply.Reply
		childrenRps []*reply.Reply
		sub         *reply.Subject
	)
	res = new(xmodel.RootFolderRes)
	if req.Cursor.Backward() {
		return
	}
	cursor := &xmodel.Cursor{
		Ps:   req.Cursor.Ps,
		Next: req.Cursor.Next,
	}
	if sub, err = s.Subject(ctx, req.Oid, req.Type); err != nil {
		return
	}
	if rpIDs, err = s.foldedReplies(ctx, sub, req.Root, cursor); err != nil {
		return
	}
	if childrenMap, err = s.repliesMap(ctx, req.Oid, req.Type, rpIDs); err != nil {
		return
	}
	for _, rpID := range rpIDs {
		if r, ok := childrenMap[rpID]; ok {
			childrenRps = append(childrenRps, r)
		}
	}
	if childrenRps != nil {
		res.Replies = childrenRps
	} else {
		res.Replies = _emptyReplies
	}
	if err = s.buildReply(ctx, sub, res.Replies, req.Mid, false); err != nil {
		return
	}
	// 只有往下翻
	res.Cursor = buildResCursor(&xmodel.Cursor{Prev: req.Cursor.Prev, Next: req.Cursor.Next, Ps: req.Cursor.Ps}, res.Replies, xmodel.CursorModeCursor)
	res.CommonRes = s.replyCommonRes(ctx, req.Mid, req.Oid, req.Type, sub)
	return
}

// foldedReplies ...
func (s *Service) foldedReplies(ctx context.Context, sub *reply.Subject, root int64, cursor *xmodel.Cursor) (rpIDs []int64, err error) {
	var (
		kind string
		ID   int64
		ok   bool
	)
	if root == 0 {
		kind = xmodel.FolderKindSub
		ID = sub.Oid
	} else {
		kind = xmodel.FolderKindRoot
		ID = root
	}
	if ok, err = s.dao.Redis.ExpireFolder(ctx, kind, ID); err != nil {
		return
	}
	if ok {
		if rpIDs, err = s.dao.Redis.FolderByCursor(ctx, kind, ID, cursor); err != nil {
			return
		}
	} else {
		if rpIDs, err = s.dao.Reply.FoldedRepliesCursor(ctx, sub.Oid, sub.Type, root, cursor); err != nil {
			return
		}
		s.cache.Do(ctx, func(ctx context.Context) {
			s.dao.Databus.RecoverFolderIdx(ctx, sub.Oid, sub.Type, root)
		})
	}
	return
}
