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

// SearchReq search request options
type SearchReq struct {
	// 搜索的短语（必须是 UTF-8 格式），会被分词
	// 当值为空字符串时关键词会从下面的 Tokens 读入
	Text string

	// 关键词（必须是 UTF-8 格式），当 Text 不为空时优先使用 Text
	// 通常你不需要自己指定关键词，除非你运行自己的分词程序
	Tokens []string

	// 文档标签（必须是 UTF-8 格式），标签不存在文档文本中，
	// 但也属于搜索键的一种
	Labels []string

	// 类别
	// Class string

	// 逻辑检索表达式
	Logic Logic

	// 当不为 nil 时，仅从这些 DocIds 包含的键中搜索（忽略值）
	DocIds map[uint64]bool

	// 排序选项
	RankOpts *RankOpts

	// 超时，单位毫秒（千分之一秒）。此值小于等于零时不设超时。
	// 搜索超时的情况下仍有可能返回部分排序结果。
	Timeout int

	// 设为 true 时仅统计搜索到的文档个数，不返回具体的文档
	CountDocsOnly bool

	// 不排序，对于可在引擎外部（比如客户端）排序情况适用
	// 对返回文档很多的情况打开此选项可以有效节省时间
	Orderless bool
}

// RankOpts rank options
type RankOpts struct {
	// 文档的评分规则，值为 nil 时使用 Engine 初始化时设定的规则
	ScoringCriteria ScoringCriteria

	// 默认情况下（ReverseOrder = false）按照分数从大到小排序，否则从小到大排序
	ReverseOrder bool

	// 从第几条结果开始输出
	OutputOffset int

	// 最大输出的搜索结果数，为 0 时无限制
	MaxOutputs int
}

// Logic logic options
type Logic struct {
	// return all doc
	// All bool

	// 与查询, 必须都存在
	Must bool

	// 或查询, 有一个存在即可
	Should bool

	// 非查询, 不包含
	NotIn bool

	LogicExpr LogicExpr
}

// LogicExpr logic expression options
type LogicExpr struct {

	// 与查询, 必须都存在
	MustLabels []string

	// 或查询, 有一个存在即可
	ShouldLabels []string

	// 非查询, 不包含
	NotInLabels []string
}
