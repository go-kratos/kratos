package resource

import (
	"context"

	accgrpc "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// MidByNickname .
func MidByNickname(c context.Context, name string) (mid int64, err error) {
	reply, err := accCli.InfosByName3(c, &accgrpc.NamesReq{Names: []string{name}})
	if err != nil || reply == nil {
		log.Error("accCli.InfosByName3 name(%s) err(%v)", name, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if reply.Infos != nil {
		for k := range reply.Infos {
			mid = k
		}
	}
	return
}

// NamesByMIDs .
func NamesByMIDs(c context.Context, mids []int64) (res map[int64]string, err error) {
	reply, err := accCli.Infos3(c, &accgrpc.MidsReq{Mids: mids})
	if err != nil || reply == nil {
		log.Error("accCli.NamesByMIDs mids(%v) err(%v)", mids, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if reply.Infos != nil {
		res = make(map[int64]string, len(reply.Infos))
		for mid, info := range reply.Infos {
			res[mid] = info.Name
		}
	}
	return
}

// NameByMID .
func NameByMID(c context.Context, mid int64) (nickname string, err error) {
	reply, err := accCli.Info3(c, &accgrpc.MidReq{Mid: mid})
	if err != nil || reply == nil {
		log.Error("accCli.Info3 mid(%d) err(%v)", mid, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if reply.Info != nil {
		nickname = reply.Info.Name
	}
	return
}
