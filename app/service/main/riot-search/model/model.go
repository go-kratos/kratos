package model

// Document id and content
type Document struct {
	ID      uint64 `json:"id"`
	Content string `json:"content"`
}

// RiotSearchReq search request params
type RiotSearchReq struct {
	IDs     []uint64 `form:"ids,split"`
	Keyword string   `form:"keyword" validate:"required"`
	Pn      int      `form:"pn" validate:"min=1"`
	Ps      int      `form:"ps" validate:"min=0"`
}

// IDsResp resp of ids
type IDsResp struct {
	IDs    []uint64 `json:"ids"`
	Tokens []string `json:"tokens"`
	Page   *Page    `json:"page"`
}

// DocumentsResp resp of documents
type DocumentsResp struct {
	Documents []Document `json:"ducuments"`
	Tokens    []string   `json:"tokens"`
	Page      *Page      `json:"page"`
}

// Page Pager
type Page struct {
	PageNum  int `json:"pn"`
	PageSize int `json:"ps"`
	Total    int `json:"total"`
}

// **********************
// * Model for archives *
// **********************

// ArchiveMessage databus message
type ArchiveMessage struct {
	Action string       `json:"action"`
	Table  string       `json:"table"`
	New    *ArchiveMeta `json:"new"`
	Old    *ArchiveMeta `json:"old"`
}

// ArchiveMeta Archive Metadata
type ArchiveMeta struct {
	AID   uint64 `json:"aid"`
	Title string `json:"title"`
	State int    `json:"state"`
}

// States archive states
type States struct {
	LegalStates map[int]bool
}

// PubStates publish states
var PubStates = &States{
	LegalStates: map[int]bool{
		-40:   true,
		0:     true,
		10000: true,
		1:     true,
		1001:  true,
		15000: true,
		20000: true,
		30000: true,
	},
}

// Legal return leagal
func (l *States) Legal(state int) bool {
	if _, ok := l.LegalStates[state]; ok {
		return true
	}
	return false
}
