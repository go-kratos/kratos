package model

import "strconv"

// ArchiveCheckParams search params.
type ArchiveCheckParams struct {
	Bsp          *BasicSearchParams
	Aids         []int64 `form:"aids,split" params:"aids"`
	TypeIds      []int64 `form:"typeids,split" params:"typeids"`
	States       []int64 `form:"states,split" params:"states"`
	Attrs        []int64 `form:"attrs,split" params:"attrs"`
	DurationFrom int64   `form:"duration_from" params:"duration_from"`
	DurationTo   int64   `form:"duration_to" params:"duration_to"`
	Mids         []int64 `form:"mids,split" params:"mids"`
	MidFrom      int64   `form:"mid_from" params:"mid_from"`
	MidTo        int64   `form:"mid_to" params:"mid_to"`
	AllKW        int     `form:"all_kw" params:"all_kw" default:"0"`
	TimeFrom     string  `form:"time_from" params:"time_from"`
	TimeTo       string  `form:"time_to" params:"time_to"`
	Time         string  `form:"time" params:"time"`
	FromIP       string  `form:"from_ip" params:"from_ip"`
}

// VideoParams search video
type VideoParams struct {
	Bsp        *BasicSearchParams
	VIDs       []int64  `form:"vids,split" params:"vids"`
	AIDs       []int64  `form:"aids,split" params:"aids"`
	CIDs       []int64  `form:"cids,split" params:"cids"`
	TIDs       []int64  `form:"tids,split" params:"tids"`
	FileNames  []string `form:"filename,split" params:"filename"`
	TagID      int64    `form:"tag_id" params:"tag_id"`
	Status     []int64  `form:"status,split" params:"status"`
	XCodeState []int64  `form:"xcode_state,split" params:"xcode_state"`
	UserType   int      `form:"user_type" params:"user_type"`
	// archive
	RelationStates []int64 `form:"relation_state,split" params:"relation_state"`
	ArcMids        []int64 `form:"arc_mids,split" params:"arc_mids"`
	DurationFrom   int     `form:"duration_from" params:"duration_from"`
	DurationTo     int     `form:"duration_to" params:"duration_to"`
	// other
	OrderType int `form:"order_type" params:"order_type"`
}

// TaskQa .
type TaskQa struct {
	Bsp           *BasicSearchParams
	Ids           []int64  `form:"ids,split" params:"ids"`
	TaskIds       []string `form:"task_ids,split" params:"task_ids"`
	Uids          []string `form:"uids,split" params:"uids"`
	ArcTagIds     []string `form:"arc_tagids,split" params:"arc_tagids"`
	AuditTagIds   []int64  `form:"audit_tagids,split" params:"audit_tagids"`
	UpGroups      []string `form:"up_groups,split" params:"up_groups"`
	ArcTitles     []string `form:"arc_titles,split" params:"arc_titles"`
	ArcTypeIds    []string `form:"arc_typeids,split" params:"arc_typeids"`
	States        []string `form:"states,split" params:"states"`
	AuditStatuses []string `form:"audit_statuses,split" params:"audit_statuses"`
	FansFrom      string   `form:"fans_from" params:"fans_from"`
	FansTo        string   `form:"fans_to" params:"fans_to"`
	CtimeFrom     string   `form:"ctime_from" params:"ctime_from"`
	CtimeTo       string   `form:"ctime_to" params:"ctime_to"`
	FtimeFrom     string   `form:"ftime_from" params:"ftime_from"`
	FtimeTo       string   `form:"ftime_to" params:"ftime_to"`
}

// ArchiveCommerce .
type ArchiveCommerce struct {
	Bsp        *BasicSearchParams
	Ids        []string `form:"ids,split" params:"ids"`
	PTypeIds   []string `form:"ptypeids,split" params:"ptypeids"`
	TypeIds    []string `form:"typeids,split" params:"typeids"`
	Mids       []string `form:"mids,split" params:"mids"`
	States     []string `form:"states,split" params:"states"`
	Copyrights []string `form:"copyrights,split" params:"copyrights"`
	OrderIds   []string `form:"order_ids,split" params:"order_ids"`
	// 逻辑判断
	Action     string `form:"action" params:"action"`                        // 获取一级分区列表、等其他定制查询
	IsOrder    int    `form:"is_order" params:"is_order" default:"-1"`       //是否商单
	IsOriginal int    `form:"is_original" params:"is_original" default:"-1"` //是否原创
}

// TaskQaFansParams .
type TaskQaFansParams struct {
	ID   int64 `json:"id"`
	Fans int64 `json:"fans"`
}

// IndexName .
func (m *TaskQaFansParams) IndexName() string {
	return "task_qa"
}

// IndexType .
func (m *TaskQaFansParams) IndexType() string {
	return "base"
}

// IndexID .
func (m *TaskQaFansParams) IndexID() string {
	return strconv.FormatInt(m.ID, 10)
}
