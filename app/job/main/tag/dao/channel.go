package dao

import (
	"context"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
)

const (
	_channelSQL     = "SELECT tid FROM channel WHERE state in (2,3)"
	_channelRuleSQL = "SELECT id,tid,in_rule,notin_rule FROM channel_rule WHERE id>? AND state=0 ORDER BY id ASC LIMIT 1000"
)

// ChannelMap .
func (d *Dao) ChannelMap(c context.Context) (tidMap map[int64]struct{}, err error) {
	tidMap = make(map[int64]struct{})
	rows, err := d.platform.Query(c, _channelSQL)
	if err != nil {
		log.Error("d.Channels() error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tid int64
		if err = rows.Scan(&tid); err != nil {
			log.Error("d.Channels() scan() error(%v)", err)
			return
		}
		tidMap[tid] = struct{}{}
	}
	err = rows.Err()
	return
}

// ChannelRules .
func (d *Dao) ChannelRules(c context.Context, lastID int64) (rules []*model.ChannelRule, err error) {
	rules = make([]*model.ChannelRule, 0, 1000)
	rows, err := d.platform.Query(c, _channelRuleSQL, lastID)
	if err != nil {
		log.Error("d.ChannelRules(%d) error(%v)", lastID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		cr := &model.ChannelRule{}
		if err = rows.Scan(&cr.ID, &cr.Tid, &cr.InRule, &cr.NotInRule); err != nil {
			log.Error("d.ChannelRules(%d) scan() error(%v)", lastID, err)
			return
		}
		rules = append(rules, cr)
	}
	err = rows.Err()
	return
}
