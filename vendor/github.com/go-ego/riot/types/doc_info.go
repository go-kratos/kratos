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
	"sync"
)

// DocInfosShard 文档信息[id]info
type DocInfosShard struct {
	DocInfos map[uint64]*DocInfo
	NumDocs  uint64 // 这实际上是总文档数的一个近似
	sync.RWMutex
}

// DocInfo document info
type DocInfo struct {
	Fields    interface{}
	TokenLens float32
}

/// inverted_index.go

// InvertedIndexShard 反向索引表([关键词]反向索引表)
type InvertedIndexShard struct {
	InvertedIndex map[string]*KeywordIndices
	TotalTokenLen float32 //总关键词数
	sync.RWMutex
}

// KeywordIndices 反向索引表的一行，收集了一个搜索键出现的所有文档，
// 按照 DocId 从小到大排序。
type KeywordIndices struct {
	// 下面的切片是否为空，取决于初始化时 IndexType 的值
	DocIds      []uint64  // 全部类型都有
	Frequencies []float32 // IndexType == FrequenciesIndex
	Locations   [][]int   // IndexType == LocsIndex
}
