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

package types

import (
	"github.com/go-ego/riot/utils"
)

// SearchResp search response options
type SearchResp struct {
	// 搜索用到的关键词
	Tokens []string

	// 类别
	// Class string

	// 搜索到的文档，已排序
	// Docs []ScoredDoc
	Docs interface{}

	// 搜索是否超时。超时的情况下也可能会返回部分结果
	Timeout bool

	// 搜索到的文档个数。注意这是全部文档中满足条件的个数，可能比返回的文档数要大
	NumDocs int
}

// Content search content
type Content struct {
	// new Content
	Content string

	// new 属性 Attri
	Attri interface{}

	// new 返回评分字段
	Fields interface{}
}

// ScoredDoc scored the document
type ScoredDoc struct {
	DocId uint64

	// new 返回文档 Content
	Content string
	// new 返回文档属性 Attri
	Attri interface{}
	// new 返回评分字段
	Fields interface{}

	// 文档的打分值
	// 搜索结果按照 Scores 的值排序，先按照第一个数排，
	// 如果相同则按照第二个数排序，依次类推。
	Scores []float32

	// 用于生成摘要的关键词在文本中的字节位置，
	// 该切片长度和 SearchResp.Tokens 的长度一样
	// 只有当 IndexType == LocsIndex 时不为空
	TokenSnippetLocs []int

	// 关键词出现的位置
	// 只有当 IndexType == LocsIndex 时不为空
	TokenLocs [][]int
}

// ScoredDocs 为了方便排序
type ScoredDocs []ScoredDoc

func (docs ScoredDocs) Len() int {
	return len(docs)
}

func (docs ScoredDocs) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}

func (docs ScoredDocs) Less(i, j int) bool {
	// 为了从大到小排序，这实际上实现的是 More 的功能
	for iScore := 0; iScore < utils.MinInt(len(docs[i].Scores), len(docs[j].Scores)); iScore++ {
		if docs[i].Scores[iScore] > docs[j].Scores[iScore] {
			return true
		} else if docs[i].Scores[iScore] < docs[j].Scores[iScore] {
			return false
		}
	}
	return len(docs[i].Scores) > len(docs[j].Scores)
}

/*
  ______   .__   __.  __      ____    ____  __   _______
 /  __  \  |  \ |  | |  |     \   \  /   / |  | |       \
|  |  |  | |   \|  | |  |      \   \/   /  |  | |  .--.  |
|  |  |  | |  . `  | |  |       \_    _/   |  | |  |  |  |
|  `--'  | |  |\   | |  `----.    |  |     |  | |  '--'  |
 \______/  |__| \__| |_______|    |__|     |__| |_______/

*/

// ScoredID scored doc only id
type ScoredID struct {
	DocId uint64

	// 文档的打分值
	// 搜索结果按照 Scores 的值排序，先按照第一个数排，
	// 如果相同则按照第二个数排序，依次类推。
	Scores []float32

	// 用于生成摘要的关键词在文本中的字节位置，
	// 该切片长度和 SearchResp.Tokens 的长度一样
	// 只有当 IndexType == LocsIndex 时不为空
	TokenSnippetLocs []int

	// 关键词出现的位置
	// 只有当 IndexType == LocsIndex 时不为空
	TokenLocs [][]int
}

// ScoredIDs 为了方便排序
type ScoredIDs []ScoredID

func (docs ScoredIDs) Len() int {
	return len(docs)
}

func (docs ScoredIDs) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}

func (docs ScoredIDs) Less(i, j int) bool {
	// 为了从大到小排序，这实际上实现的是 More 的功能
	for iScore := 0; iScore < utils.MinInt(len(docs[i].Scores), len(docs[j].Scores)); iScore++ {
		if docs[i].Scores[iScore] > docs[j].Scores[iScore] {
			return true
		} else if docs[i].Scores[iScore] < docs[j].Scores[iScore] {
			return false
		}
	}
	return len(docs[i].Scores) > len(docs[j].Scores)
}
