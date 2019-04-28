package service

import (
	"go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
)

// CheckBusiness check if exist the business id.
func (s *Service) CheckBusiness(business string) (id int64, err error) {
	b := s.dao.BusinessMap[business]
	if b == nil {
		err = ecode.ThumbupBusinessBlankErr
		return
	}
	id = b.ID
	return
}

// CheckBusinessOrigin check origin id.
func (s *Service) CheckBusinessOrigin(business string, originID int64) (id int64, err error) {
	b := s.dao.BusinessMap[business]
	if b == nil {
		err = ecode.ThumbupBusinessBlankErr
		return
	}
	if (b.EnableOriginID == 1 && originID == 0) || (b.EnableOriginID == 0 && originID != 0) {
		err = ecode.ThumbupOriginErr
		return
	}
	id = b.ID
	return
}

func (s *Service) checkItemLikeType(businessID int64, state int8) error {
	if state == model.StateLike {
		if s.dao.BusinessIDMap[businessID] == nil {
			return ecode.ThumbupBusinessBlankErr
		}
		if !s.dao.BusinessIDMap[businessID].EnableItemLikeList() {
			return ecode.ThumbupBusinessErr
		}
		return nil
	}
	// check dislikes
	if s.dao.BusinessIDMap[businessID] == nil {
		return ecode.ThumbupBusinessBlankErr
	}
	if !s.dao.BusinessIDMap[businessID].EnableItemDislikeList() {
		return ecode.ThumbupBusinessErr
	}
	return nil
}

func (s *Service) checkUserLikeType(businessID int64, state int8) error {
	if state == model.StateLike {
		if s.dao.BusinessIDMap[businessID] == nil {
			return ecode.ThumbupBusinessBlankErr
		}
		if !s.dao.BusinessIDMap[businessID].EnableUserLikeList() {
			return ecode.ThumbupBusinessErr
		}
		return nil
	}
	// check dislikes
	if s.dao.BusinessIDMap[businessID] == nil {
		return ecode.ThumbupBusinessBlankErr
	}
	if !s.dao.BusinessIDMap[businessID].EnableUserDislikeList() {
		return ecode.ThumbupBusinessErr
	}
	return nil
}
