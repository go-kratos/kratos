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
	_setChannelRuleString    = "(%d,%q,%q,%q,%d,%q)"
	_channelRuleSQL          = "SELECT id,tid,in_rule,notin_rule,editor,state,ctime,mtime FROM channel_rule WHERE tid=?"
	_mutiChannelRuleSQL      = "SELECT id,tid,in_rule,notin_rule,editor,state,ctime,mtime FROM channel_rule WHERE tid in (%s)"
	_updateStateRuleByTidSQL = "UPDATE channel_rule SET state=?,editor=? WHERE tid=?"
	_updateChannelRuleSQL    = "REPLACE INTO channel_rule(tid,in_rule,notin_rule,editor,state,ctime) VALUES %s"
	_timeLayoutFormat        = "2006-01-02 15:04:05"
)

func setChannelRules(channelRules []*model.ChannelRule) string {
	var sqls []string
	for _, v := range channelRules {
		sqls = append(sqls, fmt.Sprintf(_setChannelRuleString, v.Tid, v.InRule, v.NotInRule, v.Editor, v.State, v.CTime.Time().Format(_timeLayoutFormat)))
	}
	return strings.Join(sqls, ",")
}

// ChannelRule get channel rule by channel id.
func (d *Dao) ChannelRule(c context.Context, tid int64) (res []*model.ChannelRule, tids []int64, err error) {
	rows, err := d.db.Query(c, _channelRuleSQL, tid)
	if err != nil {
		log.Error("d.dao.ChannelRule(%d) error(%v)", tid, err)
		return
	}
	defer rows.Close()
	res = make([]*model.ChannelRule, 0)
	for rows.Next() {
		t := &model.ChannelRule{}
		var tempTids []int64
		if err = rows.Scan(&t.ID, &t.Tid, &t.InRule, &t.NotInRule, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelRule(%d) scan() error(%v)", tid, err)
			return
		}
		res = append(res, t)
		if tempTids, err = xstr.SplitInts(t.InRule); err != nil {
			log.Warn("d.dao.ChannelRule(%d) xstr.SplitInts(%s)", tid, t.InRule)
			err = nil
		}
		tids = append(tids, tempTids...)
		if tempTids, err = xstr.SplitInts(t.NotInRule); err != nil {
			log.Warn("d.dao.ChannelRule(%d) xstr.SplitInts(%s)", tid, t.NotInRule)
			err = nil
		}
		tids = append(tids, tempTids...)
	}
	return
}

// ChannelRuleMap get channel rule map by channel id.
func (d *Dao) ChannelRuleMap(c context.Context, tid int64) (res map[int64]*model.ChannelRule, err error) {
	rows, err := d.db.Query(c, _channelRuleSQL, tid)
	if err != nil {
		log.Error("d.dao.ChannelRule(%d) error(%v)", tid, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.ChannelRule)
	for rows.Next() {
		t := &model.ChannelRule{}
		if err = rows.Scan(&t.ID, &t.Tid, &t.InRule, &t.NotInRule, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelRule(%d) scan() error(%v)", tid, err)
			return
		}
		res[t.ID] = t
	}
	return
}

// ChannelRules get channel rules by channel ids.
func (d *Dao) ChannelRules(c context.Context, tids []int64) (res map[int64][]*model.ChannelRule, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_mutiChannelRuleSQL, xstr.JoinInts(tids)))
	if err != nil {
		log.Error("d.dao.ChannelRule(%v) error(%v)", tids, err)
		return
	}
	defer rows.Close()
	res = make(map[int64][]*model.ChannelRule, len(tids))
	for rows.Next() {
		t := &model.ChannelRule{}
		if err = rows.Scan(&t.ID, &t.Tid, &t.InRule, &t.NotInRule, &t.Editor, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelRule(%v) scan() error(%v)", tids, err)
			return
		}
		res[t.Tid] = append(res[t.Tid], t)
	}
	return
}

// TxUpChannelRuleState tx channel state.
func (d *Dao) TxUpChannelRuleState(tx *sql.Tx, tid int64, state int32, uname string) (affect int64, err error) {
	res, err := tx.Exec(_updateStateRuleByTidSQL, state, uname, tid)
	if err != nil {
		log.Error("d.dao.TxUpChannelRuleState(%d,%s,%d) error(%v)", tid, uname, state, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateChannelRules tx update channel rules.
func (d *Dao) TxUpdateChannelRules(tx *sql.Tx, channelRules []*model.ChannelRule) (affect int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateChannelRuleSQL, setChannelRules(channelRules)))
	if err != nil {
		log.Error("d.dao.TxUpdateChannelRules(%+v) error(%v)", channelRules, err)
		return
	}
	return res.RowsAffected()
}
