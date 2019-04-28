package http

import (
	"encoding/json"
	"strings"

	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func upsert(c *bm.Context) {
	up := &model.UpsertParams{}
	if err := c.Bind(up); err != nil {
		return
	}
	dataBody := map[string][]model.MapData{}
	decoder := json.NewDecoder(strings.NewReader(up.DataStr))
	decoder.UseNumber()
	if err := decoder.Decode(&dataBody); err != nil {
		log.Error("s.http.upsert(%v) json error(%v)", err, dataBody)
	}
	if len(dataBody) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	for _, n := range dataBody {
		for _, m := range n {
			if err := m.NumberToInt64(); err != nil {
				log.Error("s.http.upsert(%v) to int64 error(%v)", err, m)
			}
		}
	}
	c.JSON(nil, svr.Upsert(c, up, dataBody))
}
