package resource

import (
	"context"

	"go-common/library/ecode"
	"go-common/library/log"
)

// ColumnInfo column category
type columnInfo struct {
	ID   int64  `json:"id"`
	PID  int64  `json:"parent_id"`
	Name string `json:"name"`
}

// ColumnCategory get column types.
func columnCategory(c context.Context) (data []*columnInfo, err error) {
	var res struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    []*columnInfo `json:"data"`
	}
	url := articleCategoryURL
	if err = client.Get(c, url, "", nil, &res); err != nil {
		log.Error("resource.columnCategory GET error(%v) | uri(%s)", err, url)
		return
	}
	if res.Code != 0 {
		log.Error("resource.columnCategory code != 0. res.Code(%d) | uri(%s) res(%v)", res.Code, url, res)
		err = ecode.GrowupGetTypeError
		return
	}
	data = res.Data
	return
}

// ColumnCategoryNameToID .
func ColumnCategoryNameToID(c context.Context) (categories map[string]int64, err error) {
	data, err := columnCategory(c)
	if err != nil {
		return
	}
	categories = make(map[string]int64, len(data))
	for _, v := range data {
		if v.PID == 0 {
			categories[v.Name] = v.ID
		}
	}
	return
}
