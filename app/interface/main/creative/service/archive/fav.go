package archive

import (
	"context"
	"sort"
	"strconv"

	"go-common/app/interface/main/creative/model/archive"
)

// get max 5 fav types
func (s *Service) favTypes(c context.Context, mid int64) (favTps []*archive.Type) {
	favTps = make([]*archive.Type, 0)
	var res map[string]int64
	res, _ = s.arc.FavTypes(c, mid)
	if len(res) > 0 {
		type kv struct {
			TidStr    string
			Timestamp int64
		}
		var kvSlice []kv
		for k, v := range res {
			kvSlice = append(kvSlice, kv{k, v})
		}
		sort.Slice(kvSlice, func(i, j int) bool {
			return kvSlice[i].Timestamp > kvSlice[j].Timestamp
		})
		for _, v := range kvSlice {
			tid, _ := strconv.Atoi(v.TidStr)
			if tp, ok := s.p.TypeMapCache[int16(tid)]; ok && len(favTps) < 5 {
				favTps = append(favTps, tp)
			}
		}
	}
	return
}

// Fav fn
func (s *Service) Fav(c context.Context, mid int64) (res map[string]interface{}) {
	res = make(map[string]interface{})
	res["typelist"] = s.favTypes(c, mid)
	return
}
