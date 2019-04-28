package service

import (
	"context"

	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// favoriteds return aids is favs.
func (s *Service) favoriteds(c context.Context, mid int64, aids []int64) (res map[int64]bool) {
	var n = 50
	res = make(map[int64]bool, len(aids))
	for len(aids) > 0 {
		if n > len(aids) {
			n = len(aids)
		}
		arg := &favmdl.ArgIsFavs{
			Type:   favmdl.TypeVideo,
			Mid:    mid,
			Oids:   aids[:n],
			RealIP: metadata.String(c, metadata.RemoteIP),
		}
		favMap, err := s.favRPC.IsFavs(c, arg)
		if err != nil {
			log.Error("s.favRPC.IsFavs(%v) error(%v)", arg, err)
			return
		}
		aids = aids[n:]
		for k, v := range favMap {
			res[k] = v
		}
	}
	return
}
