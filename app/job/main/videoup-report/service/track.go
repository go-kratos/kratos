package service

import (
	"context"
	"encoding/json"
	"strconv"

	"go-common/app/job/main/videoup-report/model/archive"
)

func (s *Service) trackArchive(nw *archive.Archive, old *archive.Archive) (err error) {
	var (
		bs      []byte
		remarks = make(map[string]string)
	)
	if addit, _ := s.arc.Addit(context.TODO(), nw.ID); addit != nil {
		remarks["dynamic"] = addit.Dynamic
		if addit.MissionID > 0 {
			remarks["mission_id"] = strconv.FormatInt(addit.MissionID, 10)
		}
	}
	if old == nil {
		remarks["cover"] = nw.Cover
		remarks["desc"] = nw.Content
		remarks["title"] = nw.Title
		remarks["typeid"] = strconv.Itoa(int(nw.TypeID))
		remarks["copyright"] = strconv.Itoa(int(nw.Copyright))
		bs, _ = json.Marshal(remarks)
	} else if nw.State != old.State || nw.Access != old.Access || nw.Round != old.Round || nw.Content != old.Content ||
		nw.Cover != old.Cover || nw.Title != old.Title || nw.TypeID != old.TypeID || nw.Copyright != old.Copyright || nw.Attribute != old.Attribute {
		if nw.Cover != old.Cover {
			remarks["cover"] = nw.Cover
		}
		if nw.Content != old.Content {
			remarks["desc"] = nw.Content
		}
		if nw.Title != old.Title {
			remarks["title"] = nw.Title
		}
		if nw.TypeID != old.TypeID {
			remarks["typeid"] = strconv.Itoa(int(nw.TypeID))
		}
		if nw.Copyright != old.Copyright {
			remarks["copyright"] = strconv.Itoa(int(nw.Copyright))
		}
		if len(remarks) != 0 {
			bs, _ = json.Marshal(remarks)
		}
		if nw.State >= int(archive.StateOpen) && nw.Access == int(archive.AccessMember) {
			nw.State = int(archive.AccessMember)
		}
	} else {
		// NOTE: nothing modify
		return
	}
	s.arc.AddTrack(context.TODO(), nw.ID, nw.State, nw.Round, nw.Attribute, string(bs), nw.MTime, nw.MTime)
	return
}

func (s *Service) trackVideo(nw *archive.Video, old *archive.Video) (err error) {
	var (
		remarks = make(map[string]interface{})
		bs      []byte
	)
	if old == nil {
		if nw.Title != "" {
			remarks["title"] = nw.Title
		}
		if nw.Desc != "" {
			remarks["desc"] = nw.Desc
		}
	} else if nw.XcodeState != old.XcodeState || nw.Status != old.Status || nw.Title != old.Title || nw.Desc != old.Desc || nw.Attribute != old.Attribute {
		if nw.FailCode != archive.XcodeFailZero {
			remarks["xcode_fail"] = nw.FailCode
		}
		if nw.Title != old.Title && nw.Title != "" {
			remarks["title"] = nw.Title
		}
		if nw.Desc != old.Desc && nw.Desc != "" {
			remarks["desc"] = nw.Desc
		}
	} else {
		// no change
		return
	}
	if len(remarks) != 0 {
		bs, err = json.Marshal(remarks)
	}
	s.arc.AddVideoTrack(context.TODO(), nw.Aid, nw.Filename, nw.Status, nw.XcodeState, string(bs), nw.MTime, nw.MTime)
	return
}
