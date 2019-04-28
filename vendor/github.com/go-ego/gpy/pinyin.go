package gpy

import (
	// "fmt"
	"regexp"
	"strings"
	"unicode"
)

// Meta
const (
	version = "0.10.0.34"
	// License   = "MIT"
)

// GetVersion get version
func GetVersion() string {
	return version
}

// 拼音风格(推荐)
const (
	Normal      = 0 // 普通风格，不带声调（默认风格）。如： zhong guo
	Tone        = 1 // 声调风格1，拼音声调在韵母第一个字母上。如： zhōng guó
	Tone2       = 2 // 声调风格2，即拼音声调在各个韵母之后，用数字 [1-4] 进行表示。如： zho1ng guo2
	Tone3       = 8 // 声调风格3，即拼音声调在各个拼音之后，用数字 [1-4] 进行表示。如： zhong1 guo2
	Initials    = 3 // 声母风格，只返回各个拼音的声母部分。如： zh g
	FirstLetter = 4 // 首字母风格，只返回拼音的首字母部分。如： z g
	Finals      = 5 // 韵母风格，只返回各个拼音的韵母部分，不带声调。如： ong uo
	FinalsTone  = 6 // 韵母风格1，带声调，声调在韵母第一个字母上。如： ōng uó
	FinalsTone2 = 7 // 韵母风格2，带声调，声调在各个韵母之后，用数字 [1-4] 进行表示。如： o1ng uo2
	FinalsTone3 = 9 // 韵母风格3，带声调，声调在各个拼音之后，用数字 [1-4] 进行表示。如： ong1 uo2
)

// 拼音风格(兼容之前的版本)
// const (
// 	NORMAL       = Normal
// 	TONE         = Tone
// 	TONE2        = Tone2
// 	INITIALS     = Initials
// 	FIRST_LETTER = FirstLetter
// 	FINALS       = Finals
// 	FINALS_TONE  = FinalsTone
// 	FINALS_TONE2 = FinalsTone2
// )

var (
	// 声母表
	initialArray = strings.Split(
		"b,p,m,f,d,t,n,l,g,k,h,j,q,x,r,zh,ch,sh,z,c,s",
		",",
	)

	// 所有带声调的字符
	rePhoneticSymbolSource = func(m map[string]string) string {
		s := ""
		for k := range m {
			s = s + k
		}
		return s
	}(phoneticSymbol)
)

var (
	// 匹配带声调字符的正则表达式
	rePhoneticSymbol = regexp.MustCompile("[" + rePhoneticSymbolSource + "]")

	// 匹配使用数字标识声调的字符的正则表达式
	reTone2 = regexp.MustCompile("([aeoiuvnm])([1-4])$")

	// 匹配 Tone2 中标识韵母声调的正则表达式
	reTone3 = regexp.MustCompile("^([a-z]+)([1-4])([a-z]*)$")
)

// Args 配置信息
type Args struct {
	Style     int    // 拼音风格（默认： Normal)
	Heteronym bool   // 是否启用多音字模式（默认：禁用）
	Separator string // Slug 中使用的分隔符（默认：-)

	// 处理没有拼音的字符（默认忽略没有拼音的字符）
	// 函数返回的 slice 的长度为0 则表示忽略这个字符
	Fallback func(r rune, a Args) []string
}

var (
	// Style 默认配置：风格
	Style = Normal

	// Heteronym 默认配置：是否启用多音字模式
	Heteronym = false

	// Separator 默认配置： `Slug` 中 Join 所用的分隔符
	Separator = "-"

	// Fallback 默认配置: 如何处理没有拼音的字符(忽略这个字符)
	Fallback = func(r rune, a Args) []string {
		return []string{}
	}

	finalExceptionsMap = map[string]string{
		"ū": "ǖ",
		"ú": "ǘ",
		"ǔ": "ǚ",
		"ù": "ǜ",
	}

	reFinalExceptions  = regexp.MustCompile("^(j|q|x)(ū|ú|ǔ|ù)$")
	reFinal2Exceptions = regexp.MustCompile("^(j|q|x)u(\\d?)$")
)

// NewArgs 返回包含默认配置的 `Args`
func NewArgs() Args {
	return Args{Style, Heteronym, Separator, Fallback}
}

// 获取单个拼音中的声母
func initial(p string) string {
	s := ""
	for _, v := range initialArray {
		if strings.HasPrefix(p, v) {
			s = v
			break
		}
	}
	return s
}

// 获取单个拼音中的韵母
func final(p string) string {
	n := initial(p)
	if n == "" {
		return handleYW(p)
	}

	// 特例 j/q/x
	matches := reFinalExceptions.FindStringSubmatch(p)
	// jū -> jǖ
	if len(matches) == 3 && matches[1] != "" && matches[2] != "" {
		v, _ := finalExceptionsMap[matches[2]]
		return v
	}
	// ju -> jv, ju1 -> jv1
	p = reFinal2Exceptions.ReplaceAllString(p, "${1}v$2")
	return strings.Join(strings.SplitN(p, n, 2), "")
}

// 处理 y, w
func handleYW(p string) string {
	// 特例 y/w
	if strings.HasPrefix(p, "yu") {
		p = "v" + p[2:] // yu -> v
	} else if strings.HasPrefix(p, "yi") {
		p = p[1:] // yi -> i
	} else if strings.HasPrefix(p, "y") {
		p = "i" + p[1:] // y -> i
	} else if strings.HasPrefix(p, "wu") {
		p = p[1:] // wu -> u
	} else if strings.HasPrefix(p, "w") {
		p = "u" + p[1:] // w -> u
	}
	return p
}

func toFixed(p string, a Args) string {
	if a.Style == Initials {
		return initial(p)
	}
	origP := p

	// 替换拼音中的带声调字符
	py := rePhoneticSymbol.ReplaceAllStringFunc(p, func(m string) string {
		symbol, _ := phoneticSymbol[m]
		switch a.Style {
		// 不包含声调
		case Normal, FirstLetter, Finals:
			// 去掉声调: a1 -> a
			m = reTone2.ReplaceAllString(symbol, "$1")
		case Tone2, FinalsTone2, Tone3, FinalsTone3:
			// 返回使用数字标识声调的字符
			m = symbol
		default:
			// 声调在头上
		}
		return m
	})

	switch a.Style {
	// 将声调移动到最后
	case Tone3, FinalsTone3:
		py = reTone3.ReplaceAllString(py, "$1$3$2")
	}
	switch a.Style {
	// 首字母
	case FirstLetter:
		py = py[:1]
	// 韵母
	case Finals, FinalsTone, FinalsTone2, FinalsTone3:
		// 转换为 []rune unicode 编码用于获取第一个拼音字符
		// 因为 string 是 utf-8 编码不方便获取第一个拼音字符
		rs := []rune(origP)
		switch string(rs[0]) {
		// 因为鼻音没有声母所以不需要去掉声母部分
		case "ḿ", "ń", "ň", "ǹ":
		default:
			py = final(py)
		}
	}
	return py
}

func applyStyle(p []string, a Args) []string {
	newP := []string{}
	for _, v := range p {
		newP = append(newP, toFixed(v, a))
	}
	return newP
}

// SinglePinyin 把单个 `rune` 类型的汉字转换为拼音.
func SinglePinyin(r rune, a Args) []string {
	if a.Fallback == nil {
		a.Fallback = Fallback
	}
	value, ok := PinyinDict[int(r)]
	pys := []string{}
	if ok {
		pys = strings.Split(value, ",")
	} else {
		pys = a.Fallback(r, a)
	}
	if len(pys) > 0 {
		if !a.Heteronym {
			pys = pys[:1]
		}

		return applyStyle(pys, a)
	}
	return pys
}

// IsChineseChar to determine whether the Chinese string
// 判断是否为中文字符串
func IsChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}

// HanPinyin 汉字转拼音，支持多音字模式.
func HanPinyin(s string, a Args) [][]string {
	pys := [][]string{}
	for _, r := range s {
		py := SinglePinyin(r, a)
		if len(py) > 0 {
			pys = append(pys, py)
		}
	}
	return pys
}

// Pinyin 汉字转拼音，支持多音字模式和拼音与英文等字母混合.
func Pinyin(s string, a Args) [][]string {
	pys := [][]string{}
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			// if r {
			// }
			py := SinglePinyin(r, a)
			if len(py) > 0 {
				pys = append(pys, py)
			}
		}
		// else {
		// 	py := strings.Split(s, " ")
		// 	fmt.Println(py)
		// }
	}

	py := strings.Split(s, " ")
	for i := 0; i < len(py); i++ {
		var (
			pyarr []string
			cs    int64
		)

		for _, r := range py[i] {
			if unicode.Is(unicode.Scripts["Han"], r) {
				// continue
				cs++
			}
		}
		if cs == 0 {
			pyarr = append(pyarr, py[i])
			pys = append(pys, pyarr)
		}
	}

	return pys
}

// LazyPinyin 汉字转拼音，与 `Pinyin` 的区别是：
// 返回值类型不同，并且不支持多音字模式，每个汉字只取第一个音.
func LazyPinyin(s string, a Args) []string {
	a.Heteronym = false
	pys := []string{}
	for _, v := range Pinyin(s, a) {
		pys = append(pys, v[0])
	}
	return pys
}

// Slug join `LazyPinyin` 的返回值.
// 建议改用 https://github.com/mozillazg/go-slugify
func Slug(s string, a Args) string {
	separator := a.Separator
	return strings.Join(LazyPinyin(s, a), separator)
}

// Convert 跟 Pinyin 的唯一区别就是 a 参数可以是 nil
func Convert(s string, a *Args) [][]string {
	if a == nil {
		args := NewArgs()
		a = &args
	}
	return Pinyin(s, *a)
}

// LazyConvert 跟 LazyPinyin 的唯一区别就是 a 参数可以是 nil
func LazyConvert(s string, a *Args) []string {
	if a == nil {
		args := NewArgs()
		a = &args
	}
	return LazyPinyin(s, *a)
}
