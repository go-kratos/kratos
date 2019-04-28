package service

import (
	"context"
	v1 "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"sync"
	"time"
)

var (
	_wait sync.WaitGroup

	_archiveSubKeepAlive = true
)

// archiveSub 订阅B站稿件，过滤后传给生产流程
func (s *Service) archiveSub(c context.Context) {
	_wait.Add(1)
	defer _wait.Done()
	for _archiveSubKeepAlive {
		// 订阅稿件
		archive, err := s.dao.ArchiveSub()
		if err == ecode.ArchiveDatabusNilErr {
			log.Error("ArchiveSub failed archive[%v] err[%v]", archive, err)
			return
		}
		if err != nil {
			log.Error("ArchiveSub failed archive[%v] err[%v]", archive, err)
			continue
		}
		if archive == nil {
			continue
		}

		// 生产SVID
		pubtime, err := time.Parse("2006-01-02 15:04:05", archive.PubTime)
		if err != nil {
			log.Error("ArchiveSub pubtime parse failed archive[%v] err[%v]", archive, err)
			continue
		}
		res, err := s.CreateID(c, &v1.CreateIDRequest{
			Mid:  int64(archive.MID),
			Time: pubtime.Unix(),
		})
		if err != nil || res == nil {
			log.Error("ArchiveSub CreateID failed archive[%v] err[%v]", archive, err)
			continue
		}

		// 传递生产
		if err = s.dao.ArchiveKickOff(c, res.NewId, archive); err != nil {
			log.Error("ArchiveSub ArchiveKickOff failed archive[%v] SVID[%d] err[%v]", archive, res.NewId, err)
			continue
		}
	}
}

func (s *Service) archiveSubClose() {
	_archiveSubKeepAlive = false
	_wait.Wait()
}
