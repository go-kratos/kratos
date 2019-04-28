package model

// UpArtStat struct.
type UpArtStat struct {
	View      int64 `json:"view"`
	Reply     int64 `json:"reply"`
	Like      int64 `json:"like"`
	Coin      int64 `json:"coin"`
	Fav       int64 `json:"fav"`
	PreView   int64 `json:"-"`
	PreReply  int64 `json:"-"`
	PreLike   int64 `json:"-"`
	PreCoin   int64 `json:"-"`
	PreFav    int64 `json:"-"`
	IncrView  int64 `json:"incr_view"`
	IncrReply int64 `json:"incr_reply"`
	IncrLike  int64 `json:"incr_like"`
	IncrCoin  int64 `json:"incr_coin"`
	IncrFav   int64 `json:"incr_fav"`
}
