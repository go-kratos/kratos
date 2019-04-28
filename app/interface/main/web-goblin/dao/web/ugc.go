package web

import (
	"context"
	"time"

	webmdl "go-common/app/interface/main/web-goblin/model/web"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

const _ugcIncre = "web_goblin"

// UgcIncre ugc increment .
func (d *Dao) UgcIncre(ctx context.Context, pn, ps int, start, end int64) (res []*webmdl.SearchAids, err error) {
	var (
		startStr, endStr string
		rs               struct {
			Result []*webmdl.SearchAids `json:"result"`
		}
	)
	startStr = time.Unix(start, 0).Format("2006-01-02 15:04:05")
	endStr = time.Unix(end, 0).Format("2006-01-02 15:04:05")
	r := d.ela.NewRequest(_ugcIncre).WhereRange("mtime", startStr, endStr, elastic.RangeScopeLoRo).Fields("aid").Fields("action").Index(_ugcIncre).Pn(pn).Ps(ps)
	if err = r.Scan(ctx, &rs); err != nil {
		log.Error("r.Scan error(%v)", err)
		return
	}
	res = rs.Result
	return
}
