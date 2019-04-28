package model

const (
	_OldIdx     = 0
	_OldZone    = 1
	_PgcIdx     = 2
	_UgcIdx     = 3
	_homepageID = 0
)

// Module def.
type Module struct {
	ID          int        `json:"id"`
	PageID      int        `json:"page_id"`
	Type        int        `json:"type"`
	Title       string     `json:"title"`
	Icon        string     `json:"icon"`
	Source      int        `json:"source"`
	Flexible    int        `json:"flexible"`
	Capacity    int        `json:"capacity"`
	More        int        `json:"more"`
	MoreType    int        `json:"more_type"`
	MoreNewPage int        `json:"more_new_page"`
	MorePage    int        `json:"more_page"`
	Order       int        `json:"order"`
	Data        []*ModCard `json:"data"`
	SrcType     int        `json:"src_type"`
}

// JumpNewIdx tells whether this modules jumps to the new idx page
func (m *Module) JumpNewIdx() bool {
	return m.MoreType == _PgcIdx || m.MoreType == _UgcIdx
}

// OnHomepage tells whether the module is on the homepage
func (m *Module) OnHomepage() bool {
	return m.PageID == _homepageID
}

// MoreTreat treats the morepage and moretype related, used for zone/modpages, not homepage
func (m *Module) MoreTreat() {
	if m.MorePage == 0 {
		return
	}
	if m.MoreType == _OldIdx || m.MoreType == _OldZone { // if more jump setting is old zone/idx (<=1.13), set more_new_page = more_page
		m.MoreNewPage = m.MorePage
	}
	if m.MoreType == _PgcIdx || m.MoreType == _UgcIdx { // if more jump setting is new ugc/pgc idx, set more_new_page = Idx category, more_page = page_id
		m.MoreNewPage = m.MorePage
		m.MorePage = m.PageID
	}
}

// IsUGC returns whether the module is filled by ugc or not
func (m Module) IsUGC() bool {
	return m.SrcType == _TypeUGC
}

// ModCard structure, based on normal Card, 4 more fields for Follow Module
type ModCard struct {
	Card
	LastEPIndex   string `json:"last_ep_index"`
	NewestEPIndex string `json:"newest_ep_index"`
	TotalCount    string `json:"total_count"`
	IsFinish      string `json:"is_finish"`
}

// ReqModData is the request body to modData function
type ReqModData struct {
	Mod      *Module
	PGCListM map[int][]*Card
	UGCListM map[int][]*Card
}

// ReqPageFollow is the request body to PageFollow function
type ReqPageFollow struct {
	AccessKey string `form:"access_key"`
	PageID    int    `form:"page_id" validate:"min=0"`
	Build     int    `form:"build"`
}

// ReqHomeFollow is the request body to HomeFollow function
type ReqHomeFollow struct {
	AccessKey string `form:"access_key"`
	Build     int    `form:"build"`
}
