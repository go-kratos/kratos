package tag

import (
	"context"

	"go-common/app/interface/main/app-tag/conf"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

type Dao struct {
	// conf
	c                   *conf.Config
	client              *httpx.Client
	tagRPC              *tagrpc.Service
	detailURL           string
	tagHotsIDURL        string
	similarTagChangeURL string
	tagURL              string
	tagRankingURL       string
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:                   c,
		client:              httpx.NewClient(conf.Conf.HTTPClient),
		tagRPC:              tagrpc.New2(c.TagRPC),
		detailURL:           c.Host.ApiCo + _detail,
		tagHotsIDURL:        c.Host.ApiCo + _tagHotsIDURL,
		similarTagChangeURL: c.Host.ApiCo + _similarTagChangeURL,
		tagURL:              c.Host.ApiCo + _tagURL,
		tagRankingURL:       c.Host.ApiCo + _tagRankingURL,
	}
	return
}

// InfoByID by tag id
func (d *Dao) InfoByID(c context.Context, mid, tid int64) (t *tag.Tag, err error) {
	arg := &tag.ArgID{ID: tid, Mid: mid}
	if t, err = d.tagRPC.InfoByID(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			log.Error("tagRPC.InfoByID(%v) error(%v)", arg, err)
		}
	}
	return
}

// InfoByName by tag name
func (d *Dao) InfoByName(c context.Context, tname string) (t *tag.Tag, err error) {
	arg := &tag.ArgName{Name: tname}
	if t, err = d.tagRPC.InfoByName(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			log.Error("tagRPC.InfoByName(%v) error(%v)", arg, err)
		}
	}
	return
}
