package dao

import (
	"context"
	"errors"
	"fmt"
	"go-common/app/service/live/gift/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

var _getDiscountGift = "SELECT id,discount_id,gift_id,user_type,discount_price,corner_mark,corner_position FROM discount_gift WHERE discount_id in (%s)"

// GetByDiscountIds GetByDiscountIds
func (d *Dao) GetByDiscountIds(ctx context.Context, ids []int64) (res []*model.DiscountGift, err error) {
	log.Info("GetByDiscountIds,ids:%v", ids)
	if len(ids) == 0 {
		log.Error("query GetByDiscountIds params null")
		err = errors.New("params error")
		return
	}
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, fmt.Sprintf(_getDiscountGift, xstr.JoinInts(ids))); err != nil {
		log.Error("query GetByDiscountIds error,err %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		d := &model.DiscountGift{}
		if err = rows.Scan(&d.Id, &d.DiscountId, &d.GiftId, &d.UserType, &d.DiscountPrice, &d.CornerMark, &d.CornerPosition); err != nil {
			log.Error("GetByDiscountIds scan error,err %v", err)
			return
		}
		res = append(res, d)
	}
	return
}
