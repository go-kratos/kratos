package service

import (
	"context"
	"sort"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/track"
	"go-common/library/log"
)

const (
	historyTimeMin = 0
	historyTimeMax = 0x7fffffffffffffff
)

// TrackArchive get archive list.
func (s *Service) TrackArchive(c context.Context, aid int64) (archive []*track.Archive, err error) {
	if archive, err = s.track.ArchiveTrack(c, aid); err != nil {
		log.Error("s.track.ArchiveTrack(%d) error(%v)", aid, err)
		return
	}
	sort.Sort(track.Archives(archive))
	return
}

//TrackArchiveInfo 稿件信息追踪
func (s *Service) TrackArchiveInfo(c context.Context, aid int64) (info *track.ArcTrackInfo, err error) {
	relation := [][]int{}
	editHistory, err := s.AllEditHistory(c, aid)
	if err != nil {
		log.Error("TrackArchiveInfo s.AllEditHistory(aid(%d)) error(%v)", aid, err)
		return
	}

	tr, err := s.TrackArchive(c, aid)
	if err != nil {
		log.Error("TrackArchiveInfo s.TrackArchive(aid(%d)) error(%v)", aid, err)
		return
	}

	//2个降序数组，track根据history聚合
	cpHistoryTime := []int64{historyTimeMax}
	for _, h := range editHistory {
		cpHistoryTime = append(cpHistoryTime, int64(h.ArcHistory.CTime))
	}
	cpHistoryTime = append(cpHistoryTime, historyTimeMin)
	index := 0
	histlen := len(cpHistoryTime) - 1
	for i := 1; i <= histlen; i++ {
		rela := []int{}
		for ; index < len(tr); index++ {
			t := int64(tr[index].Timestamp)
			if t >= cpHistoryTime[i] && t < cpHistoryTime[i-1] {
				rela = append(rela, index)
				continue
			}

			break
		}

		if i == histlen && len(rela) == 0 {
			continue
		}
		relation = append(relation, rela)
	}

	info = &track.ArcTrackInfo{
		EditHistory: editHistory,
		Track:       tr,
		Relation:    relation,
	}
	return
}

//TrackHistoryDetail 稿件某条编辑历史的详细情况
func (s *Service) TrackHistoryDetail(c context.Context, hid int64) (h *archive.EditHistory, err error) {
	if h, err = s.EditHistory(c, hid); err != nil {
		log.Error("TrackHistoryDetail s.EditHistory(hid(%d)) error(%v)", hid, err)
		return
	}
	if h == nil {
		return
	}

	//获取最新的src_type
	cids := []int64{}
	for _, vh := range h.VHistory {
		cids = append(cids, vh.CID)
	}
	srcTypes, err := s.arc.VideoSrcTypeByIDs(c, cids)
	if err != nil {
		log.Error("TrackHistoryDetail s.arc.VideoSrcTypeByID(hid(%d)) error(%v)", hid, err)
		return
	}
	for _, vh := range h.VHistory {
		vh.SRCType = srcTypes[vh.CID]
	}

	return
}

// TrackVideo get video process.
func (s *Service) TrackVideo(c context.Context, filename string, aid int64) (video []*track.Video, err error) {
	if video, err = s.track.VideoTrack(c, filename, aid); err != nil {
		log.Error("s.track.VideoTrack(%s) error(%v)", filename, err)
		return
	}

	//以下努力均为获取当时视频审核的属性
	vid, err := s.arc.VIDByAIDFilename(c, aid, filename)
	if err != nil {
		log.Error("s.arc.VIDByAIDFilename error(%v), filename(%s)", err, filename)
		return
	}

	if vid <= 0 {
		return
	}
	attrs, ctimes, err := s.arc.VideoOperAttrsCtimes(c, vid)
	if err != nil {
		log.Error("s.arc.VideoOperAttrsCtimes error(%v) vid(%d)", err, vid)
		return
	}

	for _, tk := range video {
		i := 0
		for ; i < len(attrs); i++ {
			if ctimes[i] > int64(tk.Timestamp) {
				continue
			}

			tk.Attribute = attrs[i]
			break
		}
	}
	return
}
