package dao

import (
	"context"

	recall "go-common/app/service/bbq/recsys-recall/api/grpc/v1"
)

// SetInvertedIndex 倒排写入redis
func (d *Dao) SetInvertedIndex(c context.Context, key string, svids []int64) error {
	_, err := d.recallClient.NewIncomeVideo(c, &recall.NewIncomeVideoRequest{
		Key:   key,
		SVIDs: svids,
	})
	return err
}
