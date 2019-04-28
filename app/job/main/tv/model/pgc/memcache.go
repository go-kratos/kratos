package pgc

// SimpleEP is the structure of ep in mc
type SimpleEP struct {
	ID        int64 `json:"id"`
	EPID      int   `json:"epid"`
	SeasonID  int   `json:"season_id"`
	State     int   `json:"state"`
	Valid     int   `json:"valid"`
	IsDeleted int   `json:"is_deleted"`
	NoMark    int   `json:"no_mark"`
}

// SimpleSeason is the structure of season in mc
type SimpleSeason struct {
	ID        int64 `json:"id"`
	IsDeleted int8  `json:"is_deleted"`
	Valid     int   `json:"valid"`
	Check     int8  `json:"check"`
}
