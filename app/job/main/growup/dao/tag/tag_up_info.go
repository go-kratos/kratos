package tag

import (
	"context"
	"fmt"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_tagUpInfoByTag = "SELECT tag_id, mid FROM tag_up_info WHERE tag_id in (%s) AND is_deleted = 0 LIMIT ?,?"
)

// GetTagUpInfoByTag get tag_up_info by tag_id
func (d *Dao) GetTagUpInfoByTag(c context.Context, tags []int64, from, limit int, tagMID map[int64][]int64) (count int, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_tagUpInfoByTag, xstr.JoinInts(tags)), from, limit)
	if err != nil {
		log.Error("d.db.Query(%v), error(%v)", _tagUpInfoByTag, err)
		return
	}
	defer rows.Close()

	count = 0
	for rows.Next() {
		var tagID, mid int64
		err = rows.Scan(&tagID, &mid)
		if err != nil {
			log.Error("GetTagUpInfoByTag rows scan error(%v)", err)
			return
		}
		count++
		if _, ok := tagMID[tagID]; !ok {
			tagMID[tagID] = make([]int64, 0)
		}
		tagMID[tagID] = append(tagMID[tagID], mid)
	}
	return
}
