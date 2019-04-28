package global

import (
	"context"

	accmdl "go-common/app/service/main/account/model"
	"go-common/library/log"
)

// MID gets mid by nickname
func MID(c context.Context, nickname string) (mid int64, err error) {
	res, err := GetAccRPC().InfosByName3(c, &accmdl.ArgNames{Names: []string{nickname}})
	if err != nil {
		log.Error("InfosByName3 fail, nickname=%+v, err=%+v", nickname, err)
		return
	}
	for k := range res {
		mid = k
	}
	return
}

// Names get nicknames by mids
func Names(c context.Context, mids []int64) (res map[int64]string, err error) {
	res = make(map[int64]string)
	infos, err := GetAccRPC().Infos3(c, &accmdl.ArgMids{Mids: mids})
	if err != nil {
		log.Error("Infos3 fail, mids=%+v, err=%+v", mids, err)
		return
	}
	for k, v := range infos {
		res[k] = v.Name
	}
	return
}
