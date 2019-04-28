package service

import (
	"go-common/app/interface/main/tv/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_pgcPiece = 50
)

// PgcCards treat the slice of int, and call the PGC api
func (s *Service) PgcCards(ids []int64) (res map[int64]string, err error) {
	var (
		length   = len(ids)
		nbBatch  = length / _pgcPiece // number of batch to run
		resCards map[string]*model.SeasonCard
	)
	if length%_pgcPiece != 0 {
		nbBatch = nbBatch + 1
	}
	res = make(map[int64]string)
	for i := 0; i < nbBatch; i++ {
		begin := i * _pgcPiece
		end := (i + 1) * _pgcPiece
		if end > length {
			end = length
		}
		batchIDs := ids[begin:end]
		batchStr := xstr.JoinInts(batchIDs)
		if resCards, err = s.dao.PgcCards(ctx, batchStr); err != nil {
			log.Error("PGCCards IDs: %s, Err: %v", batchStr, err)
			return
		}
		for k, v := range resCards {
			if sid := int64(atoi(k)); sid > 0 && v.NewEP != nil && v.NewEP.IndexShow != "" {
				res[sid] = v.NewEP.IndexShow
			}
		}
	}
	return
}
