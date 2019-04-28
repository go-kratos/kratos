package ugc

import "go-common/library/time"

// ArcCMS represents the archive data structure in MC
type ArcCMS struct {
	// Media Info
	Title   string
	AID     int64
	Content string
	Cover   string
	TypeID  int32
	Pubtime time.Time
	Videos  int
	// Auth Info
	Valid   int
	Deleted int
	Result  int
}

// ArcFull is the plus version of ArcCMS
type ArcFull struct {
	ArcCMS
	Copyright int32
	State     int32
	MID       int64
	Duration  int64
}

// VideoCMS represents the video data structure in MC
type VideoCMS struct {
	// Media Info
	CID        int
	Title      string
	AID        int
	IndexOrder int
	// Auth Info
	Valid   int
	Deleted int
	Result  int
}

// ToSimple transforms an arcFull to SimpleArc
func (arc *ArcFull) ToSimple() *SimpleArc {
	return &SimpleArc{
		AID:      arc.AID,
		MID:      arc.MID,
		TypeID:   arc.TypeID,
		Videos:   int64(arc.Videos),
		Title:    arc.Title,
		Cover:    arc.Cover,
		Content:  arc.Content,
		Duration: arc.Duration,
		Pubtime:  arc.Pubtime.Time().Format("2006-01-02"),
	}
}
