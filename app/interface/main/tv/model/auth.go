package model

const (
	_epPass          = 3
	_epRejected      = 4
	_noMarkWhiteList = 1
)

// EpAuth is the structure of ep in mc
type EpAuth struct {
	ID        int64 `json:"id"`
	EPID      int64 `json:"epid"`
	SeasonID  int64 `json:"season_id"`
	State     int   `json:"state"`
	Valid     int   `json:"valid"`
	IsDeleted int   `json:"is_deleted"`
	NoMark    int   `json:"no_mark"`
}

// SnAuth is the structure of season in mc
type SnAuth struct {
	ID        int64 `json:"id"`
	IsDeleted int8  `json:"is_deleted"`
	Valid     int   `json:"valid"`
	Check     int8  `json:"check"`
}

// NotDeleted def.
func (s SnAuth) NotDeleted() bool {
	return s.IsDeleted == 0
}

// NotDeleted def.
func (s EpAuth) NotDeleted() bool {
	return s.IsDeleted == 0
}

// CanPlay returns whether the season is able to play
func (s EpAuth) CanPlay() bool {
	return s.IsDeleted == 0 && s.Valid == 1 && s.State == 3
}

// Auditing checks whether the ep is still auditing
// func (s EpAuth) Auditing() bool {
// 	return s.State != _epPass && s.State != _epRejected && s.IsDeleted == _noDel
// }

// Whitelist checks whether the ep is in the whitelist of no mark eps
func (s EpAuth) Whitelist() bool {
	return s.NoMark == _noMarkWhiteList
}

// CanPlay returns whether the season is able to play
func (s SnAuth) CanPlay() bool {
	return s.IsDeleted == 0 && s.Valid == 1 && s.Check == 1
}

// ArcType def.
type ArcType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
