package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/main/credit/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_addOpinionSQL    = "INSERT INTO blocked_opinion(vid,cid,content,attr,mid,vote,state) VALUE(?,?,?,?,?,?,?)"
	_addLikesSQL      = "UPDATE blocked_opinion SET likes=likes+1 where vid in(%s)"
	_addHatesSQL      = "UPDATE blocked_opinion SET hates=hates+1 where vid in(%s)"
	_delOpinionSQL    = "UPDATE blocked_opinion SET state = 1 WHERE vid = ?"
	_selOpinionIdxSQL = "SELECT vid,likes,hates,vote FROM blocked_opinion WHERE cid = ? AND state = 0 ORDER BY likes DESC"
	_selOpinionsSQL   = "SELECT vid,content,attr,hates,likes,mid,vote FROM blocked_opinion WHERE vid IN(%s) AND state = 0"
	_opContentMidSQL  = "SELECT content FROM blocked_opinion WHERE mid = ? AND state = 0 ORDER BY id DESC LIMIT 1"
)

// AddOpinionTx add user opinion by transaction.
func (d *Dao) AddOpinionTx(tx *sql.Tx, cid, opid, mid int64, content string, attr, vote, state int8) (affect int64, err error) {
	row, err := tx.Exec(_addOpinionSQL, opid, cid, content, attr, mid, vote, state)
	if err != nil {
		log.Error("add opinion err(%v)", err)
		return
	}
	return row.RowsAffected()
}

// AddLikes add user like opinions.
func (d *Dao) AddLikes(c context.Context, ids []int64) (affect int64, err error) {
	rows, err := d.db.Exec(c, fmt.Sprintf(_addLikesSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.AddLikes err(%v) ids(%v)", err, ids)
		return
	}
	return rows.RowsAffected()
}

// AddHates add user like opinions.
func (d *Dao) AddHates(c context.Context, ids []int64) (affect int64, err error) {
	rows, err := d.db.Exec(c, fmt.Sprintf(_addHatesSQL, xstr.JoinInts(ids)))
	if err != nil {
		log.Error("d.AddHates err(%v) ids(%v)", err, ids)
		return
	}
	return rows.RowsAffected()
}

// DelOpinion del opinion.
func (d *Dao) DelOpinion(c context.Context, opid int64) (err error) {
	if _, err = d.db.Exec(c, _delOpinionSQL, opid); err != nil {
		log.Error("d.Delopinion err(%v)", err)
	}
	return
}

// OpinionIdx get opinion index.
func (d *Dao) OpinionIdx(c context.Context, cid int64) (ops []*model.Opinion, err error) {
	rows, err := d.db.Query(c, _selOpinionIdxSQL, cid)
	if err != nil {
		log.Error("OpinionIdx err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		op := &model.Opinion{}
		if err = rows.Scan(&op.OpID, &op.Like, &op.Hate, &op.Vote); err != nil {
			log.Error("row.Scan err(%v)", err)
			return
		}
		ops = append(ops, op)
	}
	return
}

// OpinionCaseIdx get opinion case index.
func (d *Dao) OpinionCaseIdx(c context.Context, cid int64) (ops []*model.Opinion, err error) {
	rows, err := d.db.Query(c, _selOpinionIdxSQL, cid)
	if err != nil {
		log.Error("OpinionIdx err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		op := &model.Opinion{}
		if err = rows.Scan(&op.OpID, &op.Like, &op.Hate, &op.Vote); err != nil {
			log.Error("row.Scan err(%v)", err)
			return
		}
		if op.Like-op.Hate > -5 {
			ops = append(ops, op)
		}
	}
	return
}

// Opinions get opinions.
func (d *Dao) Opinions(c context.Context, opids []int64) (ops []*model.Opinion, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_selOpinionsSQL, xstr.JoinInts(opids)))
	if err != nil {
		log.Error("d.Opinions vids(%v) err(%v)", err, opids)
		return
	}
	defer rows.Close()
	for rows.Next() {
		op := &model.Opinion{}
		if err = rows.Scan(&op.OpID, &op.Content, &op.Attr, &op.Hate, &op.Like, &op.Mid, &op.Vote); err != nil {
			log.Error("row.Scan err(%v)", err)
		}
		ops = append(ops, op)
	}
	return
}

// OpContentMid get opinion content.
func (d *Dao) OpContentMid(c context.Context, mid int64) (content string, err error) {
	row := d.db.QueryRow(c, _opContentMidSQL, mid)
	if err = row.Scan(&content); err != nil {
		if err != sql.ErrNoRows {
			log.Error("row.Scan() error(%v)", err)
			return
		}
		content = ""
		err = nil
	}
	return
}
