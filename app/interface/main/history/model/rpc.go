package model

import hismdl "go-common/app/service/main/history/model"

// History video hisotry info.
type History struct {
	Mid      int64  `json:"mid,omitempty"`
	Aid      int64  `json:"aid"`
	Sid      int64  `json:"sid,omitempty"`
	Epid     int64  `json:"epid,omitempty"`
	TP       int8   `json:"tp,omitempty"`
	Business string `json:"business"`
	STP      int8   `json:"stp,omitempty"` // sub_type
	Cid      int64  `json:"cid,omitempty"`
	DT       int8   `json:"dt,omitempty"`
	Pro      int64  `json:"pro,omitempty"`
	Unix     int64  `json:"view_at"`
}

// Histories history sorted.
type Histories []*History

func (h Histories) Len() int { return len(h) }
func (h Histories) Less(i, j int) bool {
	if h[i].Unix == h[j].Unix {
		return h[i].Aid < h[j].Aid
	}
	return h[i].Unix > h[j].Unix
}
func (h Histories) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// FillBusiness add history
func (h *History) FillBusiness() {
	if h == nil {
		return
	}
	h.Business = businessIDs[h.TP]
}

// ConvertType convert old type
func (h *History) ConvertType() {
	if h == nil {
		return
	}
	switch h.TP {
	case TypeBangumi:
		h.TP = TypePGC
		h.STP = SubTypeBangumi
	case TypeMovie:
		h.TP = TypePGC
		h.STP = SubTypeFilm
	case TypePGC:
		if h.Epid == 0 || h.Sid == 0 {
			h.TP = TypeUGC
		}
	}
}

// ConvertServiceType .
func (h History) ConvertServiceType() (r *hismdl.History) {
	switch h.TP {
	case TypeOffline:
		h.TP = TypeUGC
		h.STP = SubTypeOffline
	case TypeUnknown:
		h.TP = TypeUGC
	case TypeBangumi:
		h.TP = TypePGC
		h.STP = SubTypeBangumi
	case TypeMovie:
		h.TP = TypePGC
		h.STP = SubTypeFilm
	}
	if h.TP == TypePGC && (h.Epid == 0 || h.Sid == 0) {
		h.TP = TypeUGC
	}
	h.FillBusiness()
	r = &hismdl.History{
		Mid:        h.Mid,
		BusinessID: int64(h.TP),
		Business:   h.Business,
		Kid:        h.Aid,
		Aid:        h.Aid,
		Sid:        h.Sid,
		Epid:       h.Epid,
		Cid:        h.Cid,
		SubType:    int32(h.STP),
		Device:     int32(h.DT),
		Progress:   int32(h.Pro),
		ViewAt:     h.Unix,
	}
	if h.TP == TypePGC {
		r.Kid = r.Sid
	}
	return
}

// ArgPro arg.
type ArgPro struct {
	Mid    int64
	RealIP string
	Aids   []int64
}

// ArgPos arg.
type ArgPos struct {
	Mid      int64
	Aid      int64
	Business string
	TP       int8
	RealIP   string
}

// ArgDelete arg.
type ArgDelete struct {
	Mid       int64
	RealIP    string
	Resources []*Resource
}

// ArgHistory arg.
type ArgHistory struct {
	Mid      int64
	Realtime int64
	RealIP   string
	History  *History
}

// ArgHistories arg.
type ArgHistories struct {
	Mid      int64
	TP       int8
	Business string
	Pn       int
	Ps       int
	RealIP   string
}

// ArgCursor arg.
type ArgCursor struct {
	Mid int64
	Max int64
	TP  int8
	// history business
	Business string
	ViewAt   int64
	// filter business, blank means all business
	Businesses []string
	Ps         int
	RealIP     string
}

// Resource video hisotry info .
type Resource struct {
	Mid      int64  `json:"mid,omitempty"`
	Oid      int64  `json:"oid"`
	Sid      int64  `json:"sid,omitempty"`
	Epid     int64  `json:"epid,omitempty"`
	TP       int8   `json:"tp,omitempty"`
	STP      int8   `json:"stp,omitempty"` // sub_type
	Cid      int64  `json:"cid,omitempty"`
	Business string `json:"business"`
	DT       int8   `json:"dt,omitempty"`
	Pro      int64  `json:"pro,omitempty"`
	Unix     int64  `json:"view_at"`
}

// ArgClear .
type ArgClear struct {
	Mid        int64
	RealIP     string
	Businesses []string
}
