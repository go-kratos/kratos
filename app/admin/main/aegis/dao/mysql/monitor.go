package mysql

import (
	"context"
	"encoding/json"
	"go-common/app/admin/main/aegis/model/monitor"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_moniBizRulesSql = "SELECT id,type,bid,name,state,stime,etime,rule,uid,ctime,mtime FROM monitor_rule WHERE bid = ?"
	_moniRuleSql     = "SELECT id,type,bid,name,state,stime,etime,rule,uid,ctime,mtime FROM monitor_rule WHERE id = ?"
)

// MoniBizRules 获取监控业务的所有配置
func (d *Dao) MoniBizRules(c context.Context, bid int64) (rules []*monitor.Rule, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _moniBizRulesSql, bid); err != nil {
		log.Error("db.Query() error(%v)", err)
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

// MoniRule 根据id获取监控规则
func (d *Dao) MoniRule(c context.Context, rid int64) (rule *monitor.Rule, err error) {
	rule = &monitor.Rule{}
	var confStr string
	row := d.db.QueryRow(c, _moniRuleSql, rid)
	if err = row.Scan(&rule.ID, &rule.Type, &rule.BID, &rule.Name, &rule.State, &rule.STime, &rule.ETime, &confStr, &rule.UID, &rule.CTime, &rule.MTime); err != nil {
		rule = nil
		log.Error("row.Scan error(%v)", err)
		return
	}
	conf := &monitor.RuleConf{}
	if err = json.Unmarshal([]byte(confStr), conf); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", confStr, err)
		return
	}
	rule.RuleConf = conf
	return
}
