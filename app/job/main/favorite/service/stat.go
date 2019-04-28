package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/favorite/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/log"
)

// consumeStat consumes folder's stat.
func (s *Service) consumeStat() {
	defer s.waiter.Done()
	for {
		select {
		case msg, ok := <-s.playStatSub.Messages():
			if !ok {
				break
			}
			msg.Commit()
			m := &model.PlayReport{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("fav json.Unmarshal(%s) error(%+v)", msg.Value, err)
				continue
			}
			if s.intercept(context.TODO(), m.ID, m.Mid, m.IP, m.Buvid) { // 防刷
				continue
			}
			if err := s.updatePlayStat(context.TODO(), m.ID); err != nil {
				log.Error("s.updatePlayStat(%d) error(%v)", m.ID, err)
			}
			log.Info("consumePlayStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		case msg, ok := <-s.favStatSub.Messages():
			if !ok {
				break
			}
			msg.Commit()
			m := &model.StatCount{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("fav json.Unmarshal(%s) error(%+v)", msg.Value, err)
				continue
			}
			if m.Type != "fav_playlist" {
				continue
			}
			if err := s.updateFavStat(context.TODO(), m.ID, m.Count); err != nil {
				log.Error("s.updateFav(%d,%d) error(%v)", m.ID, m.Count, err)
			}
			log.Info("consumeFavStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		case msg, ok := <-s.shareStatSub.Messages():
			if !ok {
				break
			}
			msg.Commit()
			m := &model.StatCount{}
			if err := json.Unmarshal(msg.Value, m); err != nil {
				log.Error("share json.Unmarshal(%s) error(%+v)", msg.Value, err)
				continue
			}
			if m.Type != "playlist" {
				continue
			}
			if err := s.updateShareStat(context.TODO(), m.ID, m.Count); err != nil {
				log.Error("s.statDao.UpdateShare(%d,%d) error(%v)", m.ID, m.Count, err)
			}
			log.Info("consumeShareStat key:%s partition:%d offset:%d msg: %+v)", msg.Key, msg.Partition, msg.Offset, m)
		}
	}
}

func (s *Service) intercept(c context.Context, id, mid int64, ip, buvid string) (ban bool) {
	if ban = s.statDao.IPBan(c, id, ip); ban {
		return
	}
	return s.statDao.BuvidBan(c, id, mid, ip, buvid)
}

func (s *Service) updatePlayStat(c context.Context, id int64) (err error) {
	f, err := s.stat(c, id)
	if err != nil {
		return
	}
	f.PlayCount = f.PlayCount + 1
	rows, err := s.statDao.UpdatePlay(context.TODO(), id, int64(f.PlayCount))
	if err != nil {
		log.Error("s.statDao.UpdatePlay(%d,%d) error(%v)", id, f.PlayCount, err)
		return
	}
	if rows > 0 {
		if err := s.statDao.SetFolderStatMc(c, id, f); err != nil {
			log.Error("s.SetFolderStatMc(%d,%+v) error(%v)", id, f, err)
		}
	}
	return
}

func (s *Service) updateFavStat(c context.Context, id, count int64) (err error) {
	f, err := s.stat(c, id)
	if err != nil {
		return
	}
	f.FavedCount = int32(count)
	rows, err := s.statDao.UpdateFav(context.TODO(), id, count)
	if err != nil {
		log.Error("s.statDao.UpdateFav(%d,%d) error(%v)", id, count, err)
		return
	}
	if rows > 0 {
		if err := s.statDao.SetFolderStatMc(c, id, f); err != nil {
			log.Error("s.SetFolderStatMc(%d,%+v) error(%v)", id, f, err)
		}
	}
	return
}

func (s *Service) updateShareStat(c context.Context, id, count int64) (err error) {
	f, err := s.stat(c, id)
	if err != nil {
		return
	}
	f.ShareCount = int32(count)
	rows, err := s.statDao.UpdateShare(context.TODO(), id, count)
	if err != nil {
		log.Error("s.statDao.UpdateShare(%d,%d) error(%v)", id, count, err)
		return
	}
	if rows > 0 {
		if err := s.statDao.SetFolderStatMc(c, id, f); err != nil {
			log.Error("s.SetFolderStatMc(%d,%+v) error(%v)", id, f, err)
		}
	}
	return
}

func (s *Service) stat(c context.Context, id int64) (f *favmdl.Folder, err error) {
	if f, err = s.statDao.FolderStatMc(c, id); err != nil {
		log.Error("s.statDao.FolderStatMc(%d) error(%v)", id, err)
		return
	}
	if f != nil {
		return
	}
	if f, err = s.statDao.Stat(c, id); err != nil {
		log.Error("s.statDao.FolderStatMc(%d) error(%v)", id, err)
		return
	}
	if f == nil {
		f = new(favmdl.Folder)
	}
	if err := s.statDao.SetFolderStatMc(c, id, f); err != nil {
		log.Error("s.statDao.SetFolderStatMc(%d,%+v) error(%v)", id, f, err)
	}
	return
}
