package archive

// DescFormat is archive type.
type DescFormat struct {
	ID         int64  `json:"id"`
	Copyright  int8   `json:"copyright"`
	TypeID     int64  `json:"typeid"`
	Components string `json:"components"`
	Lang       int8   `json:"lang"`
}

// AppFormat app format.
type AppFormat struct {
	ID        int64 `json:"id"`
	Copyright int8  `json:"copyright"`
	TypeID    int64 `json:"typeid"`
}

//ToLang str to int8.
func ToLang(langStr string) (lang int8) {
	if langStr == "" {
		langStr = "ch"
	}
	switch langStr {
	case "ch":
		lang = 0
	case "en":
		lang = 1
	case "jp":
		lang = 2
	}
	return
}
