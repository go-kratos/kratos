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

// 这些常数定义了反向索引表存储的数据类型
const (
	// DocIdsIndex 仅存储文档的 docId
	DocIdsIndex = 0

	// FrequenciesIndex 存储关键词的词频，用于计算BM25
	FrequenciesIndex = 1

	// LocsIndex 存储关键词在文档中出现的具体字节位置（可能有多个）
	// 如果你希望得到关键词紧邻度数据，必须使用 LocsIndex 类型的索引
	LocsIndex = 2

	// 默认插入索引表文档 CACHE SIZE
	defaultDocCacheSize = 300000
)

// IndexerOpts 初始化索引器选项
type IndexerOpts struct {
	// 索引表的类型，见上面的常数
	IndexType int

	// 待插入索引表文档 CACHE SIZE
	DocCacheSize int

	// BM25 参数
	BM25Parameters *BM25Parameters
}

// BM25Parameters 见http://en.wikipedia.org/wiki/Okapi_BM25
// 默认值见 engine_init_options.go
type BM25Parameters struct {
	K1 float32
	B  float32
}

// Init init IndexerOpts
func (options *IndexerOpts) Init() {
	if options.DocCacheSize == 0 {
		options.DocCacheSize = defaultDocCacheSize
	}
}
