package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/answer/model"
	"go-common/library/ecode"
)

func (s *Service) convertModel(rpcRes *model.AnsQuesList) (res *model.AnsQueDetailList) {
	res = &model.AnsQueDetailList{CurrentTime: rpcRes.CurrentTime.Unix(), EndTime: rpcRes.EndTime.Unix()}
	for _, q := range rpcRes.QuesList {
		que := &model.AnsQueDetail{
			ID: q.ID, AnsImg: q.Img, QsHeight: q.Height, QsPositionY: q.PositionY,
			Ans1Hash: q.Ans[0].AnsHash, Ans0Height: q.Ans[0].Height, Ans0PositionY: q.Ans[0].PositionY,
			Ans2Hash: q.Ans[1].AnsHash, Ans1Height: q.Ans[1].Height, Ans1PositionY: q.Ans[1].PositionY,
			Ans3Hash: q.Ans[2].AnsHash, Ans2Height: q.Ans[2].Height, Ans2PositionY: q.Ans[2].PositionY,
			Ans4Hash: q.Ans[3].AnsHash, Ans3Height: q.Ans[3].Height, Ans3PositionY: q.Ans[3].PositionY,
		}
		res.QuesList = append(res.QuesList, que)
	}
	return
}

func (s *Service) convertExtraModel(rpcRes *model.AnsQuesList) (res *model.AnsQueDetailList) {
	res = &model.AnsQueDetailList{CurrentTime: rpcRes.CurrentTime.Unix(), EndTime: rpcRes.EndTime.Unix()}
	for _, q := range rpcRes.QuesList {
		que := &model.AnsQueDetail{
			ID: q.ID, AnsImg: q.Img, QsHeight: q.Height, QsPositionY: q.PositionY,
			Ans1Hash: q.Ans[0].AnsHash, Ans0Height: q.Ans[0].Height, Ans0PositionY: q.Ans[0].PositionY,
			Ans2Hash: q.Ans[1].AnsHash, Ans1Height: q.Ans[1].Height, Ans1PositionY: q.Ans[1].PositionY,
		}
		res.QuesList = append(res.QuesList, que)
	}
	return
}

func (s *Service) historyByHid(ctx context.Context, hid int64) (his *model.AnswerHistory, err error) {
	his, err = s.answerDao.HidCache(ctx, hid)
	if err != nil {
		return
	}
	if his != nil {
		return
	}
	if len(fmt.Sprintf("%d", hid)) < 10 {
		i, _ := s.answerDao.SharingIndexByHid(ctx, hid)
		his, err = s.answerDao.OldHistory(ctx, hid, i)
	} else {
		his, err = s.answerDao.HistoryByHid(ctx, hid)
	}
	if err != nil || his == nil {
		err = ecode.NothingFound
		return
	}
	s.answerDao.SetHidCache(ctx, his)
	return
}
