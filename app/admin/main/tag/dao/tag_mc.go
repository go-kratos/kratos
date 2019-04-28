package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/library/log"
)

const (
	_prefixTag  = "t_%d"
	_prefixName = "n_%s"

	_spaceReplace = "_^_"
)

func keyTag(tid int64) string {
	return fmt.Sprintf(_prefixTag, tid)
}

func keyName(name string) string {
	return fmt.Sprintf(_prefixName, strings.Replace(name, " ", _spaceReplace, -1))
}

// DelTagCache delete tags cache.
func (d *Dao) DelTagCache(c context.Context, tid int64, tname string) error {
	conn := d.mc.Get(c)
	defer conn.Close()
	if err := conn.Delete(keyTag(tid)); err != nil {
		log.Error("conn.Delete(%d) error(%v)", tid, err)
	}
	if err := conn.Delete(keyName(tname)); err != nil {
		log.Error("conn.Delete(%s) error(%v)", keyName(tname), err)
	}
	return nil
}
