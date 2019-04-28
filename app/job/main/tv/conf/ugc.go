package conf

import (
	xtime "go-common/library/time"
)

// UgcSync defines the params for ugc sync
type UgcSync struct {
	Frequency *UgcFre
	Batch     *Batch
	Cfg       *UgcCfg
}

// UgcCfg is for various ugc cfg
type UgcCfg struct {
	Copyright    string
	ReportCidURL string // the url of VideoCloud for reporting Cid
	BFSPrefix    string
	CriticalCid  int64 // critical cid 12780000, under it no need to ask for transcoding
	ThreadLimit  int64 // thread limit
}

// Batch is for the number of data to pick each time
type Batch struct {
	ManualNum   int // manually added archives
	ImportNum   int // the number of uppers to import all his video
	ArcPS       int // the page size to pick the upper's archives
	SyncPS      int // the page size of sync message ( nb of videos )
	ReportCidPS int // the page size to update cid's mark status
	ProducerPS  int // producer page size
	ReshelfPS   int // reshelf arc page size
}

// UgcFre defines the ugc sync frequencies
type UgcFre struct {
	ErrorWait      int            // postpone the operation due to error
	ManualFre      xtime.Duration // re-check the manual import need frequency
	ImportFre      xtime.Duration // import upper's video frequency
	TypesCron      string         // import the ugc types cron
	SyncFre        xtime.Duration
	UpperRefresh   xtime.Duration // upper refresh duration
	ReportCid      xtime.Duration // 1 minute to check report cid
	UpInitFre      xtime.Duration // pause between each page of upper's archive
	UpperPause     xtime.Duration // pause between each import upper
	ProducerFre    xtime.Duration // producer pause time
	FullRefreshFre xtime.Duration // video refresh frequency
	FullRefArcFre  xtime.Duration // pause between each archive
}
