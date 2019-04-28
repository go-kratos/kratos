package service

import (
	"go-common/app/job/main/block/conf"
)

// MSGRemoveInfo get msg info
func (s *Service) MSGRemoveInfo() (code string, title, content string) {
	code = conf.Conf.Property.MSG.BlockRemove.Code
	title = conf.Conf.Property.MSG.BlockRemove.Title
	content = conf.Conf.Property.MSG.BlockRemove.Content
	return
}
