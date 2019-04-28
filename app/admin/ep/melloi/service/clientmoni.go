package service

import "go-common/app/admin/ep/melloi/model"

//AddClientMoni add ClientMoni
func (s *Service) AddClientMoni(clm *model.ClientMoni) (int, error) {
	return s.dao.AddClientMoni(clm)
}
