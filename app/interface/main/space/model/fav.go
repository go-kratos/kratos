package model

import "go-common/app/service/main/favorite/model"

// FavNav fav nav struct.
type FavNav struct {
	Archive  []*model.VideoFolder `json:"archive"`
	Playlist int                  `json:"playlist"`
	Topic    int                  `json:"topic"`
	Article  int                  `json:"article"`
	Album    int                  `json:"album"`
	Movie    int                  `json:"movie"`
}

// FavArcArg .
type FavArcArg struct {
	Vmid    int64  `form:"vmid" validate:"min=1"`
	Fid     int64  `form:"fid" validate:"min=-1"`
	Tid     int64  `form:"tid"`
	Keyword string `form:"keyword"`
	Order   string `form:"order"`
	Pn      int    `form:"pn" default:"1" validate:"min=1"`
	Ps      int    `form:"ps" default:"30" validate:"min=1"`
}
