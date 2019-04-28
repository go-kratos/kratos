package service

import (
	"time"

	"go-common/app/admin/ep/melloi/model"
)

//AddComment add comment for test job
func (s *Service) AddComment(comment *model.Comment) error {
	comment.Status = 1
	comment.SubmitDate = time.Now()
	return s.dao.AddComment(comment)
}

//QueryComment query comment
func (s *Service) QueryComment(comment *model.Comment) (*model.QueryCommentResponse, error) {
	return s.dao.QueryComment(comment)
}

//UpdateComment update comment
func (s *Service) UpdateComment(comment *model.Comment) error {
	return s.dao.UpdateComment(comment)
}

//DeleteComment delete comment
func (s *Service) DeleteComment(id int64) error {
	return s.dao.DeleteComment(id)
}
