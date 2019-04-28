package service

import (
	"go-common/app/admin/ep/merlin/model"
)

// QueryImage query image.
func (s *Service) QueryImage() ([]*model.Image, error) {
	return s.dao.Images()
}

// AddImage add image.
func (s *Service) AddImage(image *model.Image) error {
	return s.dao.AddImage(image)
}

// UpdateImage update image.
func (s *Service) UpdateImage(image *model.Image) error {
	return s.dao.UpdateImage(image)
}

// DeleteImage delete image.
func (s *Service) DeleteImage(iID int64) error {
	return s.dao.DelImage(iID)
}
