package service

import "context"

// UploadImg upload img
func (s *Service) UploadImg(c context.Context, imgContent []byte, imgName string) (location string, err error) {
	return s.dao.UploadImg(c, imgContent, imgName)
}
