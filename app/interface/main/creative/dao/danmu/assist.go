package danmu

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// api
	_setDmBannedURI = "/x/internal/dm/assist/banned/upt"
)

// ResetUpBanned pool 0：cancel move，1：cancel ignore
func (d *Dao) ResetUpBanned(c context.Context, mid int64, state int8, hash, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("hash", hash)
	params.Set("stat", strconv.FormatInt(mid, 10)) //0：撤销添加屏蔽，1：撤销删除屏蔽
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.assistDmBannedURL, ip, params, &res); err != nil {
		err = ecode.CreativeDanmuErr
		log.Error("d.SetDmStat.Post(%s,%s,%s) err(%v)", d.assistDmBannedURL, ip, params.Encode(), err)
		return
	}
	if res.Code != 0 {
		err = ecode.Int(res.Code)
		log.Error("d.SetDmStat.Post(%s,%s,%s) err(%v)|code(%d)", d.assistDmBannedURL, ip, params.Encode(), err, res.Code)
		return
	}
	return
}
