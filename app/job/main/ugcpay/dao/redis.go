package dao

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func assetRelationKey(mid int64) string {
	return fmt.Sprintf("up_ar_%d", mid)
}

func assetRelationField(oid int64, otype string) string {
	return fmt.Sprintf("%s_%d", otype, oid)
}

// DelCacheAssetRelationState delete asset relation state.
func (d *Dao) DelCacheAssetRelationState(c context.Context, oid int64, otype string, mid int64) (err error) {
	var (
		key   = assetRelationKey(mid)
		field = assetRelationField(oid, otype)
		conn  = d.redis.Get(c)
	)
	defer conn.Close()
	if _, err = conn.Do("HDEL", key, field); err != nil {
		err = errors.Wrapf(err, "conn.Do(HDEL, %s, %s)", key, field)
		return
	}
	return
}
