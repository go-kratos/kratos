package block

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// SendSysMsg send sys msg.
func (d *Dao) SendSysMsg(c context.Context, code string, mids []int64, title string, content string, remoteIP string) (err error) {
	params := url.Values{}
	params.Set("mc", code)
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", content)
	params.Set("mid_list", midsToParam(mids))
	var res struct {
		Code int `json:"code"`
		Data *struct {
			Status int8   `json:"status"`
			Remark string `json:"remark"`
		} `json:"data"`
	}
	if err = d.client.Post(c, d.c.BlockProperty.MSGURL, remoteIP, params, &res); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res.Code != 0 {
		err = errors.WithStack(ecode.Int(res.Code))
		return
	}
	return
}

func midsToParam(mids []int64) (str string) {
	strs := make([]string, 0, len(mids))
	for _, mid := range mids {
		strs = append(strs, fmt.Sprintf("%d", mid))
	}
	return strings.Join(strs, ",")
}
