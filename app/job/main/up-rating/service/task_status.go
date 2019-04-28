package service

import (
	"context"
	"fmt"
)

// InsertTaskStatus insert task status
func (s *Service) InsertTaskStatus(c context.Context, ctype, status int, date, message string) (err error) {
	_, err = s.dao.InsertTaskStatus(c, fmt.Sprintf("(%d, %d, '%s', '%s')", ctype, status, date, message))
	return
}
