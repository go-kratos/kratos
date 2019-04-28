package model

import (
	"go-common/library/time"
)

// APKInfo .
type APKInfo struct {
	ID          int64     `json:"id"`
	CDNAddr     string    `json:"cdn_addr"`
	CreatedAt   time.Time `json:"created_at"`
	FileMd5     string    `json:"file_md5"`
	InetAddr    string    `json:"inet_addr"`
	IsDeleted   bool      `json:"is_deleted"`
	IsGray      bool      `json:"is_gray"`
	LocalPath   string    `json:"local_path"`
	MappingAddr string    `json:"mapping_addr"`
	SignMd5     string    `json:"sign_md5"`
	Size        int       `json:"size"`
	UpdatedAt   time.Time `json:"updated_at"`
	VersionCode int       `json:"version_code"`
	VersionID   string    `json:"version_id"`
	VersionName string    `json:"version_name"`
}

// MangoRecom is mango recom table structure
type MangoRecom struct {
	ID        int64  `json:"id" gorm:"column:id"`
	RID       int64  `json:"rid" gorm:"column:rid"`
	Rtype     int    `json:"rtype"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Category  int    `json:"category"`
	Playcount int64  `json:"playcount"`
	JID       int64  `json:"jid" gorm:"column:jid"`
	Content   string `json:"content"`
	Staff     string `json:"staff"`
	Rorder    int    `json:"rorder"`
}

// MangoListResp is the mango list response structure
type MangoListResp struct {
	List    []*MangoRecom `json:"list"`
	Pubtime string        `json:"pubtime"`
	Message string        `json:"message"` // 文案提示: rid, p213,u367 ...
}

// MangoAdd is the response of mango add function
type MangoAdd struct {
	Succ     []int64 `json:"succ"`
	Invalids []int64 `json:"invalids"`
}

// TableName def.
func (*MangoRecom) TableName() string {
	return "mango_recom"
}

// ToMango def.
func (sn *TVEpSeason) ToMango() *MangoRecom {
	return &MangoRecom{
		RID:      sn.ID,
		Rtype:    1,
		Title:    sn.Title,
		Cover:    sn.Cover,
		Category: int(sn.Category),
		Content:  sn.Desc,
		Staff:    sn.Staff,
	}
}

// ToMango def.
func (arc *SimpleArc) ToMango(cat int) *MangoRecom {
	return &MangoRecom{
		RID:      arc.AID,
		Rtype:    2,
		Title:    arc.Title,
		Cover:    arc.Cover,
		Category: cat,
		Content:  arc.Content,
	}
}

// ReqMangoEdit is the request for mango edit
type ReqMangoEdit struct {
	ID        int64  `form:"id" validate:"required"`
	Title     string `form:"title" validate:"required"`
	Cover     string `form:"cover" validate:"required"`
	Playcount int64  `form:"playcount"`
	JID       int64  `form:"jid"`
	Content   string `form:"content" validate:"required"`
	Staff     string `form:"staff"`
}

// MRecomMC is mango recom struct in MC
type MRecomMC struct {
	RIDs    []int64
	Pubtime time.Time
}

// ReqUnshelve is request for unshelve
type ReqUnshelve struct {
	IDs  []int64 `form:"ids,split" validate:"required,min=1,dive,gt=0"`
	Type int     `form:"type" validate:"required,min=1,max=4"`
}

// RespUnshelve is response for unshelve
type RespUnshelve struct {
	SuccIDs []int64 `json:"succ_ids"`
	FailIDs []int64 `json:"fail_ids"`
}
