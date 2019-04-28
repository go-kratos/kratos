package service

import (
	"context"

	"go-common/app/interface/main/space/model"
)

const _albumIndexPn = 0

var _emptyAlbumList = make([]*model.Album, 0)

// AlbumIndex get album index with mid.
func (s *Service) AlbumIndex(c context.Context, mid int64, ps int) (data []*model.Album, err error) {
	if data, err = s.dao.AlbumList(c, mid, _albumIndexPn, ps); err != nil {
		err = nil
		data = _emptyAlbumList
		return
	}
	if len(data) == 0 {
		data = _emptyAlbumList
	}
	return
}
