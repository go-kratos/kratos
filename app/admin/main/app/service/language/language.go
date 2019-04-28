package language

import (
	"context"

	"go-common/app/admin/main/app/conf"
	langdao "go-common/app/admin/main/app/dao/language"
	"go-common/app/admin/main/app/model/language"
	"go-common/library/log"
)

// Service language service
type Service struct {
	dao *langdao.Dao
}

// New new a language dao
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: langdao.New(c),
	}
	return
}

// Languages select all
func (s *Service) Languages(c context.Context) (res []*language.Language, err error) {
	if res, err = s.dao.Languages(c); err != nil {
		log.Error("s.dao.Languages error(%v)", err)
		return
	}
	return
}

// LangByID select by id
func (s *Service) LangByID(c context.Context, id int64) (res *language.Language, err error) {
	if res, err = s.dao.LangByID(c, id); err != nil {
		log.Error("s.dao.LangByID error(%v)", err)
		return
	}
	return
}

// Insert insert
func (s *Service) Insert(c context.Context, a *language.Param) (err error) {
	if err = s.dao.Insert(c, a); err != nil {
		log.Error("s.dao.InsertLanguage error(%v)", err)
		return
	}
	return
}

// Update update
func (s *Service) Update(c context.Context, a *language.Param) (err error) {
	if err = s.dao.Update(c, a); err != nil {
		log.Error("s.dao.UpdateLanguage(%v)", err)
		return
	}
	return
}
