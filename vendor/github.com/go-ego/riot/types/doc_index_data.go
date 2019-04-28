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

// DocIndexData type document Index Data struct
// type DocIndexData DocData
type DocIndexData = DocData

// DocData type document Index Data struct
type DocData struct {
	// 文档全文（必须是 UTF-8 格式），用于生成待索引的关键词
	Content string

	// new 类别
	// Class string
	// new 属性
	Attri interface{}

	// 文档的关键词
	// 当 Content 不为空的时候，优先从 Content 中分词得到关键词。
	// Tokens 存在的意义在于绕过 riot 内置的分词器，在引擎外部
	// 进行分词和预处理。
	// Tokens []*TokenData
	Tokens []TokenData

	// 文档标签（必须是 UTF-8 格式），比如文档的类别属性等，
	// 这些标签并不出现在文档文本中
	Labels []string

	// 文档的评分字段，可以接纳任何类型的结构体
	Fields interface{}
}

// TokenData 文档的一个关键词
type TokenData struct {
	// 关键词的字符串
	Text string

	// 关键词的首字节在文档中出现的位置
	Locations []int
}

// Attri doc attribute
type Attri struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Time   string `json:"time"`
	Ts     int64  `json:"ts"`
}
