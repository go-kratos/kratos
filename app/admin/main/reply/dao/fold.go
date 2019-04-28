package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/reply/model"
	"go-common/library/database/sql"
)

const (
	_foldedReplies      = "SELECT id,oid,type,mid,root,parent,dialog,count,rcount,`like`,floor,state,attr,ctime,mtime FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12"
	_countFoldedReplies = "SELECT COUNT(*) FROM reply_%d WHERE oid=? AND type=? AND root=? AND state=12"
)

// TxCountFoldedReplies ...
func (d *Dao) TxCountFoldedReplies(tx *sql.Tx, oid int64, tp int32, root int64) (count int, err error) {
	if err = tx.QueryRow(fmt.Sprintf(_countFoldedReplies, hit(oid)), oid, tp, root).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}
	return
}

// FoldedReplies ...
func (d *Dao) FoldedReplies(ctx context.Context, oid int64, tp int32, root int64) (rps []*model.Reply, err error) {
	rows, err := d.dbSlave.Query(ctx, fmt.Sprintf(_foldedReplies, hit(oid)), oid, tp, root)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Reply)
		if err = rows.Scan(&r.ID, &r.Oid, &r.Type, &r.Mid, &r.Root, &r.Parent, &r.Dialog, &r.Count, &r.RCount, &r.Like, &r.Floor, &r.State, &r.Attr, &r.CTime, &r.MTime); err != nil {
			return
		}
		rps = append(rps, r)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// RemRdsByFold ...
func (d *Dao) RemRdsByFold(ctx context.Context, roots []int64, childMap map[int64][]int64, sub *model.Subject, rpMap map[int64]*model.Reply) {
	var (
		keyMap = make(map[string][]int64)
	)
	// 评论列表缓存
	keyMap[keyMainIdx(sub.Oid, sub.Type, model.SortByFloor)] = roots
	keyMap[keyMainIdx(sub.Oid, sub.Type, model.SortByCount)] = roots
	keyMap[keyMainIdx(sub.Oid, sub.Type, model.SortByLike)] = roots
	for root, children := range childMap {
		// 评论详情页缓存
		keyMap[keyRootIdx(root)] = children
		for _, child := range children {
			// 对话列表的缓存
			if rp, ok := rpMap[child]; ok && rp.Dialog != 0 {
				keyMap[keyDialogIdx(rp.Dialog)] = append(keyMap[keyDialogIdx(rp.Dialog)], rp.ID)
			}
		}
	}
	d.RemReplyFromRedis(ctx, keyMap)
}

// AddRdsByFold ...
func (d *Dao) AddRdsByFold(ctx context.Context, roots []int64, childMap map[int64][]int64, sub *model.Subject, rpMap map[int64]*model.Reply) {
	var (
		ok         bool
		err        error
		keyMapping = make(map[string][]*model.Reply)
	)
	if ok, err = d.ExpireFolder(ctx, model.FolderKindSub, sub.Oid); err != nil {
		return
	}
	if ok {
		key := keyFolderIdx(model.FolderKindSub, sub.Oid)
		for _, root := range roots {
			if rp, ok := rpMap[root]; ok {
				keyMapping[key] = append(keyMapping[key], rp)
			}
		}
	}
	// 这里不回源
	for root, children := range childMap {
		if ok, err = d.ExpireFolder(ctx, model.FolderKindRoot, root); err != nil {
			return
		}
		if ok {
			key := keyFolderIdx(model.FolderKindRoot, root)
			for _, child := range children {
				if rp, ok := rpMap[child]; ok {
					keyMapping[key] = append(keyMapping[key], rp)
				}
			}
		}
	}
	d.AddFolder(ctx, keyMapping)
}
