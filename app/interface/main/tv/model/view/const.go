package view

import (
	"fmt"

	arcwar "go-common/app/service/main/archive/api"
)

// Goto
const (
	GotoAv       = "av"
	GotoWeb      = "web"
	GotoBangumi  = "bangumi"
	GotoLive     = "live"
	GotoGame     = "game"
	GotoArticle  = "article"
	GotoSpecial  = "special"
	GotoAudio    = "audio"
	GotoSong     = "song"
	GotoAudioTag = "audio_tag"
	GotoAlbum    = "album"
	GotoClip     = "clip"
	GotoDaily    = "daily"
)

// FillURI deal app schema.
func FillURI(gt, param string, f func(uri string) string) (uri string) {
	if param == "" {
		return
	}
	switch gt {
	case GotoAv, "":
		uri = "bilibili://video/" + param
	case GotoLive:
		uri = "bilibili://live/" + param
	case GotoBangumi:
		uri = "https://www.bilibili.com/bangumi/play/ss" + param
	case GotoArticle:
		uri = "bilibili://article/" + param
	case GotoGame:
		uri = param
	case GotoAudio:
		uri = "bilibili://music/menu/detail/" + param
	case GotoSong:
		uri = "bilibili://music/detail/" + param
	case GotoAudioTag:
		uri = "bilibili://music/categorydetail/" + param
	case GotoDaily:
		uri = "bilibili://pegasus/list/daily/" + param
	case GotoAlbum:
		uri = "bilibili://album/" + param
	case GotoClip:
		uri = "bilibili://clip/" + param
	case GotoWeb:
		uri = param
	}
	if f != nil {
		uri = f(uri)
	}
	return
}

// AvHandler logic
var AvHandler = func(a *arcwar.Arc) func(uri string) string {
	return func(uri string) string {
		if a == nil {
			return uri
		}
		if a.Dimension.Height != 0 || a.Dimension.Width != 0 {
			return fmt.Sprintf("%s?player_width=%d&player_height=%d&player_rotate=%d", uri, a.Dimension.Width, a.Dimension.Height, a.Dimension.Rotate)
		}
		return uri
	}
}
