package web

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// PgcFull pgc full .
func (d *Dao) PgcFull(ctx context.Context, tp int, pn, ps int64, source string) (res interface{}, err error) {
	var (
		param = url.Values{}
		ip    = metadata.String(ctx, metadata.RemoteIP)
		rs    struct {
			Code int         `json:"code"`
			Data interface{} `json:"result"`
		}
	)
	param.Set("bsource", source)
	param.Set("season_type", strconv.Itoa(tp))
	param.Set("page_no", strconv.FormatInt(pn, 10))
	param.Set("page_size", strconv.FormatInt(ps, 10))
	if err = d.httpR.Get(ctx, d.pgcFullURL, ip, param, &rs); err != nil {
		log.Error("d.httpR.Get err(%v)", err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}

// PgcIncre pgc increment .
func (d *Dao) PgcIncre(ctx context.Context, tp int, pn, ps, start, end int64, source string) (res interface{}, err error) {
	var (
		param = url.Values{}
		ip    = metadata.String(ctx, metadata.RemoteIP)
	)
	var rs struct {
		Code int         `json:"code"`
		Data interface{} `json:"result"`
	}
	param.Set("bsource", source)
	param.Set("season_type", strconv.Itoa(tp))
	param.Set("page_no", strconv.FormatInt(pn, 10))
	param.Set("page_size", strconv.FormatInt(ps, 10))
	param.Set("start_ts", strconv.FormatInt(start, 10))
	param.Set("end_ts", strconv.FormatInt(end, 10))
	if err = d.httpR.Get(ctx, d.pgcIncreURL, ip, param, &rs); err != nil {
		log.Error("d.httpR.Get url(%s) err(%s)", d.pgcIncreURL+"?"+param.Encode(), err)
		return
	}
	if rs.Code != ecode.OK.Code() {
		err = ecode.Int(rs.Code)
		return
	}
	res = rs.Data
	return
}
