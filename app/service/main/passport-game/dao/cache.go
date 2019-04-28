package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/main/passport-game/model"
	"go-common/library/log"
)

// TokenPBCache get token pb cache.
func (d *Dao) TokenPBCache(c context.Context, key string) (res *model.Perm, err error) {
	if !strings.HasPrefix(key, _keyPrefixTokenPB) {
		err = fmt.Errorf("invalid cache key %s, key pattern is %s{token}", key, _keyPrefixInfoPB)
		return
	}
	token := key[len(_keyPrefixTokenPB):]
	return d.TokenCache(c, token)
}

// InfoPBCache get info pb cache.
func (d *Dao) InfoPBCache(c context.Context, key string) (res *model.Info, err error) {
	if !strings.HasPrefix(key, _keyPrefixInfoPB) {
		err = fmt.Errorf("invalid cache key %s, key pattern is %s{mid}", key, _keyPrefixInfoPB)
		return
	}
	midStr := key[len(_keyPrefixInfoPB):]
	mid, err := strconv.ParseInt(midStr, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s, 10, 64) error(%v)", midStr, err)
		return
	}
	return d.InfoCache(c, mid)
}
