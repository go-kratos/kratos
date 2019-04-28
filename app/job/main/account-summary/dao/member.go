package dao

// import (
// 	"context"
// 	"fmt"

// 	"go-common/app/job/main/account-summary/model"
// 	member "go-common/app/service/main/member/model"
// 	"go-common/library/log"

// 	"github.com/pkg/errors"
// )

// const (
// 	_AllBase     = `SELECT mid,name,sex,face,sign,rank FROM user_base_%02d`
// 	_AllExp      = `SELECT mid,exp FROM user_exp_%02d`
// 	_AllOfficial = `SELECT mid,role,title,description FROM user_official`
// )

// func (d *Dao) allMemberBaseFromTable(ctx context.Context, no int64) ([]*model.MemberBase, error) {
// 	rows, err := d.MemberDB.Query(ctx, fmt.Sprintf(_AllBase, no))
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}
// 	defer rows.Close()

// 	res := make([]*model.MemberBase, 0)
// 	for rows.Next() {
// 		mb := &member.BaseInfo{}
// 		if err = rows.Scan(&mb.Mid, &mb.Name, &mb.Sex, &mb.Face, &mb.Sign, &mb.Rank); err != nil {
// 			log.Error("Failed to scan row in query all member base: %+v", err)
// 			err = nil
// 			continue
// 		}
// 		mb.RandFaceURL()

// 		b := &model.MemberBase{
// 			EmbedMid: model.EmbedMid{Mid: mb.Mid},
// 			Birthday: model.Date(mb.Birthday),
// 			Face:     mb.Face,
// 			Name:     mb.Name,
// 			Rank:     mb.Rank,
// 			Sex:      mb.Sex,
// 			Sign:     mb.Sign,
// 		}
// 		res = append(res, b)
// 	}

// 	return res, nil
// }

// func (d *Dao) allMemberExpFromTable(ctx context.Context, no int64) ([]*model.MemberExp, error) {
// 	rows, err := d.MemberDB.Query(ctx, fmt.Sprintf(_AllExp, no))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	res := make([]*model.MemberExp, 0)
// 	for rows.Next() {
// 		mid := int64(0)
// 		exp := int64(0)
// 		if err = rows.Scan(&mid, &exp); err != nil {
// 			log.Error("Failed to scan row in query all member exp: %+v", err)
// 			err = nil
// 			continue
// 		}

// 		e := &model.MemberExp{
// 			EmbedMid: model.EmbedMid{Mid: mid},
// 			Exp:      exp / 100,
// 		}
// 		res = append(res, e)
// 	}

// 	return res, nil
// }

// // AllMemberBase is
// func (d *Dao) AllMemberBase(ctx context.Context) <-chan []*model.MemberBase {
// 	resCh := make(chan []*model.MemberBase)
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			res, err := d.allMemberBaseFromTable(context.Background(), int64(i))
// 			if err != nil {
// 				log.Error("Failed to get all member base from table with table id: %d: %+v", i, err)
// 				continue
// 			}
// 			resCh <- res
// 		}
// 		close(resCh)
// 	}()
// 	return resCh
// }

// // AllMemberExp is
// func (d *Dao) AllMemberExp(ctx context.Context) <-chan []*model.MemberExp {
// 	resCh := make(chan []*model.MemberExp)
// 	go func() {
// 		for i := 0; i < 100; i++ {
// 			res, err := d.allMemberExpFromTable(context.Background(), int64(i))
// 			if err != nil {
// 				log.Error("Failed to get all member exp from table with table id: %d: %+v", i, err)
// 				continue
// 			}
// 			resCh <- res
// 		}
// 		close(resCh)
// 	}()
// 	return resCh
// }

// // AllOfficial is
// func (d *Dao) AllOfficial(ctx context.Context) ([]*model.MemberOfficial, error) {
// 	rows, err := d.MemberDB.Query(ctx, _AllOfficial)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	res := make([]*model.MemberOfficial, 0)
// 	for rows.Next() {
// 		o := &model.MemberOfficial{}
// 		if err = rows.Scan(&o.Mid, &o.Role, &o.Title, &o.Description); err != nil {
// 			log.Error("Failed to scan row in query all member official: %+v", err)
// 			err = nil
// 			continue
// 		}
// 		res = append(res, o)
// 	}

// 	return res, nil
// }
