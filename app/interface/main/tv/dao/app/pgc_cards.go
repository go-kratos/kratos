package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// PgcCards get season new index from pgc
func (d *Dao) PgcCards(ctx context.Context, ids string) (result map[string]*model.SeasonCard, err error) {
	var (
		response     = &model.PgcResponse{}
		params       = url.Values{}
		pgcSeasonURL = d.conf.Host.APINewindex
	)
	params.Set("season_ids", ids)
	if err = d.client.Get(ctx, pgcSeasonURL, "", params, response); err != nil {
		log.Error("PgcCards New Ep ERROR:%v", err)
		return
	}
	if response.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(response.Code), fmt.Sprintf("PgcCards API Error %v", response.Message))
		log.Error("PgcCards ERROR:%v, URL: %s", err, pgcSeasonURL+"?"+params.Encode())
		return
	}
	if response.Result == nil {
		err = errors.Wrap(ecode.ServerErr, "PgcCards api returns empty")
		log.Error("PgcCards ERROR:%v", err)
		return
	}
	result = response.Result
	return
}
