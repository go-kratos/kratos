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

// ScoringCriteria 评分规则通用接口
type ScoringCriteria interface {
	// 给一个文档评分，文档排序时先用第一个分值比较，如果
	// 分值相同则转移到第二个分值，以此类推。
	// 返回空切片表明该文档应该从最终排序结果中剔除。
	Score(doc IndexedDoc, fields interface{}) []float32
}

// RankByBM25 一个简单的评分规则，文档分数为BM25
type RankByBM25 struct {
}

// Score score
func (rule RankByBM25) Score(doc IndexedDoc, fields interface{}) []float32 {
	return []float32{doc.BM25}
}
