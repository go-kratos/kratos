package archive

import "go-common/library/time"

//ArcHistory 稿件的用户编辑历史
type ArcHistory struct {
	ID      int64     `json:"id"`
	AID     int64     `json:"aid"`
	Title   string    `json:"title"`
	Tag     string    `json:"tag"`
	Content string    `json:"content"`
	Cover   string    `json:"cover"`
	MID     int64     `json:"mid"`
	CTime   time.Time `json:"ctime"`
}

//VideoHistory 视频的用户编辑历史
type VideoHistory struct {
	ID          int64     `json:"id"`
	CID         int64     `json:"cid"`
	EpTitle     string    `json:"eptitle"`
	Description string    `json:"description"`
	Filename    string    `json:"filename"`
	SRCType     string    `json:"src_type"`
	CTime       time.Time `json:"ctime"`
}

//EditHistory 一次完整的用户编辑历史
type EditHistory struct {
	ArcHistory *ArcHistory     `json:"arc_history"`
	VHistory   []*VideoHistory `json:"v_history"`
}

func (ah *ArcHistory) diff(one *ArcHistory) (res *ArcHistory, diff bool) {
	if one == nil {
		res = ah
		diff = true
		return
	}

	res = &ArcHistory{
		ID:    ah.ID,
		AID:   ah.AID,
		CTime: ah.CTime,
	}
	if ah.Title != one.Title {
		res.Title = ah.Title
		diff = true
	}
	if ah.Tag != one.Tag {
		res.Tag = ah.Tag
		diff = true
	}
	if ah.Content != one.Content {
		res.Content = ah.Content
		diff = true
	}
	if ah.Cover != one.Cover {
		res.Cover = ah.Cover
		diff = true
	}
	if ah.MID != one.MID {
		res.MID = ah.MID
		diff = true
	}
	return
}

//Diff only show diff between next archive edit history
func (eh *EditHistory) Diff(one *EditHistory) (res *EditHistory, diff bool) {
	if one == nil {
		res = eh
		diff = true
		return
	}

	var ah *ArcHistory
	vh := []*VideoHistory{}
	oldfs := map[string]int{}
	ah, diff = eh.ArcHistory.diff(one.ArcHistory)

	//show those  whose filenames not exist in one
	for _, v := range one.VHistory {
		oldfs[v.Filename] = 1
	}
	for _, v := range eh.VHistory {
		if oldfs[v.Filename] != 1 {
			vh = append(vh, v)
			diff = true
		}
	}

	res = &EditHistory{
		ArcHistory: ah,
		VHistory:   vh,
	}
	return
}
