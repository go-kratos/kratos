package dao

import (
	"context"
	"encoding/json"
	"net/url"

	"go-common/app/admin/main/filter/model"
	"go-common/library/log"
)

const (
	_getScore = "/nlpinfer/realtime"
)

// AiScore get ai score.
func (dao *Dao) AiScore(c context.Context, content string) (res *model.AiScore, err error) {
	params := url.Values{}
	var commentArg struct {
		Comments []string `json:"comments"`
	}
	var (
		comments []string
		cc       []byte
	)
	comments = append(comments, content)
	commentArg.Comments = comments
	cc, err = json.Marshal(commentArg)
	if err != nil {
		log.Error("AiScore json.Marshal(%+v) error(%v)", comments, err)
		return
	}
	params.Set("comments", string(cc))
	params.Set("service", "comment")
	res = &model.AiScore{}
	if err = dao.client.Post(c, dao.aiScoreURL, "", params, res); err != nil {
		log.Error("AiScore(%s) error(%v)", dao.aiScoreURL+"?"+params.Encode(), err)
		return
	}
	log.Info("AiScore(%s) res(%+v)", dao.aiScoreURL+"?"+params.Encode(), res)
	return
}
