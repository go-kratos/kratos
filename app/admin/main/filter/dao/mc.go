package dao

import (
	"context"
	"fmt"

	"go-common/library/log"
)

func mcKey(key, area string) string {
	return fmt.Sprintf("%s_%s", key, area)
}

// DelKeyAreaCache .
func (d *Dao) DelKeyAreaCache(c context.Context, key, area string) (err error) {
	conn := d.mc.Get(c)
	defer conn.Close()
	mcKey := mcKey(key, area)
	if err = conn.Delete(mcKey); err != nil {
		log.Error("conn.Delete(%s) error(%v)", mcKey, err)
		return
	}
	return
}
