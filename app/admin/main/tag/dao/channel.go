package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_setChannelString         = "(%d,%d,%d,%d,%d,%q,%q,%d,%q)"
	_countChannelByTypeSQL    = "SELECT count(tid) FROM channel WHERE type=?"
	_countRecommendChannelSQL = "SELECT count(tid) FROM channel WHERE state=3"
	_countChanneListSQL       = "SELECT count(tid) FROM channel %s"
	_channelSQL               = "SELECT tid,rank,top_rank,type,attr,operator,editor,state,ctime,mtime FROM channel where tid=?"
	_channelMutiSQL           = "SELECT tid,rank,top_rank,type,attr,operator,editor,state,ctime,mtime FROM channel where tid in (%s)"
	_channeListSQL            = "SELECT tid,rank,top_rank,type,operator,editor,attr,state,ctime,mtime FROM channel %s ORDER BY %s %s LIMIT ?,?"
	_channelAllSQL            = "SELECT tid,rank,top_rank,type,attr,operator,editor,state,ctime,mtime FROM channel"
	_channelsByTypeSQL        = "SELECT tid,rank,top_rank,type,attr,operator,editor,state,ctime,mtime FROM channel where type=?"
	_updateChannelSQL         = "UPDATE channel SET rank=?,top_rank=?,type=?,attr=?,editor=?,state=? WHERE tid=?"
	_setMutiChannelSQL        = "REPLACE INTO channel (tid,rank,top_rank,type,attr,operator,editor,state,ctime) VALUES %s"
	_insertChannelSQL         = "INSERT INTO channel (tid,rank,top_rank,type,attr,operator,editor,state) VALUES (?,?,?,?,?,?,?,?)"
	_recommandChannelSQL      = "SELECT tid,rank,top_rank,type,attr,operator,editor,state,ctime,mtime FROM channel WHERE state=3"
	_updateChannelAttrSQL     = "UPDATE channel SET attr=?,operator=? WHERE tid=?"
)

func setChannels(channels []*model.Channel) string {
	var sqls []string
	for _, v := range channels {
		sqls = append(sqls, fmt.Sprintf(_setChannelString, v.ID, v.Rank, v.TopRank, v.Type, v.Attr, v.Operator, v.Editor, v.State, v.CTime.Time().Format(_timeLayoutFormat)))
	}
	return strings.Join(sqls, ",")
}

// Channel get a channel info by channel tid.
func (d *Dao) Channel(c context.Context, tid int64) (res *model.Channel, err error) {
	res = new(model.Channel)
	row := d.db.QueryRow(c, _channelSQL, tid)
	if err = row.Scan(&res.ID, &res.Rank, &res.TopRank, &res.Type, &res.Attr, &res.Operator, &res.Editor, &res.State, &res.CTime, &res.MTime); err != nil {
		log.Error("d.dao.Channel(%d) error(%v)", tid, err)
		if err == sql.ErrNoRows {
			err = nil
			res = nil
		}
	}
	return
}

// Channels get muti channels by tids.
func (d *Dao) Channels(c context.Context, tids []int64) (res []*model.Channel, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_channelMutiSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("d.dao.Channels(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Channel, 0, len(tids))
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Attr, &t.Operator, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.Channels scan(%v) error(%v)", tids, err)
			return
		}
		res = append(res, t)
	}
	return
}

// ChannelMap get channel map by tids.
func (d *Dao) ChannelMap(c context.Context, tids []int64) (res map[int64]*model.Channel, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_channelMutiSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("d.dao.ChannelMap(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Channel, len(tids))
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Attr, &t.Operator, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelMap scan(%v) error(%v)", tids, err)
			return
		}
		res[t.ID] = t
	}
	return
}

// ChannelAll get all channels.
func (d *Dao) ChannelAll(c context.Context) (res []*model.Channel, tids []int64, err error) {
	rows, err := d.db.Query(c, _channelAllSQL)
	if err != nil {
		log.Error("d.dao.ChannelAll() error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*model.Channel, 0)
	tids = make([]int64, 0)
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Attr, &t.Operator, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelAll scan() error(%v)", err)
			return
		}
		res = append(res, t)
		tids = append(tids, t.ID)
	}
	return
}

// ChannelsByType get muti channels by type.
func (d *Dao) ChannelsByType(c context.Context, tp int64) (res []*model.Channel, tids []int64, err error) {
	rows, err := d.db.Query(c, _channelsByTypeSQL, tp)
	if err != nil {
		log.Error("d.dao.ChannelsByType(%d) error(%v)", tp, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Channel, 0)
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Attr, &t.Operator, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.Channels scan(%d) error(%v)", tp, err)
			return
		}
		res = append(res, t)
		tids = append(tids, t.ID)
	}
	return
}

// ChanneList get channels by sqls.
func (d *Dao) ChanneList(c context.Context, sqls []string, order, sort string, start, end int32) (res []*model.Channel, ids []int64, err error) {
	var (
		sql string
	)
	if len(sqls) > 0 {
		sql = fmt.Sprintf("WHERE %s", strings.Join(sqls, " AND "))
	}
	rows, err := d.db.Query(c, fmt.Sprintf(_channeListSQL, sql, order, sort), start, end)
	if err != nil {
		log.Error("d.dao.ChanneList(%s,%s,%s,%d,%d) error(%v)", sql, order, sort, start, end, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Operator, &t.Editor, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChanneList Scan() error(%v)", err)
			return
		}
		res = append(res, t)
		ids = append(ids, t.ID)
	}
	return
}

// CountChanneList get channels by sql.
func (d *Dao) CountChanneList(c context.Context, sqls []string) (count int32, err error) {
	var sqlStr string
	if len(sqls) > 0 {
		sqlStr = fmt.Sprintf("WHERE %s", strings.Join(sqls, " AND "))
	}
	row := d.db.QueryRow(c, fmt.Sprintf(_countChanneListSQL, sqlStr))
	if err = row.Scan(&count); err != nil {
		log.Error("d.dao.CountChanneList(%s) error(%v)", sqlStr, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// CountChannelByType get channels by type.
func (d *Dao) CountChannelByType(c context.Context, tp int64) (count int32, err error) {
	row := d.db.QueryRow(c, _countChannelByTypeSQL, tp)
	if err = row.Scan(&count); err != nil {
		log.Error("d.dao.CountChannelByType(%d) error(%v)", tp, err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// CountRecommendChannel get count recommend channel.
func (d *Dao) CountRecommendChannel(c context.Context) (count int32, err error) {
	row := d.db.QueryRow(c, _countRecommendChannelSQL)
	if err = row.Scan(&count); err != nil {
		log.Error("d.dao.CountRecommendChannel() error(%v)", err)
		if err == sql.ErrNoRows {
			err = nil
		}
	}
	return
}

// TxInsertChannel tx insert channel.
func (d *Dao) TxInsertChannel(tx *sql.Tx, channel *model.Channel) (id int64, err error) {
	res, err := tx.Exec(_insertChannelSQL, channel.ID, channel.Rank, channel.TopRank, channel.Type, channel.Attr, channel.Operator, channel.Editor, channel.State)
	if err != nil {
		log.Error("d.dao.InsertChannel(%v) error(%v)", channel, err)
		return
	}
	return res.LastInsertId()
}

// UpdateChannel update channel.
func (d *Dao) UpdateChannel(c context.Context, channel *model.Channel) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateChannelSQL, channel.Rank, channel.TopRank, channel.Type, channel.Attr, channel.Editor, channel.State, channel.ID)
	if err != nil {
		log.Error("d.dao.UpdateChannel(%v) error(%v)", channel, err)
		return
	}
	return res.RowsAffected()
}

// TxUpChannel tx channel.
func (d *Dao) TxUpChannel(tx *sql.Tx, channel *model.Channel) (affect int64, err error) {
	res, err := tx.Exec(_updateChannelSQL, channel.Rank, channel.TopRank, channel.Type, channel.Attr, channel.Editor, channel.State, channel.ID)
	if err != nil {
		log.Error("d.dao.TxUpChannel(%v) error(%v)", channel, err)
		return
	}
	return res.RowsAffected()
}

// UpdateChannels update channels.
func (d *Dao) UpdateChannels(c context.Context, channels []*model.Channel) (affect int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_setMutiChannelSQL, setChannels(channels)))
	if err != nil {
		log.Error("d.dao.UpdateChannelCategories(%+v) error(%v)", channels, err)
		return
	}
	return res.RowsAffected()
}

// RecommandChannel recommand channels.
func (d *Dao) RecommandChannel(c context.Context) (res []*model.Channel, tids []int64, err error) {
	rows, err := d.db.Query(c, _recommandChannelSQL)
	if err != nil {
		log.Error("d.dao.Channels(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Channel, 0)
	tids = make([]int64, 0)
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Attr, &t.Operator, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.Channels scan(%v) error(%v)", tids, err)
			return
		}
		res = append(res, t)
		tids = append(tids, t.ID)
	}
	return
}

// TxUpChannelAttr tx update channel attr.
func (d *Dao) TxUpChannelAttr(tx *sql.Tx, channel *model.Channel) (affect int64, err error) {
	res, err := tx.Exec(_updateChannelAttrSQL, channel.Attr, channel.Operator, channel.ID)
	if err != nil {
		log.Error("d.TxUpChannelAttr(%v) error(%v)", channel, err)
		return
	}
	return res.RowsAffected()
}
