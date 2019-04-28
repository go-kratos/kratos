package conf

import xtime "go-common/library/time"

// Search related config def.
type Search struct {
	URL         string
	MainVer     string
	SugNum      int
	SugPGCBuild int
	SugType     string
	Highlight   string         // use highlight or not
	HotwordFre  xtime.Duration // hotword reload frequency
	ResultURL   string
	UserSearch  string
}
