package gse

import (
	"bytes"
	"fmt"
)

// ToString segments to string  输出分词结果为字符串
//
// 有两种输出模式，以 "中华人民共和国" 为例
//
//  普通模式（searchMode=false）输出一个分词 "中华人民共和国/ns "
//  搜索模式（searchMode=true） 输出普通模式的再细致切分：
//      "中华/nz 人民/n 共和/nz 国/n 共和国/ns 人民共和国/nt 中华人民共和国/ns "
//
// 默认 searchMode=false
// 搜索模式主要用于给搜索引擎提供尽可能多的关键字，详情请见 Token 结构体的注释。
func ToString(segs []Segment, searchMode ...bool) (output string) {
	var mode bool
	if len(searchMode) > 0 {
		mode = searchMode[0]
	}

	if mode {
		for _, seg := range segs {
			output += tokenToString(seg.token)
		}
		return
	}

	for _, seg := range segs {
		output += fmt.Sprintf("%s/%s ",
			textSliceToString(seg.token.text), seg.token.pos)
	}
	return
}

func tokenToBytes(token *Token) (output []byte) {
	for _, s := range token.segments {
		output = append(output, tokenToBytes(s.token)...)
	}
	output = append(output, []byte(fmt.Sprintf("%s/%s ",
		textSliceToString(token.text), token.pos))...)

	return
}

func tokenToString(token *Token) (output string) {
	hasOnlyTerminalToken := true
	for _, s := range token.segments {
		if len(s.token.segments) > 1 || IsJp(string(s.token.text[0])) {
			hasOnlyTerminalToken = false
		}

		if !hasOnlyTerminalToken {
			if s != nil {
				output += tokenToString(s.token)
			}
		}
	}

	output += fmt.Sprintf("%s/%s ", textSliceToString(token.text), token.pos)
	return
}

// ToSlice segments to slice 输出分词结果到一个字符串 slice
//
// 有两种输出模式，以 "中华人民共和国" 为例
//
//  普通模式（searchMode=false）输出一个分词"[中华人民共和国]"
//  搜索模式（searchMode=true） 输出普通模式的再细致切分：
//      "[中华 人民 共和 国 共和国 人民共和国 中华人民共和国]"
//
// 默认 searchMode=false
// 搜索模式主要用于给搜索引擎提供尽可能多的关键字，详情请见Token结构体的注释。
func ToSlice(segs []Segment, searchMode ...bool) (output []string) {
	var mode bool
	if len(searchMode) > 0 {
		mode = searchMode[0]
	}

	if mode {
		for _, seg := range segs {
			output = append(output, tokenToSlice(seg.token)...)
		}
		return
	}

	for _, seg := range segs {
		output = append(output, seg.token.Text())
	}
	return
}

func tokenToSlice(token *Token) (output []string) {
	hasOnlyTerminalToken := true
	for _, s := range token.segments {
		if len(s.token.segments) > 1 || IsJp(string(s.token.text[0])) {
			hasOnlyTerminalToken = false
		}

		if !hasOnlyTerminalToken {
			output = append(output, tokenToSlice(s.token)...)
		}
	}

	output = append(output, textSliceToString(token.text))
	return output
}

// 将多个字元拼接一个字符串输出
func textToString(text []Text) string {
	var output string
	for _, word := range text {
		output += string(word)
	}
	return output
}

// 将多个字元拼接一个字符串输出
func textSliceToString(text []Text) string {
	return Join(text)
}

// 返回多个字元的字节总长度
func textSliceByteLen(text []Text) (length int) {
	for _, word := range text {
		length += len(word)
	}
	return
}

func textSliceToBytes(text []Text) []byte {
	var buf bytes.Buffer
	for _, word := range text {
		buf.Write(word)
	}
	return buf.Bytes()
}

// Join is better string splicing
func Join(text []Text) string {
	switch len(text) {
	case 0:
		return ""
	case 1:
		return string(text[0])
	case 2:
		// Special case for common small values.
		// Remove if github.com/golang/go/issues/6714 is fixed
		return string(text[0]) + string(text[1])
	case 3:
		// Special case for common small values.
		// Remove if #6714 is fixed
		return string(text[0]) + string(text[1]) + string(text[2])
	}
	n := 0
	for i := 0; i < len(text); i++ {
		n += len(text[i])
	}

	b := make([]byte, n)
	bp := copy(b, text[0])
	for _, str := range text[1:] {
		bp += copy(b[bp:], str)
	}
	return string(b)
}
