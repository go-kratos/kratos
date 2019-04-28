package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

//AddComment add comment for performance test task
func (d *Dao) AddComment(comment *model.Comment) error {
	//comment.Status = 1
	return d.DB.Table(model.Comment{}.TableName()).Create(comment).Error
}

//QueryComment query comment
func (d *Dao) QueryComment(comment *model.Comment) (res *model.QueryCommentResponse, err error) {
	res = &model.QueryCommentResponse{}
	//comment.Status = 1
	err = d.DB.Table(model.Comment{}.TableName()).Where(comment).Count(&res.Total).Order("id desc").Find(&res.Comments).Error
	return
}

//UpdateComment update comment
func (d *Dao) UpdateComment(comment *model.Comment) error {
	//return d.DB.Table(model.Comment{}.TableName()).Update(comment).Where("ID=?", comment.ID).Error
	return d.DB.Model(&model.Comment{}).Update(comment).Where("ID=?", comment.ID).Error
}

//DeleteComment delete comment
func (d *Dao) DeleteComment(id int64) error {
	return d.DB.Model(&model.Comment{}).Where("ID=?", id).Update("status", 2).Error
}
