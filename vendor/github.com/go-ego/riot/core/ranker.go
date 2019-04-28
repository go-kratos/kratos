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

package core

import (
	// "fmt"
	"log"
	"sort"
	"sync"

	"github.com/go-ego/riot/types"
	"github.com/go-ego/riot/utils"
)

// Ranker ranker
type Ranker struct {
	idOnly bool

	lock struct {
		sync.RWMutex
		fields map[uint64]interface{}
		docs   map[uint64]bool
		// new
		content map[uint64]string
		attri   map[uint64]interface{}
	}

	initialized bool
}

// Init init ranker
func (ranker *Ranker) Init(onlyID ...bool) {
	if ranker.initialized == true {
		log.Fatal("The Ranker can not be initialized twice.")
	}
	ranker.initialized = true

	ranker.lock.fields = make(map[uint64]interface{})
	ranker.lock.docs = make(map[uint64]bool)

	if len(onlyID) > 0 {
		ranker.idOnly = onlyID[0]
	}

	if !ranker.idOnly {
		// new
		ranker.lock.content = make(map[uint64]string)
		ranker.lock.attri = make(map[uint64]interface{})
	}
}

// AddDoc add doc
// 给某个文档添加评分字段
func (ranker *Ranker) AddDoc(
	// docId uint64, fields interface{}, content string, attri interface{}) {
	docId uint64, fields interface{}, content ...interface{}) {
	if ranker.initialized == false {
		log.Fatal("The Ranker has not been initialized.")
	}

	ranker.lock.Lock()
	ranker.lock.fields[docId] = fields
	ranker.lock.docs[docId] = true

	if !ranker.idOnly {
		// new
		if len(content) > 0 {
			ranker.lock.content[docId] = content[0].(string)
		}

		if len(content) > 1 {
			ranker.lock.attri[docId] = content[1]
			// ranker.lock.attri[docId] = attri
		}
	}

	ranker.lock.Unlock()
}

// RemoveDoc 删除某个文档的评分字段
func (ranker *Ranker) RemoveDoc(docId uint64) {
	if ranker.initialized == false {
		log.Fatal("The Ranker has not been initialized.")
	}

	ranker.lock.Lock()
	delete(ranker.lock.fields, docId)
	delete(ranker.lock.docs, docId)

	if !ranker.idOnly {
		// new
		delete(ranker.lock.content, docId)
		delete(ranker.lock.attri, docId)
	}

	ranker.lock.Unlock()
}

// RankDocId rank docs by types.ScoredIDs
func (ranker *Ranker) RankDocId(docs []types.IndexedDoc,
	options types.RankOpts, countDocsOnly bool) (types.ScoredIDs, int) {
	var outputDocs types.ScoredIDs
	numDocs := 0

	for _, d := range docs {
		ranker.lock.RLock()
		// 判断 doc 是否存在
		if _, ok := ranker.lock.docs[d.DocId]; ok {

			fs := ranker.lock.fields[d.DocId]

			ranker.lock.RUnlock()
			// 计算评分并剔除没有分值的文档
			scores := options.ScoringCriteria.Score(d, fs)
			if len(scores) > 0 {
				if !countDocsOnly {
					outputDocs = append(outputDocs, types.ScoredID{
						DocId:            d.DocId,
						Scores:           scores,
						TokenSnippetLocs: d.TokenSnippetLocs,
						TokenLocs:        d.TokenLocs})
				}
				numDocs++
			}
		} else {
			ranker.lock.RUnlock()
		}
	}

	// 排序
	if !countDocsOnly {
		if options.ReverseOrder {
			sort.Sort(sort.Reverse(outputDocs))
		} else {
			sort.Sort(outputDocs)
		}
		// 当用户要求只返回部分结果时返回部分结果
		var start, end int
		if options.MaxOutputs != 0 {
			start = utils.MinInt(options.OutputOffset, len(outputDocs))
			end = utils.MinInt(options.OutputOffset+options.MaxOutputs, len(outputDocs))
		} else {
			start = utils.MinInt(options.OutputOffset, len(outputDocs))
			end = len(outputDocs)
		}
		return outputDocs[start:end], numDocs
	}

	return outputDocs, numDocs
}

// RankDocs rank docs by types.ScoredDocs
func (ranker *Ranker) RankDocs(docs []types.IndexedDoc,
	options types.RankOpts, countDocsOnly bool) (types.ScoredDocs, int) {
	var outputDocs types.ScoredDocs
	numDocs := 0

	for _, d := range docs {
		ranker.lock.RLock()
		// 判断 doc 是否存在
		if _, ok := ranker.lock.docs[d.DocId]; ok {

			fs := ranker.lock.fields[d.DocId]
			content := ranker.lock.content[d.DocId]
			attri := ranker.lock.attri[d.DocId]

			ranker.lock.RUnlock()
			// 计算评分并剔除没有分值的文档
			scores := options.ScoringCriteria.Score(d, fs)
			if len(scores) > 0 {
				if !countDocsOnly {
					outputDocs = append(outputDocs, types.ScoredDoc{
						DocId: d.DocId,
						// new
						Fields:  fs,
						Content: content,
						Attri:   attri,
						//
						Scores:           scores,
						TokenSnippetLocs: d.TokenSnippetLocs,
						TokenLocs:        d.TokenLocs})
				}
				numDocs++
			}
		} else {
			ranker.lock.RUnlock()
		}
	}

	// 排序
	if !countDocsOnly {
		if options.ReverseOrder {
			sort.Sort(sort.Reverse(outputDocs))
		} else {
			sort.Sort(outputDocs)
		}
		// 当用户要求只返回部分结果时返回部分结果
		var start, end int
		if options.MaxOutputs != 0 {
			start = utils.MinInt(options.OutputOffset, len(outputDocs))
			end = utils.MinInt(options.OutputOffset+options.MaxOutputs, len(outputDocs))
		} else {
			start = utils.MinInt(options.OutputOffset, len(outputDocs))
			end = len(outputDocs)
		}
		return outputDocs[start:end], numDocs
	}

	return outputDocs, numDocs
}

// Rank rank docs
// 给文档评分并排序
func (ranker *Ranker) Rank(docs []types.IndexedDoc,
	options types.RankOpts, countDocsOnly bool) (
	interface{}, int) {

	if ranker.initialized == false {
		log.Fatal("The Ranker has not been initialized.")
	}

	// 对每个文档评分
	if ranker.idOnly {
		outputDocs, numDocs := ranker.RankDocId(docs, options, countDocsOnly)
		return outputDocs, numDocs
	}

	outputDocs, numDocs := ranker.RankDocs(docs, options, countDocsOnly)
	return outputDocs, numDocs
}
