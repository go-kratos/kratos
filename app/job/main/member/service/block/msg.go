package block

// MSGRemoveInfo get msg info
func (s *Service) MSGRemoveInfo() (code string, title, content string) {
	code = s.conf.BlockProperty.MSG.BlockRemove.Code
	title = s.conf.BlockProperty.MSG.BlockRemove.Title
	content = s.conf.BlockProperty.MSG.BlockRemove.Content
	return
}
