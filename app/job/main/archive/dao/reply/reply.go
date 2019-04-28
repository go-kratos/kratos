package reply

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"go-common/library/log"
)

const (
	_changeSubjectMid = "/x/internal/v2/reply/subject/mid"
)

// ChangeSubjectMid change av's owner
func (d *Dao) ChangeSubjectMid(oid, mid int64) (err error) {
	params := url.Values{}
	params.Set("adid", "0")
	params.Set("oid", strconv.FormatInt(oid, 10))
	params.Set("type", "1")
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int64 `json:"code"`
	}
	if err = d.client.Post(context.TODO(), d.changeSubMid, "", params, &res); err != nil {
		log.Error("d.client.Post(%s) error(%v)", d.changeSubMid+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = errors.New(strconv.FormatInt(res.Code, 10))
		log.Error("d.client.Post(%s) code(%v)", d.changeSubMid+"?"+params.Encode(), res.Code)
		return
	}
	return
}
