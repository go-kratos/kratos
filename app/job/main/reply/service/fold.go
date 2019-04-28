package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

func (s *Service) folderHanlder(ctx context.Context, msg *consumerMsg) {
	var d struct {
		Op   string `json:"op"`
		Oid  int64  `json:"oid"`
		Tp   int8   `json:"tp"`
		Root int64  `json:"root"`
	}
	if err := json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	switch d.Op {
	case "re_idx":
		s.recoverFolderIdx(ctx, d.Oid, d.Tp, d.Root)
	// case "marker":
	// 删除折叠评论后需要查询是否需要取消标记
	// s.marker.Do(ctx, func(ctx context.Context) {
	// 	s.handleHasFoldedMark(ctx, d.Oid, d.Tp, d.Root)
	// })
	default:
		return
	}
}

// handleFolded ...
func (s *Service) handleFolded(ctx context.Context, rp *reply.Reply) {
	sub, root, err := s.handleHasFoldedMark(ctx, rp.Oid, rp.Type, rp.Root)
	if err != nil {
		return
	}
	if sub != nil {
		s.dao.Mc.DeleteSub(ctx, sub.Oid, sub.Type)
	}
	if root != nil {
		s.dao.Mc.DeleteReply(ctx, root.RpID)
	}
	s.remFoldedCache(ctx, rp)
}

func (s *Service) recoverFolderIdx(ctx context.Context, oid int64, tp int8, root int64) {
	rps, err := s.dao.Reply.FoldedReplies(ctx, oid, tp, root)
	if err != nil || len(rps) == 0 {
		return
	}
	// 折叠根评论
	if root == 0 {
		s.dao.Redis.AddFolderBatch(ctx, reply.FolderKindSub, oid, rps)
	} else {
		s.dao.Redis.AddFolderBatch(ctx, reply.FolderKindRoot, root, rps)
	}
}

func (s *Service) handleHasFoldedMark(ctx context.Context, oid int64, tp int8, root int64) (sub *reply.Subject, reply *reply.Reply, err error) {
	var (
		tx    *sql.Tx
		count int
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	// 锁subject表
	if sub, err = s.dao.Subject.GetForUpdate(tx, oid, tp); err != nil {
		tx.Rollback()
		return
	}
	if count, err = s.dao.Reply.TxCountFoldedReplies(tx, oid, tp, root); err != nil || count > 0 {
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
		if _, err = s.dao.Subject.TxUpAttr(tx, oid, tp, sub.Attr, time.Now()); err != nil {
			tx.Rollback()
			return
		}
	} else {
		if reply, err = s.dao.Reply.GetForUpdate(tx, oid, root); err != nil {
			tx.Rollback()
			return
		}
		if !reply.HasFolded() {
			tx.Rollback()
			return
		}
		reply.UnmarkHasFolded()
		if _, err = s.dao.Reply.TxUpAttr(tx, oid, root, reply.Attr, time.Now()); err != nil {
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
func (s *Service) remFoldedCache(ctx context.Context, rp *reply.Reply) {
	if rp.IsRoot() {
		s.dao.Redis.RemFolder(ctx, reply.FolderKindSub, rp.Oid, rp.RpID)
	} else {
		s.dao.Redis.RemFolder(ctx, reply.FolderKindRoot, rp.Root, rp.RpID)
	}
}
