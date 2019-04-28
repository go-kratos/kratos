package service

import (
	"context"

	"go-common/app/admin/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Privileges Privilege list.
func (s *Service) Privileges(c context.Context, langType int8) (res []*model.PrivilegeResp, err error) {
	var (
		ps  []*model.Privilege
		prs []*model.PrivilegeResources
	)
	if ps, err = s.dao.PrivilegeList(c, langType); err != nil {
		err = errors.WithStack(err)
		return
	}
	if prs, err = s.dao.PrivilegeResourcesList(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	prmap := make(map[int64]map[int8]*model.PrivilegeResources, len(prs))
	for _, v := range prs {
		if prmap[v.PID] == nil {
			prmap[v.PID] = map[int8]*model.PrivilegeResources{}
		}
		prmap[v.PID][v.Type] = v
	}
	for i, v := range ps {
		r := &model.PrivilegeResp{
			ID:          v.ID,
			Name:        v.Name,
			Title:       v.Title,
			Explain:     v.Explain,
			Type:        v.Type,
			Operator:    v.Operator,
			State:       v.State,
			IconURL:     v.IconURL,
			IconGrayURL: v.IconGrayURL,
			LangType:    v.LangType,
			Order:       int64(i),
		}
		if prmap[v.ID] != nil && prmap[v.ID][model.WebResources] != nil {
			r.WebLink = prmap[v.ID][model.WebResources].Link
			r.WebImageURL = prmap[v.ID][model.WebResources].ImageURL
		}
		if prmap[v.ID] != nil && prmap[v.ID][model.AppResources] != nil {
			r.AppLink = prmap[v.ID][model.AppResources].Link
			r.AppImageURL = prmap[v.ID][model.AppResources].ImageURL
		}
		res = append(res, r)
	}
	return
}

// UpdatePrivilegeState update privilege state.
func (s *Service) UpdatePrivilegeState(c context.Context, p *model.Privilege) (err error) {
	if _, err = s.dao.UpdateStatePrivilege(c, p); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// DeletePrivilege delete privilege .
func (s *Service) DeletePrivilege(c context.Context, id int64) (err error) {
	if _, err = s.dao.DeletePrivilege(c, id); err != nil {
		err = errors.WithStack(err)
	}
	return
}

// AddPrivilege add privilege.
func (s *Service) AddPrivilege(c context.Context, req *model.ArgAddPrivilege, img *model.ArgImage) (err error) {
	var (
		iconURL     string
		iconGrayURL string
		order       int64
		webImageURL string
		appImageURL string
		id          int64
		tx          *gorm.DB
	)
	if iconURL, iconGrayURL, webImageURL, appImageURL, err = s.uploadImage(c, img); err != nil {
		err = errors.WithStack(err)
		return
	}
	if iconURL == "" || iconGrayURL == "" {
		err = ecode.VipFileUploadFaildErr
		return
	}
	if order, err = s.dao.MaxOrder(c); err != nil {
		err = errors.WithStack(err)
		return
	}
	order++
	tx = s.dao.BeginGormTran(c)
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				err = errors.Wrapf(err, "tx.Rollback(%+v)", req)
			}
			return
		}
		if err = tx.Commit().Error; err != nil {
			err = errors.Wrapf(err, "tx.Commit(%+v)", req)
		}
	}()
	if id, err = s.dao.AddPrivilege(tx, &model.Privilege{
		Name:        req.Name,
		Title:       req.Title,
		Explain:     req.Explain,
		Type:        req.Type,
		Operator:    req.Operator,
		State:       model.NormalPrivilege,
		IconURL:     iconURL,
		IconGrayURL: iconGrayURL,
		Order:       order,
		LangType:    req.LangType,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = s.dao.AddPrivilegeResources(tx, &model.PrivilegeResources{
		PID:      id,
		Link:     req.AppLink,
		ImageURL: appImageURL,
		Type:     model.AppResources,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = s.dao.AddPrivilegeResources(tx, &model.PrivilegeResources{
		PID:      id,
		Link:     req.WebLink,
		ImageURL: webImageURL,
		Type:     model.WebResources,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdatePrivilege update privilege.
func (s *Service) UpdatePrivilege(c context.Context, req *model.ArgUpdatePrivilege, img *model.ArgImage) (err error) {
	var (
		iconURL     string
		iconGrayURL string
		webImageURL string
		appImageURL string
		tx          *gorm.DB
	)
	if iconURL, iconGrayURL, webImageURL, appImageURL, err = s.uploadImage(c, img); err != nil {
		err = errors.WithStack(err)
		return
	}
	tx = s.dao.BeginGormTran(c)
	defer func() {
		if err != nil {
			if err1 := tx.Rollback().Error; err1 != nil {
				err = errors.Wrapf(err, "tx.Rollback(%+v)", req)
			}
			return
		}
		if err = tx.Commit().Error; err != nil {
			err = errors.Wrapf(err, "tx.Commit(%+v)", req)
		}
	}()
	if _, err = s.dao.UpdatePrivilege(tx, &model.Privilege{
		ID:          req.ID,
		Name:        req.Name,
		Title:       req.Title,
		Explain:     req.Explain,
		Type:        req.Type,
		Operator:    req.Operator,
		State:       model.NormalPrivilege,
		IconURL:     iconURL,
		IconGrayURL: iconGrayURL,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = s.dao.UpdatePrivilegeResources(tx, &model.PrivilegeResources{
		PID:      req.ID,
		Link:     req.AppLink,
		ImageURL: appImageURL,
		Type:     model.AppResources,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = s.dao.UpdatePrivilegeResources(tx, &model.PrivilegeResources{
		PID:      req.ID,
		Link:     req.WebLink,
		ImageURL: webImageURL,
		Type:     model.WebResources,
	}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (s *Service) uploadImage(c context.Context, req *model.ArgImage) (iconURL, iconGrayURL, webImageURL, appImageURL string, err error) {
	var g errgroup.Group
	if req.IconFileType == "" &&
		req.IconGrayFileType == "" &&
		req.WebImageFileType == "" &&
		req.AppImageFileType == "" {
		return
	}
	if req.IconFileType != "" {
		g.Go(func() (err error) {
			//update icon.
			if iconURL, err = s.dao.Upload(c, "", req.IconFileType, req.IconBody, s.c.Bfs); err != nil {
				log.Error("d.Upload iconURL(%+v) error(%v)", req, err)
			}
			return
		})
	}
	if req.IconGrayFileType != "" {
		g.Go(func() (err error) {
			//update gray icon.
			if iconGrayURL, err = s.dao.Upload(c, "", req.IconGrayFileType, req.IconGrayBody, s.c.Bfs); err != nil {
				log.Error("d.Upload iconGrayURL(%+v) error(%v)", req, err)
			}
			return
		})
	}
	if req.WebImageFileType != "" {
		g.Go(func() (err error) {
			if webImageURL, err = s.dao.Upload(c, "", req.WebImageFileType, req.WebImageBody, s.c.Bfs); err != nil {
				log.Error("d.Upload webImageURL(%+v) error(%v)", req, err)
			}
			return
		})
	}
	if req.AppImageFileType != "" {
		g.Go(func() (err error) {
			if appImageURL, err = s.dao.Upload(c, "", req.AppImageFileType, req.AppImageBody, s.c.Bfs); err != nil {
				log.Error("d.Upload appImageURL(%+v) error(%v)", req, err)
			}
			return
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("g.Wait error(%v)", err)
		return
	}
	return
}

// UpdateOrder update order.
func (s *Service) UpdateOrder(c context.Context, a *model.ArgOrder) (err error) {
	if _, err = s.dao.UpdateOrder(c, a.AID, a.BID); err != nil {
		err = errors.WithStack(err)
	}
	return
}
