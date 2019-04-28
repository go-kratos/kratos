package service

import (
	"context"

	"go-common/app/service/openplatform/abtest/model"
)

//AddGroup add a new group
func (s *Service) AddGroup(c context.Context, g model.Group) (id int, err error) {
	return s.d.AddGroup(c, g)
}

//UpdateGroup update group by id
func (s *Service) UpdateGroup(c context.Context, g model.Group) (id int, err error) {
	return s.d.UpdateGroup(c, g)
}

//ListGroup list all groups
func (s *Service) ListGroup(c context.Context) (m []*model.Group, err error) {
	return s.d.ListGroup(c)
}

//DeleteGroup delete group by id
func (s *Service) DeleteGroup(c context.Context, id int) (r int, err error) {
	return s.d.DeleteGroup(c, id)
}
