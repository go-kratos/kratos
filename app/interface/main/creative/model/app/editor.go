package app

import "fmt"

// EditorData str
type EditorData struct {
	AID  int64
	CID  int64
	Type int8
	Data string
}

// Editor str
type Editor struct {
	CID    int64 `json:"cid"`
	UpFrom int8  `json:"upfrom"` // filled by backend
	// ids set
	Filters     interface{} `json:"filters"`          // 滤镜
	Fonts       interface{} `json:"fonts"`            //字体
	Subtitles   interface{} `json:"subtitles"`        //字幕
	Bgms        interface{} `json:"bgms"`             //bgm
	Stickers    interface{} `json:"stickers"`         //3d拍摄贴纸
	VStickers   interface{} `json:"videoup_stickers"` //2d投稿贴纸
	Transitions interface{} `json:"trans"`            //视频转场特效
	// add from app535
	Themes     interface{} `json:"themes"`     //编辑器的主题使用相关
	Cooperates interface{} `json:"cooperates"` //拍摄之稿件合拍
	// switch env 0/1
	AudioRecord  int8 `json:"audio_record"`  //录音
	Camera       int8 `json:"camera"`        //拍摄
	Speed        int8 `json:"speed"`         //变速
	CameraRotate int8 `json:"camera_rotate"` //摄像头翻转
}

// ParseThemes fn
func (e *Editor) ParseThemes() (valStr string) {
	Themes, ok := e.Themes.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Themes))
	} else if e.Themes != nil {
		valStr = fmt.Sprintf("%v", e.Themes)
	}
	return
}

// ParseCooperates fn
func (e *Editor) ParseCooperates() (valStr string) {
	Cooperates, ok := e.Cooperates.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Cooperates))
	} else if e.Cooperates != nil {
		valStr = fmt.Sprintf("%v", e.Cooperates)
	}
	return
}

// ParseFilters fn
func (e *Editor) ParseFilters() (valStr string) {
	Filters, ok := e.Filters.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Filters))
	} else if e.Filters != nil {
		valStr = fmt.Sprintf("%v", e.Filters)
	}
	return
}

// ParseFonts fn
func (e *Editor) ParseFonts() (valStr string) {
	Fonts, ok := e.Fonts.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Fonts))
	} else if e.Fonts != nil {
		valStr = fmt.Sprintf("%v", e.Fonts)
	}
	return
}

// ParseSubtitles fn
func (e *Editor) ParseSubtitles() (valStr string) {
	Subtitles, ok := e.Subtitles.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Subtitles))
	} else if e.Subtitles != nil {
		valStr = fmt.Sprintf("%v", e.Subtitles)
	}
	return
}

// ParseBgms fn
func (e *Editor) ParseBgms() (valStr string) {
	Bgms, ok := e.Bgms.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Bgms))
	} else if e.Bgms != nil {
		valStr = fmt.Sprintf("%v", e.Bgms)
	}
	return
}

// ParseStickers fn
func (e *Editor) ParseStickers() (valStr string) {
	Stickers, ok := e.Stickers.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Stickers))
	} else if e.Stickers != nil {
		valStr = fmt.Sprintf("%v", e.Stickers)
	}
	return
}

// ParseVStickers fn
func (e *Editor) ParseVStickers() (valStr string) {
	VStickers, ok := e.VStickers.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(VStickers))
	} else if e.VStickers != nil {
		valStr = fmt.Sprintf("%v", e.VStickers)
	}
	return
}

// ParseTransitions fn
func (e *Editor) ParseTransitions() (valStr string) {
	Transitions, ok := e.Transitions.(float64)
	if ok {
		valStr = fmt.Sprintf("%d", int64(Transitions))
	} else if e.Transitions != nil {
		valStr = fmt.Sprintf("%v", e.Transitions)
	}
	return
}
