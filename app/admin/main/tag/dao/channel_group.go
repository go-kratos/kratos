package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_setChannelSynonymString  = "(%d,%d,%q,%d,%d,%q,%q,%q)"
	_channelSynonymSQL        = "SELECT id,ptid,tid,alias,rank,state,operator,ctime,mtime FROM channel_group WHERE ptid=?"
	_upStateChannelSynonymSQL = "UPDATE channel_group SET state=? AND operator=%q WHERE ptid=?"
	_insertChannelSynonymSQL  = "REPLACE INTO channel_group(ptid,tid,alias,rank,state,operator,ctime,mtime) VALUES %s"
)

func setChannelSynonyms(synonyms []*model.ChannelSynonym) string {
	var sqls []string
	for _, v := range synonyms {
		sqls = append(sqls, fmt.Sprintf(_setChannelSynonymString, v.PTid, v.Tid, v.Alias, v.Rank, v.State, v.Operator, v.CTime.Time().Format(_timeLayoutFormat), v.MTime.Time().Format(_timeLayoutFormat)))
	}
	return strings.Join(sqls, ",")
}

// ChannelSynonymMap get channel synonyms map.
func (d *Dao) ChannelSynonymMap(c context.Context, ptid int64) (res map[int64]*model.ChannelSynonym, tids []int64, err error) {
	rows, err := d.db.Query(c, _channelSynonymSQL, ptid)
	if err != nil {
		log.Error("d.dao.ChannelSynonymMap(%v) error(%v)", ptid, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.ChannelSynonym)
	tids = make([]int64, 0)
	for rows.Next() {
		t := &model.ChannelSynonym{}
		if err = rows.Scan(&t.ID, &t.PTid, &t.Tid, &t.Alias, &t.Rank, &t.State, &t.Operator, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelSynonymMap scan(%v) error(%v)", ptid, err)
			return
		}
		res[t.Tid] = t
		tids = append(tids, t.Tid)
	}
	err = rows.Err()
	return
}

// TxUpStateChannelSynonym .
func (d *Dao) TxUpStateChannelSynonym(tx *sql.Tx, ptid int64, state int32, uname string) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_upStateChannelSynonymSQL, uname), state, ptid)
	if err != nil {
		log.Error("d.dao.UpStateChannelSynonym(%d,%d,%q) error(%v)", ptid, state, uname, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateChannelSynonyms tx update channel synonyms.
func (d *Dao) TxUpdateChannelSynonyms(tx *sql.Tx, synonyms []*model.ChannelSynonym) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_insertChannelSynonymSQL, setChannelSynonyms(synonyms)))
	if err != nil {
		log.Error("d.dao.TxUpdateChannelSynonyms(%+v) error(%v)", synonyms, err)
		return
	}
	return res.RowsAffected()
}
