package dao

import (
	"fmt"

	"go-common/app/admin/main/apm/model/need"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_tableNeed    = "needs"
	_tableLike    = "user_likes"
	_addCountsSQL = "UPDATE needs SET like_counts=like_counts+?,dislike_counts=dislike_counts+? WHERE id=?"
)

//NeedInfoAdd add need info
func (d *Dao) NeedInfoAdd(r *need.NAddReq, username string) (err error) {
	ni := &need.NInfo{
		Title:    r.Title,
		Content:  r.Content,
		Reporter: username,
		Status:   1,
	}
	if err = d.DB.Create(&ni).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//NeedInfoList all need info
func (d *Dao) NeedInfoList(arg *need.NListReq) (res []*need.NInfo, err error) {

	where := d.needInfoCondition(arg)
	if err = d.DB.Table(_tableNeed).Where(where).Order("id DESC").Offset((arg.Pn - 1) * arg.Ps).Limit(arg.Ps).Find(&res).Error; err != nil {
		err = errors.WithStack(err)
		log.Error("NeedInfoList:%s", err)
		return
	}
	return
}

//NeedInfoCount  need info count
func (d *Dao) NeedInfoCount(arg *need.NListReq) (count int64, err error) {

	where := d.needInfoCondition(arg)
	if err = d.DB.Table(_tableNeed).Where(where).Count(&count).Error; err != nil {
		log.Error("NeedInfoCount:%s", err)
		return
	}
	return
}

//needInfoCondition is
func (d *Dao) needInfoCondition(arg *need.NListReq) (where string) {
	where = "status !=5"
	if arg.Status > 0 {
		where += fmt.Sprintf(" and `status`='%d'", arg.Status)
	}
	if arg.Reporter != "" {
		where += fmt.Sprintf(" and `reporter`='%s'", arg.Reporter)
	}
	return where
}

// GetNeedInfo is
func (d *Dao) GetNeedInfo(id int64) (r *need.NInfo, err error) {
	r = &need.NInfo{}
	if err = d.DB.Table(_tableNeed).Where("id=?", id).Find(r).Error; err != nil {
		log.Error("GetNeedInfo:%s", err)
		return
	}
	return
}

//NeedInfoEdit is
func (d *Dao) NeedInfoEdit(arg *need.NEditReq) (err error) {

	if err = d.DB.Table(_tableNeed).Where("id=?", arg.ID).Updates(map[string]interface{}{"content": arg.Content, "title": arg.Title}).Error; err != nil {
		log.Error("NeedInfoEdit:%s", err)
		return
	}
	return
}

//NeedVerify is
func (d *Dao) NeedVerify(r *need.NVerifyReq) (err error) {
	tx := d.DB.Begin()
	if err = d.DB.Table(_tableNeed).Where("id=?", r.ID).Update("status", r.Status).Error; err != nil {
		log.Error("NeedVerify:%s", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//LikeCountsStats thumbsup counts
func (d *Dao) LikeCountsStats(r *need.Likereq, like, dislike int) (err error) {
	tx := d.DB.Begin()
	if err = d.DB.Exec(_addCountsSQL, like, dislike, r.ReqID).Error; err != nil {
		err = errors.WithStack(err)
		log.Error("LikeCountsStats:%s", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//GetVoteInfo is
func (d *Dao) GetVoteInfo(r *need.Likereq, username string) (u *need.UserLikes, err error) {
	u = &need.UserLikes{}
	if err = d.DB.Table(_tableLike).Where("req_id=? and `user`=?", r.ReqID, username).Find(u).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//UpdateVoteInfo is
func (d *Dao) UpdateVoteInfo(r *need.Likereq, username string) (err error) {
	if err = d.DB.Table(_tableLike).Where("req_id=? and `user`=?", r.ReqID, username).Update("like_type", r.LikeType).Error; err != nil {
		err = errors.WithStack(err)
		log.Error("UpdateVoteInfo:%s", err)
		return
	}
	return
}

//AddVoteInfo is
func (d *Dao) AddVoteInfo(r *need.Likereq, username string) (err error) {
	ul := &need.UserLikes{
		ReqID:    r.ReqID,
		User:     username,
		LikeType: r.LikeType,
	}
	if err = d.DB.Create(&ul).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

//voteInfoCondition is
func (d *Dao) voteInfoCondition(arg *need.Likereq) (where string) {

	where = " like_type != 0 "
	if arg.ReqID > 0 {
		where += fmt.Sprintf(" and `req_id`='%d'", arg.ReqID)
	}
	if arg.LikeType > 0 {
		where += fmt.Sprintf(" and `like_type`='%d'", arg.LikeType)
	}
	return where
}

//VoteInfoList is  vote info
func (d *Dao) VoteInfoList(arg *need.Likereq) (res []*need.UserLikes, err error) {

	where := d.voteInfoCondition(arg)
	if err = d.DB.Table(_tableLike).Where(where).Find(&res).Error; err != nil {
		err = errors.WithStack(err)
		log.Error("VoteInfoList:%s", err)
		return
	}
	return
}

//VoteInfoCounts vote info count
func (d *Dao) VoteInfoCounts(arg *need.Likereq) (count int64, err error) {

	where := d.voteInfoCondition(arg)
	if err = d.DB.Table(_tableLike).Where(where).Count(&count).Error; err != nil {
		log.Error("VoteInfoCount:%s", err)
		return
	}
	return
}
