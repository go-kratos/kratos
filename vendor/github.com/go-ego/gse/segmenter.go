// Copyright 2013 Hui Chen
// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

/*

package gse Go efficient text segmentation, Go 语言分词
*/

package gse

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	version string = "v0.10.0.106, Danube River!"

	minTokenFrequency = 2 // 仅从字典文件中读取大于等于此频率的分词
)

// GetVersion get the gse version
func GetVersion() string {
	return version
}

// Segmenter 分词器结构体
type Segmenter struct {
	dict *Dictionary
}

// jumper 该结构体用于记录 Viterbi 算法中某字元处的向前分词跳转信息
type jumper struct {
	minDistance float32
	token       *Token
}

// Dictionary 返回分词器使用的词典
func (seg *Segmenter) Dictionary() *Dictionary {
	return seg.dict
}

// getCurrentFilePath get current file path
func getCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(1)
	return filePath
}

// Read read the dict flie
func (seg *Segmenter) Read(file string) error {
	log.Printf("Load the gse dictionary: \"%s\" ", file)
	dictFile, err := os.Open(file)
	if err != nil {
		log.Printf("Could not load dictionaries: \"%s\", %v \n", file, err)
		return err
	}
	defer dictFile.Close()

	reader := bufio.NewReader(dictFile)
	var (
		text      string
		freqText  string
		frequency int
		pos       string
	)

	// 逐行读入分词
	line := 0
	for {
		line++
		size, fsErr := fmt.Fscanln(reader, &text, &freqText, &pos)
		if fsErr != nil {
			if fsErr == io.EOF {
				// End of file
				break
			}

			if size > 0 {
				log.Printf("File '%v' line \"%v\" read error: %v, skip",
					file, line, fsErr.Error())
			} else {
				log.Printf("File '%v' line \"%v\" is empty, read error: %v, skip",
					file, line, fsErr.Error())
			}
		}

		if size == 0 {
			// 文件结束或错误行
			// break
			continue
		} else if size < 2 {
			// 无效行
			continue
		} else if size == 2 {
			// 没有词性标注时设为空字符串
			pos = ""
		}

		// 解析词频
		var err error
		frequency, err = strconv.Atoi(freqText)
		if err != nil {
			continue
		}

		// 过滤频率太小的词
		if frequency < minTokenFrequency {
			continue
		}
		// 过滤, 降低词频
		if len([]rune(text)) < 2 {
			// continue
			frequency = 2
		}

		// 将分词添加到字典中
		words := splitTextToWords([]byte(text))
		token := Token{text: words, frequency: frequency, pos: pos}
		seg.dict.addToken(token)
	}

	return nil
}

// DictPaths get the dict's paths
func DictPaths(dictDir, filePath string) (files []string) {
	var dictPath string

	if filePath == "en" {
		return
	}

	if filePath == "zh" {
		dictPath = path.Join(dictDir, "dict/dictionary.txt")
		files = []string{dictPath}

		return
	}

	if filePath == "jp" {
		dictPath = path.Join(dictDir, "dict/jp/dict.txt")
		files = []string{dictPath}

		return
	}

	// if strings.Contains(filePath, ",") {
	fileName := strings.Split(filePath, ",")
	for i := 0; i < len(fileName); i++ {
		if fileName[i] == "jp" {
			dictPath = path.Join(dictDir, "dict/jp/dict.txt")
		}

		if fileName[i] == "zh" {
			dictPath = path.Join(dictDir, "dict/dictionary.txt")
		}

		// if str[i] == "ti" {
		// }

		dictName := fileName[i] != "en" && fileName[i] != "zh" &&
			fileName[i] != "jp" && fileName[i] != "ti"

		if dictName {
			dictPath = fileName[i]
		}

		if dictPath != "" {
			files = append(files, dictPath)
		}
	}
	// }
	log.Println("Dict files path: ", files)

	return
}

// IsJp is jp char return true
func IsJp(segText string) bool {
	for _, r := range segText {
		jp := unicode.Is(unicode.Scripts["Hiragana"], r) ||
			unicode.Is(unicode.Scripts["Katakana"], r)
		if jp {
			return true
		}
	}
	return false
}

// SegToken add segmenter token
func (seg *Segmenter) SegToken() {
	// 计算每个分词的路径值，路径值含义见 Token 结构体的注释
	logTotalFrequency := float32(math.Log2(float64(seg.dict.totalFrequency)))
	for i := range seg.dict.tokens {
		token := &seg.dict.tokens[i]
		token.distance = logTotalFrequency - float32(math.Log2(float64(token.frequency)))
	}

	// 对每个分词进行细致划分，用于搜索引擎模式，
	// 该模式用法见 Token 结构体的注释。
	for i := range seg.dict.tokens {
		token := &seg.dict.tokens[i]
		segments := seg.segmentWords(token.text, true)

		// 计算需要添加的子分词数目
		numTokensToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			// if len(segments[iToken].token.text) > 1 {
			// 略去字元长度为一的分词
			// TODO: 这值得进一步推敲，特别是当字典中有英文复合词的时候
			if len(segments[iToken].token.text) > 0 {
				hasJp := false
				if len(segments[iToken].token.text) == 1 {
					segText := string(segments[iToken].token.text[0])
					hasJp = IsJp(segText)
				}

				if !hasJp {
					numTokensToAdd++
				}
			}
		}
		token.segments = make([]*Segment, numTokensToAdd)

		// 添加子分词
		iSegmentsToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			// if len(segments[iToken].token.text) > 1 {
			if len(segments[iToken].token.text) > 0 {
				hasJp := false
				if len(segments[iToken].token.text) == 1 {
					segText := string(segments[iToken].token.text[0])
					hasJp = IsJp(segText)
				}

				if !hasJp {
					token.segments[iSegmentsToAdd] = &segments[iToken]
					iSegmentsToAdd++
				}
			}
		}
	}

}

// LoadDict load the dictionary from the file
//
// The format of the dictionary is (one for each participle):
//	participle text, frequency, part of speech
//
// Can load multiple dictionary files, the file name separated by ","
// the front of the dictionary preferentially load the participle,
//	such as: "user_dictionary.txt,common_dictionary.txt"
// When a participle appears both in the user dictionary and
// in the `common dictionary`, the `user dictionary` is given priority.
//
// 从文件中载入词典
//
// 可以载入多个词典文件，文件名用 "," 分隔，排在前面的词典优先载入分词，比如:
// 	"用户词典.txt,通用词典.txt"
// 当一个分词既出现在用户词典也出现在 `通用词典` 中，则优先使用 `用户词典`。
//
// 词典的格式为（每个分词一行）：
//	分词文本 频率 词性
func (seg *Segmenter) LoadDict(files ...string) error {
	seg.dict = NewDict()

	var (
		dictDir  = path.Join(path.Dir(getCurrentFilePath()), "data")
		dictPath string
		// load     bool
	)

	if len(files) > 0 {
		dictFiles := DictPaths(dictDir, files[0])
		if len(dictFiles) > 0 {
			// load = true
			// files = dictFiles
			for i := 0; i < len(dictFiles); i++ {
				err := seg.Read(dictFiles[i])
				if err != nil {
					return err
				}
			}
		}
	}

	if len(files) == 0 {
		dictPath = path.Join(dictDir, "dict/dictionary.txt")
		// files = []string{dictPath}
		err := seg.Read(dictPath)
		if err != nil {
			return err
		}
	}

	// if files[0] != "" && files[0] != "en" && !load {
	// 	for _, file := range strings.Split(files[0], ",") {
	// 		// for _, file := range files {
	// 		err := seg.Read(file)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	seg.SegToken()
	log.Println("Gse dictionary loaded finished.")

	return nil
}

// Segment 对文本分词
//
// 输入参数：
//	bytes	UTF8 文本的字节数组
//
// 输出：
//	[]Segment	划分的分词
func (seg *Segmenter) Segment(bytes []byte) []Segment {
	return seg.internalSegment(bytes, false)
}

// ModeSegment segment using search mode if searchMode is true
func (seg *Segmenter) ModeSegment(bytes []byte, searchMode ...bool) []Segment {
	var mode bool
	if len(searchMode) > 0 {
		mode = searchMode[0]
	}

	return seg.internalSegment(bytes, mode)
}

// Slice use modeSegment segment retrun []string
// using search mode if searchMode is true
func (seg *Segmenter) Slice(bytes []byte, searchMode ...bool) []string {
	segs := seg.ModeSegment(bytes, searchMode...)
	return ToSlice(segs, searchMode...)
}

// Slice use modeSegment segment retrun string
// using search mode if searchMode is true
func (seg *Segmenter) String(bytes []byte, searchMode ...bool) string {
	segs := seg.ModeSegment(bytes, searchMode...)
	return ToString(segs, searchMode...)
}

func (seg *Segmenter) internalSegment(bytes []byte, searchMode bool) []Segment {
	// 处理特殊情况
	if len(bytes) == 0 {
		// return []Segment{}
		return nil
	}

	// 划分字元
	text := splitTextToWords(bytes)

	return seg.segmentWords(text, searchMode)
}

func (seg *Segmenter) segmentWords(text []Text, searchMode bool) []Segment {
	// 搜索模式下该分词已无继续划分可能的情况
	if searchMode && len(text) == 1 {
		return nil
	}

	// jumpers 定义了每个字元处的向前跳转信息，
	// 包括这个跳转对应的分词，
	// 以及从文本段开始到该字元的最短路径值
	jumpers := make([]jumper, len(text))

	if seg.dict == nil {
		return nil
	}

	tokens := make([]*Token, seg.dict.maxTokenLen)
	for current := 0; current < len(text); current++ {
		// 找到前一个字元处的最短路径，以便计算后续路径值
		var baseDistance float32
		if current == 0 {
			// 当本字元在文本首部时，基础距离应该是零
			baseDistance = 0
		} else {
			baseDistance = jumpers[current-1].minDistance
		}

		// 寻找所有以当前字元开头的分词
		numTokens := seg.dict.lookupTokens(
			text[current:minInt(current+seg.dict.maxTokenLen, len(text))], tokens)

		// 对所有可能的分词，更新分词结束字元处的跳转信息
		for iToken := 0; iToken < numTokens; iToken++ {
			location := current + len(tokens[iToken].text) - 1
			if !searchMode || current != 0 || location != len(text)-1 {
				updateJumper(&jumpers[location], baseDistance, tokens[iToken])
			}
		}

		// 当前字元没有对应分词时补加一个伪分词
		if numTokens == 0 || len(tokens[0].text) > 1 {
			updateJumper(&jumpers[current], baseDistance,
				&Token{text: []Text{text[current]}, frequency: 1, distance: 32, pos: "x"})
		}
	}

	// 从后向前扫描第一遍得到需要添加的分词数目
	numSeg := 0
	for index := len(text) - 1; index >= 0; {
		location := index - len(jumpers[index].token.text) + 1
		numSeg++
		index = location - 1
	}

	// 从后向前扫描第二遍添加分词到最终结果
	outputSegments := make([]Segment, numSeg)
	for index := len(text) - 1; index >= 0; {
		location := index - len(jumpers[index].token.text) + 1
		numSeg--
		outputSegments[numSeg].token = jumpers[index].token
		index = location - 1
	}

	// 计算各个分词的字节位置
	bytePosition := 0
	for iSeg := 0; iSeg < len(outputSegments); iSeg++ {
		outputSegments[iSeg].start = bytePosition
		bytePosition += textSliceByteLen(outputSegments[iSeg].token.text)
		outputSegments[iSeg].end = bytePosition
	}
	return outputSegments
}

// updateJumper 更新跳转信息:
// 	1. 当该位置从未被访问过时 (jumper.minDistance 为零的情况)，或者
//	2. 当该位置的当前最短路径大于新的最短路径时
// 将当前位置的最短路径值更新为 baseDistance 加上新分词的概率
func updateJumper(jumper *jumper, baseDistance float32, token *Token) {
	newDistance := baseDistance + token.distance
	if jumper.minDistance == 0 || jumper.minDistance > newDistance {
		jumper.minDistance = newDistance
		jumper.token = token
	}
}

// minInt 取两整数较小值
func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// maxInt 取两整数较大值
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// splitTextToWords 将文本划分成字元
func splitTextToWords(text Text) []Text {
	output := make([]Text, 0, len(text)/3)
	current := 0
	inAlphanumeric := true
	alphanumericStart := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			// 当前是拉丁字母或数字（非中日韩文字）
			if !inAlphanumeric {
				alphanumericStart = current
				inAlphanumeric = true
			}
		} else {
			if inAlphanumeric {
				inAlphanumeric = false
				if current != 0 {
					output = append(output, toLower(text[alphanumericStart:current]))
				}
			}
			output = append(output, text[current:current+size])
		}
		current += size
	}

	// 处理最后一个字元是英文的情况
	if inAlphanumeric {
		if current != 0 {
			output = append(output, toLower(text[alphanumericStart:current]))
		}
	}

	return output
}

// toLower 将英文词转化为小写
func toLower(text []byte) []byte {
	output := make([]byte, len(text))
	for i, t := range text {
		if t >= 'A' && t <= 'Z' {
			output[i] = t - 'A' + 'a'
		} else {
			output[i] = t
		}
	}
	return output
}
