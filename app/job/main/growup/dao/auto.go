package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"go-common/library/log"
	"go-common/library/xstr"
)

// DoAvBreach av breach by api.
func (d *Dao) DoAvBreach(c context.Context, mid int64, aid int64, ctype int, reason string) (err error) {
	params := url.Values{}
	params.Set("type", strconv.Itoa(ctype))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("aids", strconv.FormatInt(aid, 10))
	params.Set("reason", reason)

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	url := d.breachURL
	if err = d.client.Post(c, url, "", params, &res); err != nil {
		log.Error("d.client.Post url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DoAvBreach code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), res.Message)
		err = fmt.Errorf("DoAvBreach error(%s)", res.Message)
	}
	return
}

// DoUpForbid up forbid by api
func (d *Dao) DoUpForbid(c context.Context, mid int64, days int, ctype int, reason string) (err error) {
	params := url.Values{}
	params.Set("type", strconv.Itoa(ctype))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("days", strconv.Itoa(days))
	params.Set("reason", reason)

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	url := d.forbidURL
	if err = d.client.Post(c, url, "", params, &res); err != nil {
		log.Error("d.client.Post url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DoUpForbid code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), res.Message)
		err = fmt.Errorf("DoUpForbid error(%s)", res.Message)
	}
	return
}

// DoUpDismiss up dismiss by api
func (d *Dao) DoUpDismiss(c context.Context, mid int64, ctype int, reason string) (err error) {
	params := url.Values{}
	params.Set("type", strconv.Itoa(ctype))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("reason", reason)

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	url := d.dismissURL
	if err = d.client.Post(c, url, "", params, &res); err != nil {
		log.Error("d.client.Post url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DoUpDismiss code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), res.Message)
		err = fmt.Errorf("DoUpDismiss error(%s)", res.Message)
	}
	return
}

// DoUpPass up pass by api
func (d *Dao) DoUpPass(c context.Context, mids []int64, ctype int) (err error) {
	params := url.Values{}
	params.Set("type", strconv.Itoa(ctype))
	params.Set("mids", xstr.JoinInts(mids))

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	url := d.passURL
	if err = d.client.Post(c, url, "", params, &res); err != nil {
		log.Error("d.client.Post url(%s) error(%v)", url+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("DoUpPass code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, url+"?"+params.Encode(), res.Message)
		err = fmt.Errorf("DoUpPass error(%s)", res.Message)
	}
	return
}
