package ugc

import (
	v1 "go-common/app/service/main/archive/api"
	"go-common/library/time"
)

// Archive archive def. corresponding to our table structure
type Archive struct {
	ID        int
	AID       int64
	MID       int64
	TypeID    int32
	Videos    int64
	Title     string
	Cover     string
	Content   string
	Duration  int64
	Copyright int32
	Pubtime   time.Time
	Ctime     time.Time
	Mtime     time.Time
	State     int32
	Manual    int
	Valid     int
	Submit    int
	Retry     int
	Result    int
	Deleted   int
}

// FromArcReply def
func (a *Archive) FromArcReply(arc *v1.Arc) {
	a.AID = arc.Aid
	a.MID = arc.Author.Mid
	a.Videos = arc.Videos
	a.TypeID = arc.TypeID
	a.Title = arc.Title
	a.Cover = arc.Pic
	a.Content = arc.Desc
	a.Duration = arc.Duration
	a.Copyright = arc.Copyright
	a.Pubtime = arc.PubDate
	a.State = arc.State
}

// ArcAllow is the struct used to check whether the arc is allowed to enter TV database
type ArcAllow struct {
	Aid       int64
	State     int32
	Ugcpay    int32
	Typeid    int32
	Copyright int32
}

// FromArcReply takes info from grpc result
func (a *ArcAllow) FromArcReply(reply *v1.Arc) {
	a.Aid = reply.Aid
	a.State = reply.State
	a.Ugcpay = reply.Rights.UGCPay
	a.Typeid = reply.TypeID
	a.Copyright = reply.Copyright
}

// FromArcmdl takes info from gorpc result
func (a *ArcAllow) FromArcmdl(mdl *v1.Arc) {
	a.Aid = mdl.Aid
	a.State = mdl.State
	a.Ugcpay = mdl.Rights.UGCPay
	a.Typeid = mdl.TypeID
	a.Copyright = mdl.Copyright
}

// FromDatabus takes info from databus result ( archive-notify T )
func (a *ArcAllow) FromDatabus(db *ArchDatabus) {
	a.Aid = db.Aid
	a.State = db.State
	a.Typeid = db.TypeID
	a.Copyright = db.Copyright
}

// FromArcFull takes info from arcFull structure ( db )
func (a *ArcAllow) FromArcFull(full *ArcFull) {
	a.Aid = full.AID
	a.State = full.State
	a.Copyright = full.Copyright
	a.Typeid = full.TypeID
}

// FromArchive takes info from DB
func (a *ArcAllow) FromArchive(arc *Archive) {
	a.Aid = arc.AID
	a.State = arc.State
	a.Copyright = arc.Copyright
	a.Typeid = arc.TypeID
}

// CanPlay distinguishes whether an archive can play or not
func (a *ArcAllow) CanPlay() bool {
	return a.State >= 0 || a.State == -6
}

// IsOrigin distinguishes whether an archive is original or not
func (a *ArcAllow) IsOrigin() bool {
	return a.Copyright == 1
}

// ArcMedia is the archive media struct in MC
type ArcMedia struct {
	Title   string
	AID     int64
	Cover   string
	TypeID  int32
	Pubtime time.Time
	Videos  int64
	Deleted int
}

// DelVideos is used to delete videos of an archive
type DelVideos struct {
	AID  int64
	CIDs []int64
}

// ToSimple def.
func (a *Archive) ToSimple() *SimpleArc {
	return &SimpleArc{
		AID:      a.AID,
		MID:      a.MID,
		TypeID:   a.TypeID,
		Videos:   a.Videos,
		Title:    a.Title,
		Cover:    a.Cover,
		Content:  a.Content,
		Duration: a.Duration,
		Pubtime:  a.Pubtime.Time().Format("2006-01-02"),
	}
}
