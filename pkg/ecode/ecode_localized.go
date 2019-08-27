package ecode

// The lang used following the specification defined at
// http://www.rfc-editor.org/rfc/bcp/bcp47.txt.
// Examples are: "en-US", "fr-CH", "es-MX"
const (
	LangDefault = "default"
	LangZhCN    = "zh-CN"
	LangZhTW    = "zh-TW"
	LangZhHK    = "zh-HK"
	LangEnUS    = "en-US"
	LangJaJP    = "ja-JP"
)

type localizedCode struct {
	code    int
	message string
}

func (l localizedCode) Error() string {
	return l.message
}

func (l localizedCode) Code() int {
	return l.code
}

func (l localizedCode) Message() string {
	return l.message
}

func (l localizedCode) Details() []interface{} {
	return nil
}

func (l localizedCode) Equal(err error) bool {
	return EqualError(l, err)
}

func fixLanguage(lang string) string {
	switch lang {
	case LangZhCN:
		lang = LangDefault
	case LangZhHK:
		lang = LangZhTW
	default:
	}
	return lang
}

// LocalizedError error
func LocalizedError(err error, langs []string) Codes {
	ec := Cause(err)
	if len(langs) == 0 {
		return ec
	}
	if localizedCodes, ok := _localizedMessage.Load().(map[int]map[string]string); ok {
		if codes, ok := localizedCodes[ec.Code()]; ok {
			for _, lang := range langs {
				if message, ok := codes[fixLanguage(lang)]; ok {
					return localizedCode{code: ec.Code(), message: message}
				}
			}
		}
	}
	return ec
}
