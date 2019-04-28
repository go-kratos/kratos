package archive

import (
	"context"
	"go-common/library/log"
	"net/url"
	"strconv"
)

//GetUserCard get user card
func (d *Dao) GetUserCard(c context.Context, mid int64) (card map[string]interface{}, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))

	res := new(struct {
		Code int                    `json:"code"`
		Card map[string]interface{} `json:"card"`
	})
	card = map[string]interface{}{}
	if err = d.clientR.Get(c, d.userCardURL, "", params, res); err != nil {
		log.Error("GetUserCard d.clientR.Get error(%v) mid(%d)", err, mid)
		return
	}

	if res == nil || res.Code != 0 {
		log.Warn("GetUserCard request failed res(%+v)", res)
		return
	}

	card = res.Card
	return
}
