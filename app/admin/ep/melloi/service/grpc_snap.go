package service

import "go-common/app/admin/ep/melloi/model"

// QueryGRPCSnapByID  query grpcsnap by id
func (s *Service) QueryGRPCSnapByID(id int) (*model.GRPCSnap, error) {
	return s.dao.QueryGRPCSnapByID(id)
}

// UpdateGRPCSnap  update grpc snap
func (s *Service) UpdateGRPCSnap(grpcSnap *model.GRPCSnap) error {
	return s.dao.UpdateGRPCSnap(grpcSnap)
}

// CreateGRPCSnap Create GRPC Snap
func (s *Service) CreateGRPCSnap(grpcSnap *model.GRPCSnap) error {
	return s.dao.CreateGRPCSnap(grpcSnap)
}

// DeleteGRPCSnap Delete GRPC Snap
func (s *Service) DeleteGRPCSnap(id int) error {
	return s.dao.DeleteGRPCSnap(id)
}
