package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// BusinessList business list.
func (s *Service) BusinessList(c context.Context, pn, ps, status int) (res []*model.VipBusinessInfo, count int64, err error) {
	if count, err = s.dao.BussinessCount(c, status); err != nil {
		log.Error("%+v", err)
		return
	}
	if count == 0 {
		return
	}
	if res, err = s.dao.BussinessList(c, pn, ps, status); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

// UpdateBusinessInfo update business info.
func (s *Service) UpdateBusinessInfo(c context.Context, req *model.VipBusinessInfo) (err error) {
	var (
		business     *model.VipBusinessInfo
		appkey       string
		businessName string
	)
	if business, err = s.dao.SelBusiness(c, req.ID); err != nil {
		err = errors.WithStack(err)
		return
	}

	if business == nil {
		err = ecode.VipBusinessNotExitErr
		return
	}
	req.Secret = business.Secret

	appkey = business.AppKey
	businessName = business.BusinessName
	arg := new(model.QueryBusinessInfo)
	if appkey != req.AppKey {

		arg.Appkey = req.AppKey

		if business, err = s.dao.SelBusinessByQuery(c, arg); err != nil {
			err = errors.WithStack(err)
			return
		}
		if business != nil {
			err = ecode.VipAppkeyExitErr
			return
		}
		req.Secret = s.makeSecret(req.AppKey)
	}
	if businessName != req.BusinessName {
		arg = new(model.QueryBusinessInfo)
		arg.Name = req.BusinessName
		if business, err = s.dao.SelBusinessByQuery(c, arg); err != nil {
			err = errors.WithStack(err)
			return
		}
		if business != nil {
			err = ecode.VipBusinessNameExitErr
			return
		}
	}

	if _, err = s.dao.UpdateBusiness(c, req); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

// AddBusinessInfo add business info.
func (s *Service) AddBusinessInfo(c context.Context, req *model.VipBusinessInfo) (err error) {
	var (
		business *model.VipBusinessInfo
	)
	arg := new(model.QueryBusinessInfo)
	arg.Appkey = req.AppKey
	if business, err = s.dao.SelBusinessByQuery(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if business != nil {
		err = ecode.VipAppkeyExitErr
		return
	}
	arg = new(model.QueryBusinessInfo)
	arg.Name = req.BusinessName
	if business, err = s.dao.SelBusinessByQuery(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if business != nil {
		err = ecode.VipBusinessNameExitErr
		return
	}
	req.Secret = s.makeSecret(req.AppKey)
	if _, err = s.dao.AddBusiness(c, req); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}

func (s *Service) makeSecret(appkey string) string {
	key := fmt.Sprintf("%v,%v,%v", appkey, time.Now().UnixNano(), rand.Intn(100000000))
	hash := md5.New()
	hash.Write([]byte(key))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

// BusinessInfo business info
func (s *Service) BusinessInfo(c context.Context, id int) (r *model.VipBusinessInfo, err error) {
	if r, err = s.dao.SelBusiness(c, id); err != nil {
		log.Error("%+v", err)
		return
	}
	return
}
