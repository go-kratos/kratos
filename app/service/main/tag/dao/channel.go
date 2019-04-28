package dao

import (
	"context"

	"go-common/app/service/main/tag/model"
	"go-common/library/log"
)

const (
	_channelCategorieSQL = "SELECT id,name,`order`,attr,state,ctime,mtime FROM channel_type WHERE state=? AND id>? ORDER BY id ASC LIMIT ?"
	_channelsSQL         = "SELECT tid,rank,top_rank,type,operator,state,attr,ctime,mtime FROM channel WHERE tid>? ORDER BY tid ASC LIMIT ?"
	_channelRulesSQL     = "SELECT id,tid,in_rule,notin_rule,state FROM channel_rule WHERE state=? AND id>? ORDER BY id ASC LIMIT ?"
	_channelGroupSQL     = "SELECT id,ptid,tid,`alias`,rank,ctime,mtime FROM channel_group WHERE ptid=? AND state=0"
)

// ChannelCategories get channel categopries.
func (d *Dao) ChannelCategories(c context.Context, arg *model.ArgChannelCategory) (res []*model.ChannelCategory, err error) {
	rows, err := d.db.Query(c, _channelCategorieSQL, arg.State, arg.LastID, arg.Size)
	if err != nil {
		log.Error("d.dao.ChannelCategories(%v) error(%v)", arg, err)
		return
	}
	defer rows.Close()
	res = make([]*model.ChannelCategory, 0, arg.Size)
	for rows.Next() {
		t := &model.ChannelCategory{}
		if err = rows.Scan(&t.ID, &t.Name, &t.Order, &t.Attr, &t.State, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelCategories(%v) rows.Scan() error(%v)", arg, err)
			return
		}
		res = append(res, t)
	}
	return
}

// Channels get channel listby tids.
func (d *Dao) Channels(c context.Context, arg *model.ArgChannels) (res []*model.Channel, err error) {
	rows, err := d.db.Query(c, _channelsSQL, arg.LastID, arg.Size)
	if err != nil {
		log.Error("d.dao.Channels(%v) error(%v)", arg, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Channel, 0, arg.Size)
	for rows.Next() {
		t := &model.Channel{}
		if err = rows.Scan(&t.ID, &t.Rank, &t.TopRank, &t.Type, &t.Operator, &t.State, &t.Attr, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.Channels(%v) rows.Scan() error(%v)", arg, err)
			return
		}
		res = append(res, t)
	}
	return
}

// ChannelRule get channel rules .
func (d *Dao) ChannelRule(c context.Context, arg *model.ArgChannelRule) (res []*model.ChannelRule, err error) {
	rows, err := d.db.Query(c, _channelRulesSQL, arg.State, arg.LastID, arg.Size)
	if err != nil {
		log.Error("d.dao.ChannelRule(%v) error(%v)", arg, err)
		return
	}
	defer rows.Close()
	res = make([]*model.ChannelRule, 0, arg.Size)
	for rows.Next() {
		t := &model.ChannelRule{}
		if err = rows.Scan(&t.Id, &t.Tid, &t.InRule, &t.NotinRule, &t.State); err != nil {
			log.Error("d.dao.ChannelRule(%v) rows.Scan() error(%v)", arg, err)
			return
		}
		res = append(res, t)
	}
	return
}

// ChannelGroup channel group.
func (d *Dao) ChannelGroup(c context.Context, tid int64) (res []*model.ChannelGroup, err error) {
	res = make([]*model.ChannelGroup, 0, model.ChannelMaxGroups)
	rows, err := d.db.Query(c, _channelGroupSQL, tid)
	if err != nil {
		log.Error("d.dao.ChannelGroup(%d) error(%v)", tid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		t := &model.ChannelGroup{}
		if err = rows.Scan(&t.Id, &t.Ptid, &t.Tid, &t.Alias, &t.Rank, &t.CTime, &t.MTime); err != nil {
			log.Error("d.dao.ChannelGroup(%d) rows.scan() error(%v)", tid, err)
			return
		}
		res = append(res, t)
	}
	err = rows.Err()
	return
}
