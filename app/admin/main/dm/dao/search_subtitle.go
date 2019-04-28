package dao

import (
	"context"

	"go-common/app/admin/main/dm/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

var (
	_subtitleFields = []string{"oid", "id"}
)

// SearchSubtitle .
func (d *Dao) SearchSubtitle(c context.Context, arg *model.SubtitleSearchArg) (res *model.SearchSubtitleResult, err error) {
	var (
		req    *elastic.Request
		fields []string
	)
	fields = _subtitleFields
	req = d.esCli.NewRequest("dm_subtitle").Index("subtitle").Fields(fields...).Pn(int(arg.Pn)).Ps(int(arg.Ps))
	if arg.Aid > 0 {
		req.WhereEq("aid", arg.Aid)
	}
	if arg.Mid > 0 {
		req.WhereEq("mid", arg.Mid)
	}
	if arg.Oid > 0 {
		req.WhereEq("oid", arg.Oid)
	}
	if arg.Mid > 0 {
		req.WhereEq("mid", arg.Mid)
	}
	if arg.Status > 0 {
		req.WhereEq("status", arg.Status)
	}
	if arg.UpperMid > 0 {
		req.WhereEq("up_mid", arg.UpperMid)
	}
	if arg.Lan > 0 {
		req.WhereEq("lan", arg.Lan)
	}
	req.Order("mtime", "desc")
	if err = req.Scan(c, &res); err != nil {
		log.Error("elastic search(%s) error(%v)", req.Params(), err)
		return
	}
	return
}
