package model

import "go-common/library/conf/env"

const (
	GotoAv      = "av"
	GotoArticle = "article"
	GotoClip    = "clip"
	GotoAlbum   = "album"
	GotoAudio   = "audio"
)

// env sh001 run
func EnvRun() (res bool) {
	var _zone = "sh001"
	if env.Zone == _zone {
		return true
	}
	return false
}
