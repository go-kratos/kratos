package service

import (
	"context"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/model"
	spymdl "go-common/app/service/main/spy/model"
	"go-common/library/log"
)

// FilterAiScore get ai score .
func (s *Service) FilterAiScore(c context.Context, content string, mid, adid, oid, rpid int64, avtype int8) (err error) {
	if _, ok := s.whiteMids.White(mid); ok {
		log.Info("FilterAiScore white hit s.conf.Ai.Whites(%d %s)", mid, content)
		return
	}
	var (
		aiScore *model.AiScore
	)
	if aiScore, err = s.dao.AiScore(c, content, "reply"); err != nil {
		return
	}
	if len(aiScore.Scores) > 0 && (aiScore.Scores[0] < aiScore.Threshold) {
		log.Info("FilterAiScore aiScore.Scores(%+v) aiScore.Threshold(%+v)", aiScore.Scores, aiScore.Threshold)
		return
	}
	var (
		score     *spymdl.UserScore
		threshold float64
	)
	if mid > 0 {
		arg := &spymdl.ArgUserScore{
			Mid: mid,
		}
		score, err = s.spyRPC.UserScore(c, arg)
		if err != nil {
			log.Error("s.spyRPC(%d) err(%v)", mid, err)
			return
		}
		log.Info("FilterAiScore s.spyRPC.UserScore(%d) res(%+v)", mid, score)
	}
	if mid == 0 {
		log.Errorv(c, log.KV("log", "ai filter mid=0"))
	}
	log.Info("FilterAiScore s.conf.Ai.Threshold(%+v) s.conf.Ai.TrueScore(+v)", conf.Conf.Property.AI.Threshold, conf.Conf.Property.AI.TrueScore)
	threshold = aiScore.Threshold
	if conf.Conf.Property.AI.Threshold > 0 {
		threshold = conf.Conf.Property.AI.Threshold
	}
	if (len(aiScore.Scores) > 0 && (aiScore.Scores[0] > threshold)) && (float64(score.Score) < conf.Conf.Property.AI.TrueScore) {
		s.addAICh(func() {
			s.dao.ReplyDel(c, adid, oid, rpid, avtype)
		})
		return
	}
	s.addAICh(func() {
		s.dao.ReplyLabel(c, adid, oid, rpid, avtype)
	})
	return
}
