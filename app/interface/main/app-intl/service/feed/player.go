package feed

import (
	"context"

	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

// ArchivesWithPlayer archives witch player
func (s *Service) ArchivesWithPlayer(c context.Context, aids []int64, qn int, platform string, fnver, fnval int) (res map[int64]*archive.ArchiveWithPlayer, err error) {
	if res, err = s.arc.ArchivesWithPlayer(c, aids, qn, platform, fnver, fnval); err != nil {
		log.Error("%+v", err)
	}
	if len(res) != 0 {
		return
	}
	am, err := s.arc.Archives(c, aids)
	if err != nil {
		return
	}
	if len(am) == 0 {
		return
	}
	res = make(map[int64]*archive.ArchiveWithPlayer, len(am))
	for aid, a := range am {
		res[aid] = &archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(a)}
	}
	return
}
