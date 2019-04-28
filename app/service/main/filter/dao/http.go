package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"go-common/app/service/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_getScore   = "/nlpinfer/realtime"
	_replyDel   = "/x/admin/reply/internal/del"
	_replyLabel = "/x/admin/reply/internal/spam"
)

// AiScore get ai score.
func (d *Dao) AiScore(c context.Context, content string, stype string) (res *model.AiScore, err error) {
	params := url.Values{}
	var commentArg struct {
		Comments []string `json:"comments"`
	}
	var comments []string
	comments = append(comments, content)
	commentArg.Comments = comments
	cc, err := json.Marshal(commentArg)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	params.Set("comments", string(cc))
	var area = ""
	switch stype {
	case "reply":
		area = "comment"
	case "live_danmu":
		area = "dmlive"
	case "danmu":
		area = "main"
	default:
		err = ecode.RequestErr
		return
	}
	if area != "" {
		params.Set("service", area)
	}
	res = &model.AiScore{}
	if err = d.httpClient.Post(c, d.aiScoreURL, "", params, res); err != nil {
		err = errors.Wrapf(err, "AiScore(%s) error(%v)", d.aiScoreURL+"?"+params.Encode(), err)
		return
	}
	log.Info("AiScore(%s) res(%+v)", d.aiScoreURL+"?"+params.Encode(), res)
	return
}

// ReplyDel reply delete.
func (d *Dao) ReplyDel(c context.Context, adid, oid, rpid int64, avType int8) (err error) {
	params := url.Values{}
	params.Set("adid", fmt.Sprintf("%d", adid))
	params.Set("oid", fmt.Sprintf("%d", oid))
	params.Set("rpid", fmt.Sprintf("%d", rpid))
	params.Set("type", fmt.Sprintf("%d", avType))
	params.Set("moral", "0")
	params.Set("notify", fmt.Sprintf("%t", true))

	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, d.mngReplyDelURL, "", params, &res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		err = errors.Errorf("filter service ReplyDel failed code(%d)", res.Code)
		return
	}
	return
}

// ReplyLabel reply label.
func (d *Dao) ReplyLabel(c context.Context, adid, oid, rpid int64, avType int8) (err error) {
	params := url.Values{}
	params.Set("adid", fmt.Sprintf("%d", adid))
	params.Set("oid", fmt.Sprintf("%d", oid))
	params.Set("rpid", fmt.Sprintf("%d", rpid))
	params.Set("type", fmt.Sprintf("%d", avType))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.httpClient.Post(c, d.mngReplyLabelURL, "", params, &res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		err = errors.Errorf("filter service ReplyLabel failed code (%d)", res.Code)
		return
	}
	return
}
