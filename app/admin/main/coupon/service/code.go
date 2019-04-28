package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CodePage code page.
func (s *Service) CodePage(c context.Context, a *model.ArgCouponCode) (res *model.CodePage, err error) {
	res = new(model.CodePage)
	var (
		count int64
		b     *model.CouponBatchInfo
		now   = time.Now().Unix()
	)
	if count, err = s.dao.CountCode(c, a); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	if res.CodeList, err = s.dao.CodeList(c, a); err != nil {
		return
	}
	for _, v := range res.CodeList {
		if b, err = s.BatchInfo(c, v.BatchToken); err != nil {
			return
		}
		if b.ExpireDay == -1 && b.ExpireTime <= now {
			v.State = model.CodeStateExpire
		}
	}
	return
}

// CodeBlock code block.
func (s *Service) CodeBlock(c context.Context, a *model.ArgCouponCode) (err error) {
	var code *model.CouponCode
	if code, err = s.dao.CodeByID(c, a.ID); err != nil {
		return
	}
	if code == nil || code.State != model.CodeStateNotUse {
		return ecode.RequestErr
	}
	return s.dao.UpdateCodeBlock(c, &model.CouponCode{ID: a.ID, State: model.CodeStateBlock, Ver: code.Ver})
}

// CodeUnBlock code un block.
func (s *Service) CodeUnBlock(c context.Context, a *model.ArgCouponCode) (err error) {
	var code *model.CouponCode
	if code, err = s.dao.CodeByID(c, a.ID); err != nil {
		return
	}
	if code == nil || code.State != model.CodeStateBlock {
		return ecode.RequestErr
	}
	return s.dao.UpdateCodeBlock(c, &model.CouponCode{ID: a.ID, State: model.CodeStateNotUse, Ver: code.Ver})
}

// ExportCode export code.
func (s *Service) ExportCode(c context.Context, a *model.ArgCouponCode) (res []string, err error) {
	if a.BatchToken == "" {
		return nil, ecode.RequestErr
	}
	var cs []*model.CouponCode
	a.Ps = 1000
	a.Pn = 1
	for {
		if cs, err = s.dao.CodeList(context.Background(), a); err != nil {
			return
		}
		if len(cs) == 0 {
			break
		}
		for _, v := range cs {
			res = append(res, v.Code)
		}
		a.Pn++
	}
	return
}

// InitCodes init codes
func (s *Service) InitCodes(c context.Context, token string) (err error) {
	var info *model.CouponBatchInfo
	if info, err = s.dao.BatchInfo(c, token); err != nil {
		return
	}
	if info == nil {
		return
	}
	if info.MaxCount == 0 || info.MaxCount > model.BatchCodeMaxCount {
		return
	}
	// init code.
	go func() {
		log.Info("init code start arg[%s,%d]", token, info.MaxCount)
		cs := []*model.CouponCode{}
		for i := int64(0); i < info.MaxCount; i++ {
			cs = append(cs, &model.CouponCode{
				BatchToken: token,
				State:      model.CodeStateNotUse,
				Code:       codeToken(i),
				CouponType: model.CouponAllowance,
			})
			if len(cs) == int(info.MaxCount) || len(cs)%model.BatchAddCodeSlice == 0 {
				if err = s.dao.BatchAddCode(context.Background(), cs); err != nil {
					log.Error("init code error[%s,%v]", token, err)
					return
				}
				log.Info("init code ing arg[%s,%d,%d]", token, info.MaxCount, i)
				time.Sleep(50 * time.Millisecond)
				cs = []*model.CouponCode{}
			}
		}
		log.Info("init code end arg[%s,%d]", token, info.MaxCount)
	}()
	return
}

func codeToken(i int64) string {
	hash := md5.New()
	unix := time.Now().UnixNano()
	key := fmt.Sprintf("%v,%v,%v", unix, i, rand.Intn(100000000))
	hash.Write([]byte(key))
	sum := hash.Sum(nil)
	code := hex.EncodeToString(sum)
	return code[12:24]
}
