package model

// ArchiveInfo archive info.
type ArchiveInfo struct {
	UpCnt      int    `json:"up_cnt"`
	ArchiveCnt int    `json:"archive_cnt"`
	UploadDate string `json:"log_date"`
}

// ArchiveID archive id
type ArchiveID struct {
	ID int64 `json:"id"`
}

// ArchiveStat for archive stat.
type ArchiveStat struct {
	ID    int64
	AvID  int64
	State int
	Play  int64 `json:"play"`
	Dm    int64 `json:"dm"`
	Reply int64 `json:"reply"`
	Like  int64 `json:"like"`
	Share int64 `json:"share"`
}
