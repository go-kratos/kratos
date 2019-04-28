package module

// PlayerSha1 one module for sha1
type PlayerSha1 struct {
	// special fields
	PlayerWebDanmakuAutoscaling        bool    `json:"player_web_danmaku_autoscaling,omitempty"`
	PlayerWebHTML5DanmakuRenderingtype string  `json:"player_web_html5_danmaku_renderingtype,omitempty"`
	PlayerAppPlaybackMode              int     `json:"player_app_playback_mode,omitempty"`
	PlayerAppPlaybackBackground        bool    `json:"player_app_playback_background,omitempty"`
	PlayerAppDanmakuStrokewidth        float64 `json:"player_app_danmaku_strokewidth,omitempty"`
	// common fileds
	PlayerDanmakuOpacity         float64 `json:"player_danmaku_opacity,omitempty"`
	PlayerDanmakuSpeed           float64 `json:"player_danmaku_speed,omitempty"`
	PlayerDanmakuDensity         int     `json:"player_danmaku_density,omitempty"`
	PlayerDanmakuScalingfactor   float64 `json:"player_danmaku_scalingfactor,omitempty"`
	PlayerDanmakuStrokestyle     int     `json:"player_danmaku_strokestyle,omitempty"`
	PlayerDanmakuFontname        string  `json:"player_danmaku_fontname,omitempty"`
	PlayerDanmakuFontbold        bool    `json:"player_danmaku_fontbold,omitempty"`
	PlayerDanmakuDefensivebottom bool    `json:"player_danmaku_defensivebottom,omitempty"`
	PlayerDanmakuEnableblocklist bool    `json:"player_danmaku_enableblocklist,omitempty"`
	PlayerDanmakuBlockrepeat     bool    `json:"player_danmaku_blockrepeat,omitempty"`
	PlayerDanmakuBlocktop        bool    `json:"player_danmaku_blocktop,omitempty"`
	PlayerDanmakuBlockscroll     bool    `json:"player_danmaku_blockscroll,omitempty"`
	PlayerDanmakuBlockbottom     bool    `json:"player_danmaku_blockbottom,omitempty"`
	PlayerDanmakuBlockcolorful   bool    `json:"player_danmaku_blockcolorful,omitempty"`
	PlayerDanmakuBlockcommon     bool    `json:"player_danmaku_blockcommon,omitempty"`
	PlayerDanmakuBlocksubtitle   bool    `json:"player_danmaku_blocksubtitle,omitempty"`
	PlayerDanmakuBlockspecial    bool    `json:"player_danmaku_blockspecial,omitempty"`
	PlayerDanmakuDomain          float64 `json:"player_danmaku_domain,omitempty"`
	PlayerSubtitleSwitch         int     `json:"player_subtitle_switch,omitempty"`
}

// Player one module return json
type Player struct {
	// special fields
	PlayerWebDanmakuAutoscaling        bool    `json:"player_web_danmaku_autoscaling"`
	PlayerWebHTML5DanmakuRenderingtype string  `json:"player_web_html5_danmaku_renderingtype"`
	PlayerAppPlaybackMode              int     `json:"player_app_playback_mode"`
	PlayerAppPlaybackBackground        bool    `json:"player_app_playback_background"`
	PlayerAppDanmakuStrokewidth        float64 `json:"player_app_danmaku_strokewidth"`
	// common fileds
	PlayerDanmakuOpacity         float64 `json:"player_danmaku_opacity"`
	PlayerDanmakuSpeed           float64 `json:"player_danmaku_speed"`
	PlayerDanmakuDensity         int     `json:"player_danmaku_density"`
	PlayerDanmakuScalingfactor   float64 `json:"player_danmaku_scalingfactor"`
	PlayerDanmakuStrokestyle     int     `json:"player_danmaku_strokestyle"`
	PlayerDanmakuFontname        string  `json:"player_danmaku_fontname"`
	PlayerDanmakuFontbold        bool    `json:"player_danmaku_fontbold"`
	PlayerDanmakuDefensivebottom bool    `json:"player_danmaku_defensivebottom"`
	PlayerDanmakuEnableblocklist bool    `json:"player_danmaku_enableblocklist"`
	PlayerDanmakuBlockrepeat     bool    `json:"player_danmaku_blockrepeat"`
	PlayerDanmakuBlocktop        bool    `json:"player_danmaku_blocktop"`
	PlayerDanmakuBlockscroll     bool    `json:"player_danmaku_blockscroll"`
	PlayerDanmakuBlockbottom     bool    `json:"player_danmaku_blockbottom"`
	PlayerDanmakuBlockcolorful   bool    `json:"player_danmaku_blockcolorful"`
	PlayerDanmakuBlockcommon     bool    `json:"player_danmaku_blockcommon"`
	PlayerDanmakuBlocksubtitle   bool    `json:"player_danmaku_blocksubtitle"`
	PlayerDanmakuBlockspecial    bool    `json:"player_danmaku_blockspecial"`
	PlayerDanmakuDomain          float64 `json:"player_danmaku_domain"`
	PlayerSubtitleSwitch         int     `json:"player_subtitle_switch"`
}
