package service

import (
	"context"
	http "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
)

// GetLocaitonAll .
func (s *Service) GetLocaitonAll(c context.Context, arg *http.LocationRequest) (*http.LocationResponse, error) {
	result := &http.LocationResponse{}
	m, err := s.dao.GetLocationAll(c)
	if err != nil {
		return result, err
	}

	var coutries []*model.Location
	for _, item := range (*m)[arg.PID] {
		coutry := &model.Location{
			ID:   item.ID,
			PID:  item.PID,
			Name: item.Name,
		}
		var provices []*model.Location
		for _, v := range (*m)[item.ID] {
			provice := &model.Location{
				ID:   v.ID,
				PID:  v.PID,
				Name: v.Name,
			}
			var citys []*model.Location
			for _, u := range (*m)[v.ID] {
				city := &model.Location{
					ID:   u.ID,
					PID:  u.PID,
					Name: u.Name,
				}
				var area []*model.Location
				for _, w := range (*m)[u.ID] {
					var child []*model.Location
					area = append(area, &model.Location{
						ID:    w.ID,
						PID:   w.PID,
						Name:  w.Name,
						Child: child,
					})
				}
				city.Child = area
				citys = append(citys, city)
			}
			provice.Child = citys
			provices = append(provices, provice)
		}
		coutry.Child = provices
		coutries = append(coutries, coutry)
	}

	result.List = coutries

	return result, err
}

// GetLocationChild .
func (s *Service) GetLocationChild(c context.Context, arg *http.LocationRequest) (*http.LocationResponse, error) {
	result := &http.LocationResponse{}
	m, err := s.dao.GetLocationChild(c, arg.PID)
	if err != nil {
		return result, err
	}

	var provices []*model.Location
	for _, v := range (*m)[arg.PID] {
		var child []*model.Location
		provice := &model.Location{
			ID:    v.ID,
			PID:   v.PID,
			Name:  v.Name,
			Child: child,
		}
		provices = append(provices, provice)
	}

	result.List = provices

	return result, err
}
