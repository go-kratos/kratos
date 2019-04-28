package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"

	"github.com/pkg/errors"
)

// PgcCond picks pgc condition
func (d *Dao) PgcCond(c context.Context, snType int32) (result *model.PgcCond, err error) {
	var (
		host   = d.c.Cfg.RefLabel.PgcAPI
		params = url.Values{}
		resp   = model.PgcCondResp{}
	)
	params.Set("season_type", fmt.Sprintf("%d", snType))
	if err = d.client.Get(c, host, "", params, &resp); err != nil {
		return
	}
	if resp.Code != ecode.OK.Code() {
		err = errors.Wrapf(ecode.Int(resp.Code), resp.Message)
		return
	}
	result = resp.Result
	return
}
