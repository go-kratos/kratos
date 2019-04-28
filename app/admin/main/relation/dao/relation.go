package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/relation/model"
	relationPB "go-common/app/service/main/relation/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_shard   = 500
	_maxSize = 20000
)

func midTable(mid int64) string {
	return fmt.Sprintf("user_relation_mid_%03d", mid%_shard)
}

func fidTable(fid int64) string {
	return fmt.Sprintf("user_relation_fid_%03d", fid%_shard)
}

// Followers is
func (d *Dao) Followers(ctx context.Context, fid int64, mid int64) (model.RelationList, error) {
	list := model.RelationList{}
	db := d.ReadORM.Table(fidTable(fid)).Where("status=?", 0).Where("fid=?", fid).Limit(_maxSize).Order("id desc")
	if mid > 0 {
		db = db.Where("mid=?", mid).Limit(1)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	for _, f := range list {
		f.ParseRelation()
	}
	return list, nil
}

// Followings is
func (d *Dao) Followings(ctx context.Context, mid int64, fid int64) (model.RelationList, error) {
	list := model.RelationList{}
	db := d.ReadORM.Table(midTable(mid)).Where("status=?", 0).Where("mid=?", mid).Limit(_maxSize).Order("id desc")
	if fid > 0 {
		db = db.Where("fid=?", fid).Limit(1)
	}
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	for _, f := range list {
		f.ParseRelation()
	}
	return list, nil
}

// Stat is
func (d *Dao) Stat(ctx context.Context, mid int64) (*relationPB.StatReply, error) {
	stat, err := d.relationClient.Stat(ctx, &relationPB.MidReq{
		Mid:    mid,
		RealIp: metadata.String(ctx, metadata.RemoteIP),
	})
	if err != nil {
		log.Error("d.relationRPC.Stat err(%+v)", err)
		return nil, err
	}
	return stat, nil
}

// Stats is
func (d *Dao) Stats(ctx context.Context, mids []int64) (map[int64]*relationPB.StatReply, error) {
	statReply, err := d.relationClient.Stats(ctx, &relationPB.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(ctx, metadata.RemoteIP),
	})
	if err != nil {
		log.Error("d.relationRPC.Stats err(%+v)", err)
		return nil, err
	}
	if len(statReply.StatReplyMap) == 0 {
		return nil, nil
	}
	return statReply.StatReplyMap, nil
}
