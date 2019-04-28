package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/answer/model"
	"go-common/library/log"
)

// AddLabourQuestion  add labour question.
func (s *Service) AddLabourQuestion(c context.Context, msg *model.MsgCanal) (err error) {
	var que = &model.LabourQs{}
	if err = json.Unmarshal(msg.New, que); err != nil {
		log.Error("json.Unmarshal(%v) err(%v)", msg, err)
		return
	}
	log.Info("labour add (%+v)", que)
	if err = s.createBFSImg(c, que); err != nil {
		log.Error("createBFSImg(%v) err(%v)", que, err)
		return
	}
	que.State = model.HadCreateImg
	s.dao.AddQs(c, que)
	return
}

// ModifyLabourQuestion  nodify labour question.
func (s *Service) ModifyLabourQuestion(c context.Context, msg *model.MsgCanal) (err error) {
	var (
		newq = &model.LabourQs{}
		oldq = &model.LabourQs{}
	)
	if err = json.Unmarshal(msg.New, newq); err != nil {
		log.Error("newlqs json.Unmarshal(%v) err(%v)", msg, err)
		return
	}
	if err = json.Unmarshal(msg.Old, oldq); err != nil {
		log.Error("oldlqs json.Unmarshal(%v) err(%v)", msg, err)
		return
	}
	log.Info("labour modify (%+v)(%+v)", newq, oldq)
	if newq.Status == oldq.Status && newq.Ans == oldq.Ans && newq.Isdel == oldq.Isdel {
		log.Error("ModifyLabourQuestion no change(%v, %v)", newq, oldq)
		return
	}
	s.dao.UpdateQs(c, newq)
	return
}

// UploadQueImg uplaod que img.
func (s *Service) UploadQueImg(c context.Context, que *model.LabourQs) (err error) {
	if err = s.createBFSImg(c, que); err != nil {
		log.Error("createBFSImg(%v) err(%v)", que, err)
		return
	}
	que.State = model.HadCreateImg
	s.dao.UpdateState(c, que)
	return
}
