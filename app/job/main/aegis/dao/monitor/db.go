package monitor

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go-common/app/job/main/aegis/model/monitor"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

const (
	_rulesByBid    = "SELECT id,type,bid,name,state,stime,etime,rule,uid,ctime,mtime FROM monitor_rule WHERE bid = ? AND state = 1 AND stime < ? AND etime > ?"
	_allValidRules = "SELECT id,type,bid,name,state,stime,etime,rule,uid,ctime,mtime FROM monitor_rule WHERE state = 1 AND stime < ? AND etime > ?"
)

// RulesByBid 获取某业务的监控
func (d *Dao) RulesByBid(c context.Context, bid int64) (rules []*monitor.Rule, err error) {
	var (
		rows *xsql.Rows
		now  = time.Now()
	)
	if rows, err = d.db.Query(c, _rulesByBid, bid, now, now); err != nil {
		log.Error("d.db.Exec error(%v)", errors.WithStack(err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		rule := &monitor.Rule{}
		var confStr string
		if err = rows.Scan(&rule.ID, &rule.Type, &rule.BID, &rule.Name, &rule.State, &rule.STime, &rule.ETime, &confStr, &rule.UID, &rule.CTime, &rule.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		conf := &monitor.RuleConf{}
		if err = json.Unmarshal([]byte(confStr), conf); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", confStr, err)
			return
		}
		rule.RuleConf = conf
		rules = append(rules, rule)
	}
	return
}

// ValidRules 获取有效的监控
func (d *Dao) ValidRules(c context.Context) (rules []*monitor.Rule, err error) {
	var (
		rows *xsql.Rows
		now  = time.Now()
	)
	if rows, err = d.db.Query(c, _allValidRules, now, now); err != nil {
		log.Error("d.db.Exec error(%v)", errors.WithStack(err))
		return
	}
	defer rows.Close()
	for rows.Next() {
		rule := &monitor.Rule{}
		var confStr string
		if err = rows.Scan(&rule.ID, &rule.Type, &rule.BID, &rule.Name, &rule.State, &rule.STime, &rule.ETime, &confStr, &rule.UID, &rule.CTime, &rule.MTime); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		conf := &monitor.RuleConf{}
		if err = json.Unmarshal([]byte(confStr), conf); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", confStr, err)
			return
		}
		rule.RuleConf = conf
		rules = append(rules, rule)
	}
	return
}
