package search

import (
	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/music"
	"go-common/app/interface/main/creative/model/reply"
)

const (
	//All for all of reply type.
	All = 0
	//Archive for reply type.
	Archive = 1
	//Article for reply type.
	Article = 12
	//Audio for reply type.
	Audio = 14
	//SmallVideo for reply type.
	SmallVideo = 5
)

// Result search list.
type Result struct {
	Class       *ClassCount                     `json:"class"`
	Applies     *ApplyStateCount                `json:"apply_count"`
	Type        map[int16]*TypeCount            `json:"-"`
	ArrType     []*TypeCount                    `json:"type"`
	OldArchives []*archive.OldArchiveVideoAudit `json:"archives"`
	Archives    []*archive.ArcVideoAudit        `json:"arc_audits"`
	Page        struct {
		Pn    int `json:"pn"`
		Ps    int `json:"ps"`
		Count int `json:"count"`
	} `json:"page"`
	Aids []int64 `json:"-"` // from search, call archiveRPC
	Tip  string  `json:"tip"`
}

// StaffApplyResult search list.
type StaffApplyResult struct {
	StateCount *ApplyStateCount     `json:"state_count"`
	Type       map[int16]*TypeCount `json:"-"`
	ArrType    []*TypeCount         `json:"type"`
	Applies    []*StaffApply        `json:"applies"`
	Page       struct {
		Pn    int `json:"pn"`
		Ps    int `json:"ps"`
		Count int `json:"count"`
	} `json:"page"`
	Aids []int64 `json:"-"` // from search, call archiveRPC
}

// StaffApply str
type StaffApply struct {
	ID         int64                  `json:"id"`
	Type       int8                   `json:"type"`
	Mid        int64                  `json:"mid"`
	Uname      string                 `json:"uname"`
	State      int8                   `json:"state"`
	ApplyTitle string                 `json:"apply_title"`
	ApplyState int8                   `json:"apply_state"`
	Archive    *archive.ArcVideoAudit `json:"arc_audits"`
}

// ApplyStateCount pub count.
type ApplyStateCount struct {
	Neglected int `json:"neglected"`
	Pending   int `json:"pending"`
	Processed int `json:"processed"`
}

// ClassCount pub count.
type ClassCount struct {
	Pubed    int `json:"pubed"`
	NotPubed int `json:"not_pubed"`
	Pubing   int `json:"is_pubing"`
}

// TypeCount archive count for a type.
type TypeCount struct {
	Tid   int16  `json:"tid"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

// Reply  str
type Reply struct {
	Message    string       `json:"message"`
	ID         int64        `json:"id"`
	Floor      int64        `json:"floor"`
	Count      int          `json:"count"`
	Root       int64        `json:"root"`
	Oid        int64        `json:"oid"`
	CTime      string       `json:"ctime"`
	MTime      string       `json:"mtime"`
	State      int          `json:"state"`
	Parent     int64        `json:"parent"`
	Mid        int64        `json:"mid"`
	Like       int          `json:"like"`
	Replier    string       `json:"replier"`
	Uface      string       `json:"uface"`
	Cover      string       `json:"cover"`
	Title      string       `json:"title"`
	Relation   int          `json:"relation"`
	IsElec     int          `json:"is_elec"`
	Type       int          `json:"type"`
	RootInfo   *reply.Reply `json:"root_info"`
	ParentInfo *reply.Reply `json:"parent_info"`
}

//Replies for reply list.
type Replies struct {
	SeID       string          `json:"seid"`
	Order      string          `json:"order"`
	Keyword    string          `json:"keyword"`
	Total      int             `json:"total"`
	PageCount  int             `json:"pagecount"`
	Repliers   []int64         `json:"repliers"`
	DeriveOids []int64         `json:"-"`
	DeriveIds  []int64         `json:"-"`
	Oids       []int64         `json:"-"`
	TyOids     map[int][]int64 `json:"-"`
	Result     []*Reply        `json:"result"`
}

//SimpleResult for archives simple result.
type SimpleResult struct {
	ArchivesVideos []*SimpleArcVideos `json:"simple_arc_videos"`
	Class          *ClassCount        `json:"class"`
	Page           struct {
		Pn    int `json:"pn"`
		Ps    int `json:"ps"`
		Count int `json:"count"`
	} `json:"page"`
}

//SimpleArcVideos for search archive & vidoes.
type SimpleArcVideos struct {
	Archive *archive.SimpleArchive `json:"archive"`
	Videos  []*archive.SimpleVideo `json:"videos"`
}

// ArcParam for es search param.
type ArcParam struct {
	MID     int64
	AID     int64
	TypeID  int64
	Pn      int
	Ps      int
	State   string
	Keyword string
	Order   string
}

// Arc for search archive.
type Arc struct {
	ID       int64  `json:"id"`
	TypeID   int64  `json:"typeid"`
	PID      int64  `json:"pid"`
	State    int64  `json:"state"`
	Duration int64  `json:"duration"`
	Title    string `json:"title"`
	Cover    string `json:"cover"`
	Desc     string `json:"description"`
	PubDate  string `json:"pubdate"`
}

// Pager for es page.
type Pager struct {
	Num   int `json:"num"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

// DCount for state count.
type DCount struct {
	Count int `json:"doc_count"`
}

// PList for state count list.
type PList struct {
	IsPubing DCount `json:"is_pubing"`
	NotPubed DCount `json:"not_pubed"`
	Pubed    DCount `json:"pubed"`
	Pending  DCount `json:"pending"`
}

// TList for type count list.
type TList struct {
	DCount
	Key string `json:"key"`
}

// ArcResult archive list from search.
type ArcResult struct {
	Page   *Pager `json:"page"`
	Result struct {
		Vlist []*Arc   `json:"vlist"`
		PList *PList   `json:"plist"`
		TList []*TList `json:"tlist"`
	} `json:"result"`
}

// ReliesES str
type ReliesES struct {
	Order  string     `json:"order"`
	Sort   string     `json:"sort"`
	Page   *Pager     `json:"page"`
	Result []*ReplyES `json:"result"`
}

// ReplyES  str
type ReplyES struct {
	Count   int    `json:"count"`
	CTime   string `json:"ctime"`
	Floor   int64  `json:"floor"`
	Hate    int64  `json:"hate"`
	ID      int64  `json:"id"`
	Like    int    `json:"like"`
	Message string `json:"message"`
	Mid     int64  `json:"mid"`
	MTime   string `json:"mtime"`
	OMid    int64  `json:"o_mid"`
	Oid     int64  `json:"oid"`
	Parent  int64  `json:"parent"`
	Rcount  int64  `json:"rcount"`
	Root    int64  `json:"root"`
	State   int    `json:"state"`
	Type    int    `json:"type"`
}

//ReplyParam str
type ReplyParam struct {
	Ak          string
	Ck          string
	OMID        int64
	OID         int64
	Pn          int
	Ps          int
	IsReport    int8
	Type        int8
	ResMdlPlat  int8
	FilterCtime string
	Kw          string
	Order       string
	IP          string
}

// Bgm str
type Bgm struct {
	SID int64 `json:"sid"`
}

// BgmResult str
type BgmResult struct {
	Page   *Pager `json:"page"`
	Result []*Bgm `json:"result"`
}

// MaterialRel str
type MaterialRel struct {
	AID int64 `json:"aid"`
}

// BgmExtResult str
type BgmExtResult struct {
	Page   *Pager         `json:"page"`
	Result []*MaterialRel `json:"result"`
}

// BgmSearchRes str
type BgmSearchRes struct {
	Pager *Pager         `json:"pager"`
	Bgms  []*music.Music `json:"bgm"`
}

// ApplyResult apply list from search.
type ApplyResult struct {
	Page   *Pager `json:"page"`
	Result struct {
		Vlist      []*Arc      `json:"vlist"`
		ApplyPList *ApplyPList `json:"plist"`
		TList      []*TList    `json:"tlist"`
	} `json:"result"`
}

// ApplyPList for apply state count list.
type ApplyPList struct {
	Neglected DCount `json:"neglected"`
	Pending   DCount `json:"pending"`
	Processed DCount `json:"processed"`
}
