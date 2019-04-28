package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/reply/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"

	"go-common/library/sync/errgroup.v2"
)

func compose(oids, tps, rpIDs []int64) (rpMap map[int64][]int64, tpMap map[int64]int64) {
	if len(oids) != len(rpIDs) {
		return
	}
	rpMap = make(map[int64][]int64)
	tpMap = make(map[int64]int64)
	for i, oid := range oids {
		if _, ok := rpMap[oid]; ok {
			rpMap[oid] = append(rpMap[oid], rpIDs[i])
		} else {
			rpMap[oid] = []int64{rpIDs[i]}
		}
		tpMap[oid] = tps[i]
	}
	return
}

func extend(oid int64, length int) []int64 {
	oids := make([]int64, 0, length)
	for i := 0; i < length; i++ {
		oids = append(oids, oid)
	}
	return oids
}

// FoldReplies ...
func (s *Service) FoldReplies(ctx context.Context, oids, tps, rpIDs []int64) (err error) {
	g := errgroup.WithContext(ctx)
	g.GOMAXPROCS(2)
	rpMap, tpMap := compose(oids, tps, rpIDs)
	for oid, IDs := range rpMap {
		if tp, ok := tpMap[oid]; ok {
			oid, tp, IDs := oid, tp, IDs
			g.Go(func(ctx context.Context) error {
				return s.foldReplies(ctx, oid, tp, IDs)
			})
		}
	}
	return g.Wait()
}

func (s *Service) foldReplies(ctx context.Context, oid, tp int64, rpIDs []int64) (err error) {
	var (
		sub *model.Subject
		// 所有的评论，包含需要被折叠的子评论的根评论
		rpMap = make(map[int64]*model.Reply)
		// 需要被折叠的子评论的根评论IDs
		roots   []int64
		rootMap map[int64]*model.Reply
		// 应该被折叠的评论map
		rps = make(map[int64]*model.Reply)
		mu  sync.Mutex
	)
	if sub, err = s.subject(ctx, oid, int32(tp)); err != nil {
		return
	}
	if rpMap, err = s.dao.Replies(ctx, extend(oid, len(rpIDs)), rpIDs); err != nil {
		return
	}
	for _, rp := range rpMap {
		if !rp.IsRoot() {
			roots = append(roots, rp.Root)
		}
		rps[rp.ID] = rp
	}
	if len(roots) > 0 {
		if rootMap, err = s.dao.Replies(ctx, extend(oid, len(roots)), roots); err != nil {
			return
		}
		for _, root := range rootMap {
			rpMap[root.ID] = root
		}
	}
	g := errgroup.WithContext(ctx)
	g.GOMAXPROCS(4)
	for rpID := range rps {
		rpID := rpID
		g.Go(func(ctx context.Context) error {
			if err := s.tranFoldReply(ctx, sub.Oid, rpID); err != nil {
				mu.Lock()
				delete(rps, rpID)
				mu.Unlock()
				return err
			}
			return nil
		})
	}
	err = g.Wait()
	for _, rp := range rps {
		// 这里不是删除，是为了让reply-feed去掉热评, 折叠评论是可以有互动行为的， 所以这里异步，丢了也无所谓
		rp := rp
		s.cache.Do(ctx, func(ctx context.Context) {
			s.pubEvent(ctx, "reply_del", 0, sub, rp, nil)
		})
	}
	// 标记数据库有折叠评论, rpMap 是所有的被折叠评论以及他们的根评论
	s.markHasFolded(ctx, rps, rpMap, sub)
	s.handleCacheByFold(ctx, rps, rpMap, sub)
	return err
}

// handleCacheByFold ...
func (s *Service) handleCacheByFold(ctx context.Context, rps map[int64]*model.Reply, rpMap map[int64]*model.Reply, sub *model.Subject) {
	var (
		roots    []int64
		childMap = make(map[int64][]int64)
	)
	for _, rp := range rps {
		if rp.IsRoot() {
			roots = append(roots, rp.ID)
		} else {
			childMap[rp.Root] = append(childMap[rp.Root], rp.ID)
		}
	}
	s.cacheOperater.Do(ctx, func(ctx context.Context) {
		// 删正常列表的redis缓存
		s.dao.RemRdsByFold(ctx, roots, childMap, sub, rpMap)
		// 加折叠列表redis缓存
		s.dao.AddRdsByFold(ctx, roots, childMap, sub, rpMap)
	})
}

// markAsFolded ...
func (s *Service) markHasFolded(ctx context.Context, foldedRp map[int64]*model.Reply, rpMap map[int64]*model.Reply, sub *model.Subject) {
	var (
		dirtyCacheRpIDs []int64
		markedRpIDs     []int64
	)
	for _, rp := range foldedRp {
		rp := rp
		if rp.IsRoot() {
			// 如果subject还没被标记过
			if !sub.HasFolded() {
				sub.MarkHasFolded()
				// 修改数据库标记
				s.marker.Do(ctx, func(ctx context.Context) {
					if err := s.tranMarkSubHasFolded(ctx, sub.Oid, sub.Type); err != nil {
						return
					}
				})
				// 删掉 subject mc
				s.cacheOperater.Do(ctx, func(ctx context.Context) {
					s.dao.DelSubjectCache(ctx, sub.Oid, sub.Type)
				})
			}
		} else {
			if _, ok := rpMap[rp.Root]; ok {
				markedRpIDs = append(markedRpIDs, rp.Root)
			}
		}
		dirtyCacheRpIDs = append(dirtyCacheRpIDs, rp.ID)
	}
	for _, rpID := range markedRpIDs {
		// 修改数据库标记
		rpID := rpID
		s.marker.Do(ctx, func(ctx context.Context) {
			if err := s.tranMarkReplyHasFolded(ctx, sub.Oid, rpID); err != nil {
				return
			}
		})
	}
	dirtyCacheRpIDs = append(dirtyCacheRpIDs, markedRpIDs...)
	for _, rpID := range dirtyCacheRpIDs {
		// 删除被折叠子评论的根评论以及被折叠的子评论和根评论
		rpID := rpID
		s.cacheOperater.Do(ctx, func(ctx context.Context) {
			s.dao.DelReplyCache(ctx, rpID)
		})
	}
}

func (s *Service) tranFoldReply(ctx context.Context, oid, rpID int64) (err error) {
	var (
		tx *sql.Tx
		rp *model.Reply
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if rp, err = s.dao.TxReplyForUpdate(tx, oid, rpID); err != nil {
		tx.Rollback()
		return
	}
	if rp.DenyFolded() {
		tx.Rollback()
		return ecode.ReplyForbidFolded
	}
	if _, err = s.dao.TxUpdateReplyState(tx, oid, rpID, model.StateFolded, time.Now()); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

func (s *Service) tranMarkReplyHasFolded(ctx context.Context, oid, rpID int64) (err error) {
	var (
		tx *sql.Tx
		rp *model.Reply
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if rp, err = s.dao.TxReplyForUpdate(tx, oid, rpID); err != nil {
		tx.Rollback()
		return
	}
	rp.MarkHasFolded()
	if _, err = s.dao.TxUpReplyAttr(tx, oid, rpID, rp.Attr, time.Now()); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

func (s *Service) tranMarkSubHasFolded(ctx context.Context, oid int64, tp int32) (err error) {
	var (
		tx  *sql.Tx
		sub *model.Subject
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if sub, err = s.dao.TxSubjectForUpdate(tx, oid, tp); err != nil {
		tx.Rollback()
		return
	}
	sub.MarkHasFolded()
	if _, err = s.dao.TxUpSubAttr(tx, oid, tp, sub.Attr, time.Now()); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

// handleFolded 处理折叠评论的逻辑，包括折叠评论被删除等， 状态改变之后标记改变的问题...
func (s *Service) handleFolded(ctx context.Context, rp *model.Reply) {
	sub, root, err := s.handleHasFoldedMark(ctx, rp.Oid, rp.Type, rp.Root)
	if err != nil {
		return
	}
	if sub != nil {
		s.dao.DelSubjectCache(ctx, sub.Oid, sub.Type)
	}
	if root != nil {
		s.dao.DelReplyCache(ctx, root.ID)
	}
	s.remFoldedCache(ctx, rp)
}

func (s *Service) handleHasFoldedMark(ctx context.Context, oid int64, tp int32, root int64) (sub *model.Subject, reply *model.Reply, err error) {
	var (
		tx    *sql.Tx
		count int
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	// 锁subject表
	if sub, err = s.dao.TxSubjectForUpdate(tx, oid, tp); err != nil {
		tx.Rollback()
		return
	}
	if count, err = s.dao.TxCountFoldedReplies(tx, oid, tp, root); err != nil || count > 0 {
		tx.Rollback()
		return
	}
	// 折叠根评论
	if root == 0 {
		if !sub.HasFolded() {
			tx.Rollback()
			return
		}
		sub.UnmarkHasFolded()
		if _, err = s.dao.TxUpSubAttr(tx, oid, tp, sub.Attr, time.Now()); err != nil {
			tx.Rollback()
			return
		}
	} else {
		if reply, err = s.dao.TxReplyForUpdate(tx, oid, root); err != nil {
			tx.Rollback()
			return
		}
		if !reply.HasFolded() {
			tx.Rollback()
			return
		}
		reply.UnmarkHasFolded()
		if _, err = s.dao.TxUpReplyAttr(tx, oid, root, reply.Attr, time.Now()); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	return
}

// remFoldedCache ...
func (s *Service) remFoldedCache(ctx context.Context, rp *model.Reply) {
	if rp.IsRoot() {
		s.dao.RemFolder(ctx, model.FolderKindSub, rp.Oid, rp.ID)
	} else {
		s.dao.RemFolder(ctx, model.FolderKindRoot, rp.Root, rp.ID)
	}
}
