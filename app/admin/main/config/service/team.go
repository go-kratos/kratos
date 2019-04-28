package service

import (
	"go-common/app/admin/main/config/model"
	"go-common/library/log"
)

// CreateTeam create App.
func (s *Service) CreateTeam(name, env, zone string) error {
	app := &model.Team{Name: name, Env: env, Zone: zone}
	return s.dao.DB.Create(app).Error
}

//TeamByName get team by Name.
func (s *Service) TeamByName(name, env, zone string) (team *model.Team, err error) {
	team = &model.Team{}
	if err = s.dao.DB.FirstOrCreate(&team, &model.Team{Name: name, Env: env, Zone: zone}).Error; err != nil {
		log.Error("TeamByName() error(%v)", err)
	}
	return
}
