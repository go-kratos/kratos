package bigdata

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"go-common/library/log"

	"github.com/pkg/errors"
)

type topicResponse struct {
	Code    int      `json:"error_code"`
	Message string   `json:"error_message"`
	Topics  []string `json:"topics"`
}

type topicReq struct {
	Mid     int64  `json:"mid"`
	Oid     int64  `json:"oid"`
	Type    int8   `json:"type"`
	Message string `json:"message"`
}

// Topics return topics
func (dao *Dao) Topics(c context.Context, mid int64, oid int64, typ int8, msg string) ([]string, error) {
	res := &topicResponse{}
	content, err := json.Marshal(&topicReq{
		Mid:     mid,
		Oid:     oid,
		Type:    typ,
		Message: msg,
	})
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", dao.topicURL, bytes.NewReader(content))
	if err != nil {
		log.Error("bigdata.Topics(%d,%d,%d,%s) url(%s) req(%s)send POST error(%v)", mid, oid, typ, msg, dao.topicURL, string(content), err)
		err = errors.WithStack(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	err = dao.httpClient.Do(c, req, res)
	if err != nil {
		log.Error("bigdata.Topics(%d,%d,%d,%s) url(%s) req(%s) error(%v)", mid, oid, typ, msg, dao.topicURL, string(content), err)
		return nil, err
	}
	if res.Code != 0 {
		log.Error("bigdata.Topics(%d,%d,%d,%s)  url(%s) req(%s)  return not success,error_msg(%d,%v)", mid, oid, typ, msg, dao.topicURL, string(content), res.Code, res.Message)
		return nil, err
	}
	return res.Topics, nil
}
