package dao

import (
	"context"
	"fmt"
	"go-common/library/log"
)

const (
	_setOriginStreamingStatus = "UPDATE `sv_ls_stream` SET `up_rank` = ? where `room_id` = ? and `src` = ? and `up_rank` = ? LIMIT 1;"
)

// SetOriginStreamingStatus 用于设置 老版本数据结构的 推流状态
func (d *Dao) SetOriginStreamingStatus(c context.Context, rid int64, src, from, to int) error {
	res, err := d.stmtLegacyStreamNotify.Exec(c, to, rid, src, from)
	if err != nil {
		return err
	}
	er, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if er == 0 {
		log.Infow(c, "no_record_updated", fmt.Sprintf("%d_%d_%d_%d", rid, src, from, to))
		return nil
	}
	return err
}
