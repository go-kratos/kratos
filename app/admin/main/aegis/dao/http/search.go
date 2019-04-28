package http

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_upsertES = "/x/admin/search/upsert"
)

// ResourceES search archives by es.
func (d *Dao) ResourceES(c context.Context, arg *model.SearchParams) (sres *model.SearchRes, err error) {
	r := d.es.NewRequest("aegis_resource").Index("aegis_resource").Fields(
		"id",
		"business_id",
		"flow_id",
		"oid",
		"mid",
		"content",
		"extra1",
		"extra2",
		"extra3",
		"extra4",
		"extra5",
		"extra6",
		"extra1s",
		"extra2s",
		"extra3s",
		"extra4s",
		"extratime1",
		"octime",
		"ptime",
		"metadata",
		"note",
		"reject_reason",
		"reason_id",
		"state",
		"ctime",
	).OrderScoreFirst(false)

	escm := model.EsCommon{
		Ps:    arg.Ps,
		Pn:    arg.Pn,
		Order: "ctime",
		Sort:  strings.ToLower(arg.CtimeOrder),
	}
	if escm.Sort != "asc" && escm.Sort != "desc" {
		escm.Sort = "desc"
	}

	setESParams(r, arg, escm)
	if arg.KeyWord != "" { //描述
		arg.KeyWord = strings.Replace(arg.KeyWord, "，", ",", -1)
		r.WhereLike([]string{"content"}, strings.Split(arg.KeyWord, ","), true, elastic.LikeLevelHigh)
	}

	log.Info("ResourceES params(%s)", r.Params())

	sres = &model.SearchRes{}
	if err = r.Scan(c, sres); err != nil {
		log.Error("ResourceES r.Scan params(%s)|error(%v)", r.Params(), err)
		return
	}
	arg.Pn = sres.Page.Num
	arg.Ps = sres.Page.Size
	arg.Total = sres.Page.Total
	return
}

//UpsertES 更新搜索
func (d *Dao) UpsertES(c context.Context, rsc []*model.UpsertItem) (err error) {
	if len(rsc) == 0 {
		return
	}

	items := []*model.UpsertItem{}
	for _, item := range rsc {
		if item == nil || item.ID <= 0 {
			continue
		}
		items = append(items, item)
	}

	data := map[string][]*model.UpsertItem{
		"aegis_resource": items,
	}
	datab, err := json.Marshal(data)
	if err != nil {
		log.Error("UpsertES json.Marshal error(%v) resource(%+v)", err, rsc)
		return err
	}

	res := new(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	})
	params := url.Values{}
	params.Set("business", "aegis_resource")
	params.Set("insert", "false")
	params.Set("data", string(datab))
	if err = d.clientW.Post(c, d.c.Host.Manager+_upsertES, "", params, res); err != nil {
		log.Error("UpsertES d.clientW.Post error(%v) params(%+v)", err, params)
		return
	}
	if res.Code != ecode.OK.Code() {
		log.Error("UpsertES d.clientW.Post failed, response(%+v) params(%+v)", res, params)
		return
	}

	log.Info("response(%+v)  url=%s%s?%s", res, d.c.Host.Manager, _upsertES, params.Encode())
	return
}
