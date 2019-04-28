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

Package types is riot types
*/
package types

// DocIndex document's index
type DocIndex struct {
	// DocId 文本的 DocId
	DocId uint64

	// TokenLen 文本的关键词长
	TokenLen float32

	// Keywords 加入的索引键
	Keywords []KeywordIndex
}

// KeywordIndex 反向索引项，这实际上标注了一个（搜索键，文档）对。
type KeywordIndex struct {
	// Text 搜索键的 UTF-8 文本
	Text string

	// Frequency 搜索键词频
	Frequency float32

	// Starts 搜索键在文档中的起始字节位置，按照升序排列
	Starts []int
}

// IndexedDoc 索引器返回结果
type IndexedDoc struct {
	// DocId document id
	DocId uint64

	// BM25，仅当索引类型为 FrequenciesIndex 或者 LocsIndex 时返回有效值
	BM25 float32

	// TokenProximity 关键词在文档中的紧邻距离，
	// 紧邻距离的含义见 computeTokenProximity 的注释。
	// 仅当索引类型为 LocsIndex 时返回有效值。
	TokenProximity int32

	// TokenSnippetLocs 紧邻距离计算得到的关键词位置，
	// 和 Lookup 函数输入 tokens 的长度一样且一一对应。
	// 仅当索引类型为 LocsIndex 时返回有效值。
	TokenSnippetLocs []int

	// TokenLocs 关键词在文本中的具体位置。
	// 仅当索引类型为 LocsIndex 时返回有效值。
	TokenLocs [][]int
}

// DocsIndex 方便批量加入文档索引
type DocsIndex []*DocIndex

func (docs DocsIndex) Len() int {
	return len(docs)
}

func (docs DocsIndex) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}

func (docs DocsIndex) Less(i, j int) bool {
	return docs[i].DocId < docs[j].DocId
}

// DocsId 方便批量删除文档索引
type DocsId []uint64

func (docs DocsId) Len() int {
	return len(docs)
}

func (docs DocsId) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}

func (docs DocsId) Less(i, j int) bool {
	return docs[i] < docs[j]
}
