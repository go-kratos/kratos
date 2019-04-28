package tag

import (
	"context"

	"go-common/app/interface/main/app-feed/conf"
	tag "go-common/app/interface/main/tag/model"
	tagrpc "go-common/app/interface/main/tag/rpc/client"
	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

type Dao struct {
	// conf
	// http client
	hot    string
	add    string
	cancel string
	tags   string
	detail string
	client *httpx.Client
	tagRPC *tagrpc.Service
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPTag),
		hot:    c.Host.APICo + _hot,
		add:    c.Host.APICo + _add,
		cancel: c.Host.APICo + _cancel,
		tags:   c.Host.APICo + _tags,
		detail: c.Host.APICo + _detail,
		tagRPC: tagrpc.New2(c.TagRPC),
	}
	return
}

func (d *Dao) InfoByIDs(c context.Context, mid int64, tids []int64) (tm map[int64]*tag.Tag, err error) {
	var ts []*tag.Tag
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &tag.ArgIDs{IDs: tids, Mid: mid, RealIP: ip}
	if ts, err = d.tagRPC.InfoByIDs(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	tm = make(map[int64]*tag.Tag, len(ts))
	for _, t := range ts {
		tm[t.ID] = t
	}
	return
}

func (d *Dao) SubTags(c context.Context, mid, vmid int64, pn, ps int) (sub *tag.Sub, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	arg := &tag.ArgSub{Mid: mid, Vmid: vmid, Pn: pn, Ps: ps, RealIP: ip}
	if sub, err = d.tagRPC.SubTags(c, arg); err != nil {
		if err == ecode.NothingFound {
			err = nil
			return
		}
		err = errors.Wrapf(err, "%v", arg)
	}
	return
}
