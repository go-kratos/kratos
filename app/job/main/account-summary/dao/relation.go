package dao

// import (
// 	"context"
// 	"fmt"

// 	"go-common/app/job/main/account-summary/model"
// 	"go-common/library/log"
// )

// var (
// 	_AllStat = `SELECT mid,following,whisper,black,follower FROM user_relation_stat_%02d`
// )

// func (d *Dao) allRelationStatFromTable(ctx context.Context, no int64) ([]*model.RelationStat, error) {
// 	rows, err := d.RelationDB.Query(ctx, fmt.Sprintf(_AllStat, no))
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := make([]*model.RelationStat, 0)
// 	defer rows.Close()
// 	for rows.Next() {
// 		rs := &model.RelationStat{}
// 		if err = rows.Scan(&rs.Mid, &rs.Following, &rs.Whisper, &rs.Black, &rs.Follower); err != nil {
// 			log.Error("Failed to scan row in query all relation stat: %+v", err)
// 			err = nil
// 			continue
// 		}
// 		res = append(res, rs)
// 	}

// 	return res, nil
// }

// // AllRelationStat is
// func (d *Dao) AllRelationStat(ctx context.Context) <-chan []*model.RelationStat {
// 	resCh := make(chan []*model.RelationStat)
// 	go func() {
// 		for i := 0; i < 50; i++ {
// 			res, err := d.allRelationStatFromTable(ctx, int64(i))
// 			if err != nil {
// 				log.Error("Failed to get all relation stat from table with table id: %d: %+v", i, err)
// 				continue
// 			}
// 			resCh <- res
// 		}
// 		close(resCh)
// 	}()
// 	return resCh
// }
