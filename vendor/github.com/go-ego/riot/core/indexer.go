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

Package core is riot core
*/
package core

import (
	"log"
	"math"
	"sort"
	"sync"

	"github.com/go-ego/riot/types"
	"github.com/go-ego/riot/utils"
)

// Indexer 索引器
type Indexer struct {
	// 从搜索键到文档列表的反向索引
	// 加了读写锁以保证读写安全
	tableLock struct {
		sync.RWMutex
		table     map[string]*KeywordIndices
		docsState map[uint64]int // nil: 表示无状态记录，0: 存在于索引中，1: 等待删除，2: 等待加入
	}

	addCacheLock struct {
		sync.RWMutex
		addCachePointer int
		addCache        types.DocsIndex
	}

	removeCacheLock struct {
		sync.RWMutex
		removeCachePointer int
		removeCache        types.DocsId
	}

	initOptions types.IndexerOpts
	initialized bool

	// 这实际上是总文档数的一个近似
	numDocs uint64

	// 所有被索引文本的总关键词数
	totalTokenLen float32

	// 每个文档的关键词长度
	docTokenLens map[uint64]float32
}

// KeywordIndices 反向索引表的一行，收集了一个搜索键出现的所有文档，按照DocId从小到大排序。
type KeywordIndices struct {
	// 下面的切片是否为空，取决于初始化时IndexType的值
	docIds      []uint64  // 全部类型都有
	frequencies []float32 // IndexType == FrequenciesIndex
	locations   [][]int   // IndexType == LocsIndex
}

// Init 初始化索引器
func (indexer *Indexer) Init(options types.IndexerOpts) {
	if indexer.initialized == true {
		log.Fatal("The Indexer can not be initialized twice.")
	}
	options.Init()
	indexer.initOptions = options
	indexer.initialized = true

	indexer.tableLock.table = make(map[string]*KeywordIndices)
	indexer.tableLock.docsState = make(map[uint64]int)
	indexer.addCacheLock.addCache = make(
		[]*types.DocIndex, indexer.initOptions.DocCacheSize)

	indexer.removeCacheLock.removeCache = make(
		[]uint64, indexer.initOptions.DocCacheSize*2)
	indexer.docTokenLens = make(map[uint64]float32)
}

// getDocId 从 KeywordIndices 中得到第i个文档的 DocId
func (indexer *Indexer) getDocId(ti *KeywordIndices, i int) uint64 {
	return ti.docIds[i]
}

// HasDoc doc is exist return true
func (indexer *Indexer) HasDoc(docId uint64) bool {
	docState, ok := indexer.tableLock.docsState[docId]
	if ok && docState == 0 {
		return true
	}

	return false
}

// getIndexLen 得到 KeywordIndices 中文档总数
func (indexer *Indexer) getIndexLen(ti *KeywordIndices) int {
	return len(ti.docIds)
}

// AddDocToCache 向 ADDCACHE 中加入一个文档
func (indexer *Indexer) AddDocToCache(doc *types.DocIndex, forceUpdate bool) {
	if indexer.initialized == false {
		log.Fatal("The Indexer has not been initialized.")
	}

	indexer.addCacheLock.Lock()
	if doc != nil {
		indexer.addCacheLock.addCache[indexer.addCacheLock.addCachePointer] = doc
		indexer.addCacheLock.addCachePointer++
	}

	if indexer.addCacheLock.addCachePointer >= indexer.initOptions.DocCacheSize ||
		forceUpdate {
		indexer.tableLock.Lock()

		position := 0
		for i := 0; i < indexer.addCacheLock.addCachePointer; i++ {
			docIndex := indexer.addCacheLock.addCache[i]

			docState, ok := indexer.tableLock.docsState[docIndex.DocId]
			if ok && docState <= 1 {
				// ok && docState == 0 表示存在于索引中，需先删除再添加
				// ok && docState == 1 表示不一定存在于索引中，等待删除，需先删除再添加
				if position != i {
					indexer.addCacheLock.addCache[position], indexer.addCacheLock.addCache[i] =
						indexer.addCacheLock.addCache[i], indexer.addCacheLock.addCache[position]
				}
				if docState == 0 {
					indexer.removeCacheLock.Lock()
					indexer.removeCacheLock.removeCache[indexer.removeCacheLock.removeCachePointer] =
						docIndex.DocId
					indexer.removeCacheLock.removeCachePointer++
					indexer.removeCacheLock.Unlock()

					indexer.tableLock.docsState[docIndex.DocId] = 1
					indexer.numDocs--
				}
				position++
			} else if !ok {
				indexer.tableLock.docsState[docIndex.DocId] = 2
			}
		}

		indexer.tableLock.Unlock()
		if indexer.RemoveDocToCache(0, forceUpdate) {
			// 只有当存在于索引表中的文档已被删除，其才可以重新加入到索引表中
			position = 0
		}

		addCachedDocs := indexer.addCacheLock.addCache[position:indexer.addCacheLock.addCachePointer]
		indexer.addCacheLock.addCachePointer = position

		indexer.addCacheLock.Unlock()
		sort.Sort(addCachedDocs)
		indexer.AddDocs(&addCachedDocs)
	} else {
		indexer.addCacheLock.Unlock()
	}
}

// AddDocs 向反向索引表中加入 ADDCACHE 中所有文档
func (indexer *Indexer) AddDocs(docs *types.DocsIndex) {
	if indexer.initialized == false {
		log.Fatal("The Indexer has not been initialized.")
	}

	indexer.tableLock.Lock()
	defer indexer.tableLock.Unlock()
	indexPointers := make(map[string]int, len(indexer.tableLock.table))

	// DocId 递增顺序遍历插入文档保证索引移动次数最少
	for i, doc := range *docs {
		if i < len(*docs)-1 && (*docs)[i].DocId == (*docs)[i+1].DocId {
			// 如果有重复文档加入，因为稳定排序，只加入最后一个
			continue
		}

		docState, ok := indexer.tableLock.docsState[doc.DocId]
		if ok && docState == 1 {
			// 如果此时 docState 仍为 1，说明该文档需被删除
			// docState 合法状态为 nil & 2，保证一定不会插入已经在索引表中的文档
			continue
		}

		// 更新文档关键词总长度
		if doc.TokenLen != 0 {
			indexer.docTokenLens[doc.DocId] = float32(doc.TokenLen)
			indexer.totalTokenLen += doc.TokenLen
		}

		docIdIsNew := true
		for _, keyword := range doc.Keywords {
			indices, foundKeyword := indexer.tableLock.table[keyword.Text]
			if !foundKeyword {
				// 如果没找到该搜索键则加入
				ti := KeywordIndices{}
				switch indexer.initOptions.IndexType {
				case types.LocsIndex:
					ti.locations = [][]int{keyword.Starts}
				case types.FrequenciesIndex:
					ti.frequencies = []float32{keyword.Frequency}
				}
				ti.docIds = []uint64{doc.DocId}
				indexer.tableLock.table[keyword.Text] = &ti
				continue
			}

			// 查找应该插入的位置，且索引一定不存在
			position, _ := indexer.searchIndex(
				indices, indexPointers[keyword.Text], indexer.getIndexLen(indices)-1, doc.DocId)
			indexPointers[keyword.Text] = position

			switch indexer.initOptions.IndexType {
			case types.LocsIndex:
				indices.locations = append(indices.locations, []int{})
				copy(indices.locations[position+1:], indices.locations[position:])
				indices.locations[position] = keyword.Starts
			case types.FrequenciesIndex:
				indices.frequencies = append(indices.frequencies, float32(0))
				copy(indices.frequencies[position+1:], indices.frequencies[position:])
				indices.frequencies[position] = keyword.Frequency
			}

			indices.docIds = append(indices.docIds, 0)
			copy(indices.docIds[position+1:], indices.docIds[position:])
			indices.docIds[position] = doc.DocId
		}

		// 更新文章状态和总数
		if docIdIsNew {
			indexer.tableLock.docsState[doc.DocId] = 0
			indexer.numDocs++
		}
	}
}

// RemoveDocToCache 向 REMOVECACHE 中加入一个待删除文档
// 返回值表示文档是否在索引表中被删除
func (indexer *Indexer) RemoveDocToCache(docId uint64, forceUpdate bool) bool {
	if indexer.initialized == false {
		log.Fatal("The Indexer has not been initialized.")
	}

	indexer.removeCacheLock.Lock()
	if docId != 0 {
		indexer.tableLock.Lock()
		if docState, ok := indexer.tableLock.docsState[docId]; ok && docState == 0 {
			indexer.removeCacheLock.removeCache[indexer.removeCacheLock.removeCachePointer] = docId
			indexer.removeCacheLock.removeCachePointer++
			indexer.tableLock.docsState[docId] = 1
			indexer.numDocs--
		} else if ok && docState == 2 {
			// 删除一个等待加入的文档
			indexer.tableLock.docsState[docId] = 1
		} else if !ok {
			// 若文档不存在，则无法判断其是否在 addCache 中，需避免这样的操作
		}
		indexer.tableLock.Unlock()
	}

	if indexer.removeCacheLock.removeCachePointer > 0 &&
		(indexer.removeCacheLock.removeCachePointer >= indexer.initOptions.DocCacheSize ||
			forceUpdate) {
		removeCacheddocs := indexer.removeCacheLock.removeCache[:indexer.removeCacheLock.removeCachePointer]
		indexer.removeCacheLock.removeCachePointer = 0
		indexer.removeCacheLock.Unlock()
		sort.Sort(removeCacheddocs)
		indexer.RemoveDocs(&removeCacheddocs)
		return true
	}

	indexer.removeCacheLock.Unlock()
	return false
}

// RemoveDocs 向反向索引表中删除 REMOVECACHE 中所有文档
func (indexer *Indexer) RemoveDocs(docs *types.DocsId) {
	if indexer.initialized == false {
		log.Fatal("The Indexer has not been initialized.")
	}

	indexer.tableLock.Lock()
	defer indexer.tableLock.Unlock()

	// 更新文档关键词总长度，删除文档状态
	for _, docId := range *docs {
		indexer.totalTokenLen -= indexer.docTokenLens[docId]
		delete(indexer.docTokenLens, docId)
		delete(indexer.tableLock.docsState, docId)
	}

	for keyword, indices := range indexer.tableLock.table {
		indicesTop, indicesPointer := 0, 0
		docsPointer := sort.Search(
			len(*docs), func(i int) bool { return (*docs)[i] >= indices.docIds[0] })
		// 双指针扫描，进行批量删除操作
		for docsPointer < len(*docs) && indicesPointer < indexer.getIndexLen(indices) {
			if indices.docIds[indicesPointer] < (*docs)[docsPointer] {
				if indicesTop != indicesPointer {
					switch indexer.initOptions.IndexType {
					case types.LocsIndex:
						indices.locations[indicesTop] = indices.locations[indicesPointer]
					case types.FrequenciesIndex:
						indices.frequencies[indicesTop] = indices.frequencies[indicesPointer]
					}

					indices.docIds[indicesTop] = indices.docIds[indicesPointer]
				}

				indicesTop++
				indicesPointer++
			} else if indices.docIds[indicesPointer] == (*docs)[docsPointer] {
				indicesPointer++
				docsPointer++
			} else {
				docsPointer++
			}
		}
		if indicesTop != indicesPointer {
			switch indexer.initOptions.IndexType {
			case types.LocsIndex:
				indices.locations = append(
					indices.locations[:indicesTop], indices.locations[indicesPointer:]...)
			case types.FrequenciesIndex:
				indices.frequencies = append(
					indices.frequencies[:indicesTop], indices.frequencies[indicesPointer:]...)
			}

			indices.docIds = append(
				indices.docIds[:indicesTop], indices.docIds[indicesPointer:]...)
		}

		if len(indices.docIds) == 0 {
			delete(indexer.tableLock.table, keyword)
		}
	}
}

// Lookup lookup docs
// 查找包含全部搜索键(AND操作)的文档
// 当 docIds 不为 nil 时仅从 docIds 指定的文档中查找
func (indexer *Indexer) Lookup(
	tokens []string, labels []string, docIds map[uint64]bool, countDocsOnly bool,
	logic ...types.Logic) (docs []types.IndexedDoc, numDocs int) {

	if indexer.initialized == false {
		log.Fatal("The Indexer has not been initialized.")
	}

	if indexer.numDocs == 0 {
		return
	}
	numDocs = 0

	// 合并关键词和标签为搜索键
	keywords := make([]string, len(tokens)+len(labels))
	copy(keywords, tokens)
	copy(keywords[len(tokens):], labels)

	if len(logic) > 0 {
		if logic != nil && len(keywords) > 0 && logic[0].Must == true ||
			logic[0].Should == true || logic[0].NotIn == true {

			docs, numDocs = indexer.LogicLookup(
				docIds, countDocsOnly, keywords, logic[0])

			return
		}

		if logic != nil && (len(logic[0].LogicExpr.MustLabels) > 0 ||
			len(logic[0].LogicExpr.ShouldLabels) > 0) &&
			len(logic[0].LogicExpr.NotInLabels) >= 0 {

			docs, numDocs = indexer.LogicLookup(
				docIds, countDocsOnly, keywords, logic[0])

			return
		}
	}

	indexer.tableLock.RLock()
	defer indexer.tableLock.RUnlock()

	table := make([]*KeywordIndices, len(keywords))
	for i, keyword := range keywords {
		indices, found := indexer.tableLock.table[keyword]
		if !found {
			// 当反向索引表中无此搜索键时直接返回
			return
		}
		// 否则加入反向表中
		table[i] = indices
	}

	// 当没有找到时直接返回
	if len(table) == 0 {
		return
	}

	// 归并查找各个搜索键出现文档的交集
	// 从后向前查保证先输出 DocId 较大文档
	indexPointers := make([]int, len(table))
	for iTable := 0; iTable < len(table); iTable++ {
		indexPointers[iTable] = indexer.getIndexLen(table[iTable]) - 1
	}

	// 平均文本关键词长度，用于计算BM25
	avgDocLength := indexer.totalTokenLen / float32(indexer.numDocs)
	for ; indexPointers[0] >= 0; indexPointers[0]-- {
		// 以第一个搜索键出现的文档作为基准，并遍历其他搜索键搜索同一文档
		baseDocId := indexer.getDocId(table[0], indexPointers[0])
		if docIds != nil {
			if _, found := docIds[baseDocId]; !found {
				continue
			}
		}

		iTable := 1
		found := true
		for ; iTable < len(table); iTable++ {
			// 二分法比简单的顺序归并效率高，也有更高效率的算法，
			// 但顺序归并也许是更好的选择，考虑到将来需要用链表重新实现
			// 以避免反向表添加新文档时的写锁。
			// TODO: 进一步研究不同求交集算法的速度和可扩展性。
			position, foundBaseDocId := indexer.searchIndex(table[iTable],
				0, indexPointers[iTable], baseDocId)

			if foundBaseDocId {
				indexPointers[iTable] = position
			} else {
				if position == 0 {
					// 该搜索键中所有的文档 ID 都比 baseDocId 大，因此已经没有
					// 继续查找的必要。
					return
				}

				// 继续下一 indexPointers[0] 的查找
				indexPointers[iTable] = position - 1
				found = false
				break
			}
		}

		if found {
			if docState, ok := indexer.tableLock.docsState[baseDocId]; !ok || docState != 0 {
				continue
			}
			indexedDoc := types.IndexedDoc{}

			// 当为 LocsIndex 时计算关键词紧邻距离
			if indexer.initOptions.IndexType == types.LocsIndex {
				// 计算有多少关键词是带有距离信息的
				numTokensWithLocations := 0
				for i, t := range table[:len(tokens)] {
					if len(t.locations[indexPointers[i]]) > 0 {
						numTokensWithLocations++
					}
				}
				if numTokensWithLocations != len(tokens) {
					if !countDocsOnly {
						docs = append(docs, types.IndexedDoc{
							DocId: baseDocId,
						})
					}
					numDocs++
					//当某个关键字对应多个文档且有 lable 关键字存在时，若直接 break,
					// 将会丢失相当一部分搜索结果
					continue
				}

				// 计算搜索键在文档中的紧邻距离
				tokenProximity, TokenLocs := computeTokenProximity(
					table[:len(tokens)], indexPointers, tokens)

				indexedDoc.TokenProximity = int32(tokenProximity)
				indexedDoc.TokenSnippetLocs = TokenLocs

				// 添加 TokenLocs
				indexedDoc.TokenLocs = make([][]int, len(tokens))
				for i, t := range table[:len(tokens)] {
					indexedDoc.TokenLocs[i] = t.locations[indexPointers[i]]
				}
			}

			// 当为 LocsIndex 或者 FrequenciesIndex 时计算BM25
			if indexer.initOptions.IndexType == types.LocsIndex ||
				indexer.initOptions.IndexType == types.FrequenciesIndex {
				bm25 := float32(0)
				d := indexer.docTokenLens[baseDocId]
				for i, t := range table[:len(tokens)] {
					var frequency float32
					if indexer.initOptions.IndexType == types.LocsIndex {
						frequency = float32(len(t.locations[indexPointers[i]]))
					} else {
						frequency = t.frequencies[indexPointers[i]]
					}

					// 计算 BM25
					if len(t.docIds) > 0 && frequency > 0 &&
						indexer.initOptions.BM25Parameters != nil && avgDocLength != 0 {
						// 带平滑的 idf
						idf := float32(math.Log2(float64(indexer.numDocs)/float64(len(t.docIds)) + 1))
						k1 := indexer.initOptions.BM25Parameters.K1
						b := indexer.initOptions.BM25Parameters.B
						bm25 += idf * frequency * (k1 + 1) / (frequency + k1*(1-b+b*d/avgDocLength))
					}
				}
				indexedDoc.BM25 = float32(bm25)
			}

			indexedDoc.DocId = baseDocId
			if !countDocsOnly {
				docs = append(docs, indexedDoc)
			}
			numDocs++
		}
	}

	return
}

// searchIndex 二分法查找 indices 中某文档的索引项
// 第一个返回参数为找到的位置或需要插入的位置
// 第二个返回参数标明是否找到
func (indexer *Indexer) searchIndex(indices *KeywordIndices,
	start int, end int, docId uint64) (int, bool) {
	// 特殊情况
	if indexer.getIndexLen(indices) == start {
		return start, false
	}
	if docId < indexer.getDocId(indices, start) {
		return start, false
	} else if docId == indexer.getDocId(indices, start) {
		return start, true
	}
	if docId > indexer.getDocId(indices, end) {
		return end + 1, false
	} else if docId == indexer.getDocId(indices, end) {
		return end, true
	}

	// 二分
	var middle int
	for end-start > 1 {
		middle = (start + end) / 2
		if docId == indexer.getDocId(indices, middle) {
			return middle, true
		} else if docId > indexer.getDocId(indices, middle) {
			start = middle
		} else {
			end = middle
		}
	}

	return end, false
}

// computeTokenProximity 计算搜索键在文本中的紧邻距离
//
// 假定第 i 个搜索键首字节出现在文本中的位置为 P_i，长度 L_i
// 紧邻距离计算公式为
//
// 	ArgMin(Sum(Abs(P_(i+1) - P_i - L_i)))
//
// 具体由动态规划实现，依次计算前 i 个 token 在每个出现位置的最优值。
// 选定的 P_i 通过 TokenLocs 参数传回。
func computeTokenProximity(table []*KeywordIndices,
	indexPointers []int, tokens []string) (
	minTokenProximity int, TokenLocs []int) {
	minTokenProximity = -1
	TokenLocs = make([]int, len(tokens))

	var (
		currentLocations, nextLocations []int
		currentMinValues, nextMinValues []int
		path                            [][]int
	)

	// 初始化路径数组
	path = make([][]int, len(tokens))
	for i := 1; i < len(path); i++ {
		path[i] = make([]int, len(table[i].locations[indexPointers[i]]))
	}

	// 动态规划
	currentLocations = table[0].locations[indexPointers[0]]
	currentMinValues = make([]int, len(currentLocations))
	for i := 1; i < len(tokens); i++ {
		nextLocations = table[i].locations[indexPointers[i]]
		nextMinValues = make([]int, len(nextLocations))
		for j := range nextMinValues {
			nextMinValues[j] = -1
		}

		var iNext int
		for iCurrent, currentLocation := range currentLocations {
			if currentMinValues[iCurrent] == -1 {
				continue
			}
			for iNext+1 < len(nextLocations) &&
				nextLocations[iNext+1] < currentLocation {
				iNext++
			}

			update := func(from int, to int) {
				if to >= len(nextLocations) {
					return
				}
				value := currentMinValues[from] +
					utils.AbsInt(nextLocations[to]-currentLocations[from]-len(tokens[i-1]))

				if nextMinValues[to] == -1 || value < nextMinValues[to] {
					nextMinValues[to] = value
					path[i][to] = from
				}
			}

			// 最优解的状态转移只发生在左右最接近的位置
			update(iCurrent, iNext)
			update(iCurrent, iNext+1)
		}

		currentLocations = nextLocations
		currentMinValues = nextMinValues
	}

	// 找出最优解
	var cursor int
	for i, value := range currentMinValues {
		if value == -1 {
			continue
		}
		if minTokenProximity == -1 || value < minTokenProximity {
			minTokenProximity = value
			cursor = i
		}
	}

	// 从路径倒推出最优解的位置
	for i := len(tokens) - 1; i >= 0; i-- {
		if i != len(tokens)-1 {
			cursor = path[i+1][cursor]
		}
		TokenLocs[i] = table[i].locations[indexPointers[i]][cursor]
	}

	return
}

// LogicLookup logic Lookup
func (indexer *Indexer) LogicLookup(
	docIds map[uint64]bool, countDocsOnly bool, LogicExpr []string,
	logic types.Logic) (docs []types.IndexedDoc, numDocs int) {

	indexer.tableLock.RLock()
	defer indexer.tableLock.RUnlock()

	// // 有效性检查, 不允许只出现逻辑非检索, 也不允许与或非都不存在
	// if Logic.Must == true && Logic.Should == true && Logic.NotIn == true {
	// 	return
	// }

	// MustTable 中的搜索键检查
	// 如果存在与搜索键, 则要求所有的与搜索键都有对应的反向表
	MustTable := make([]*KeywordIndices, 0)

	if len(logic.LogicExpr.MustLabels) > 0 {
		LogicExpr = logic.LogicExpr.MustLabels
	}
	if logic.Must == true || len(logic.LogicExpr.MustLabels) > 0 {
		for _, keyword := range LogicExpr {
			indices, found := indexer.tableLock.table[keyword]
			if !found {
				return
			}

			MustTable = append(MustTable, indices)
		}
	}

	// 逻辑或搜索键检查
	// 1. 如果存在逻辑或搜索键, 则至少有一个存在反向表
	// 2. 逻辑或和逻辑与之间是与关系
	ShouldTable := make([]*KeywordIndices, 0)

	if len(logic.LogicExpr.ShouldLabels) > 0 {
		LogicExpr = logic.LogicExpr.ShouldLabels
	}

	if logic.Should == true || len(logic.LogicExpr.ShouldLabels) > 0 {
		for _, keyword := range LogicExpr {
			indices, found := indexer.tableLock.table[keyword]
			if found {
				ShouldTable = append(ShouldTable, indices)
			}
		}
		if len(ShouldTable) == 0 {
			// 如果存在逻辑或搜索键， 但是对应的反向表全部为空， 则返回
			return
		}
	}

	// 逻辑非中的搜索键检查
	// 可以不存在逻辑非搜索（NotInTable为空), 允许逻辑非搜索键对应的反向表为空
	NotInTable := make([]*KeywordIndices, 0)

	if len(logic.LogicExpr.NotInLabels) > 0 {
		LogicExpr = logic.LogicExpr.NotInLabels
	}
	if logic.NotIn == true || len(logic.LogicExpr.NotInLabels) > 0 {
		for _, keyword := range LogicExpr {
			indices, found := indexer.tableLock.table[keyword]
			if found {
				NotInTable = append(NotInTable, indices)
			}
		}
	}

	// 开始检索
	numDocs = 0
	if logic.Must == true || len(logic.LogicExpr.MustLabels) > 0 {
		// 如果存在逻辑与检索
		for idx := indexer.getIndexLen(MustTable[0]) - 1; idx >= 0; idx-- {
			baseDocId := indexer.getDocId(MustTable[0], idx)
			if docIds != nil {
				_, found := docIds[baseDocId]
				if !found {
					continue
				}
			}

			mustFound := indexer.findInMustTable(MustTable[1:], baseDocId)
			shouldFound := indexer.findInShouldTable(ShouldTable, baseDocId)
			notInFound := indexer.findInNotInTable(NotInTable, baseDocId)

			if mustFound && shouldFound && !notInFound {
				indexedDoc := types.IndexedDoc{}
				indexedDoc.DocId = baseDocId
				if !countDocsOnly {
					docs = append(docs, indexedDoc)
				}
				numDocs++
			}
		}
	} else {
		// 不存在逻辑与检索, 则必须存在逻辑或检索
		// 这时进行求并集操作
		if logic.Should == true || len(logic.LogicExpr.ShouldLabels) > 0 {
			docs, numDocs = indexer.unionTable(ShouldTable, NotInTable, countDocsOnly)
		} else {
			uintDocIds := make([]uint64, 0)
			// 当前直接返回 Not 逻辑数据
			for i := 0; i < len(NotInTable); i++ {
				for _, docid := range NotInTable[i].docIds {
					if indexer.findInNotInTable(NotInTable, docid) {
						uintDocIds = append(uintDocIds, docid)
					}
				}
			}

			StableDesc(uintDocIds)

			numDocs = 0
			for _, doc := range uintDocIds {
				indexedDoc := types.IndexedDoc{}
				indexedDoc.DocId = doc
				if !countDocsOnly {
					docs = append(docs, indexedDoc)
				}
				numDocs++
			}
		}

		// fmt.Println(docs, numDocs)
	}

	return
}

// 在逻辑与反向表中对docid进行查找, 若每个反向表都找到,
// 则返回 true, 有一个找不到则返回 false
func (indexer *Indexer) findInMustTable(table []*KeywordIndices, docId uint64) bool {
	for i := 0; i < len(table); i++ {
		_, foundDocId := indexer.searchIndex(table[i],
			0, indexer.getIndexLen(table[i])-1, docId)
		if !foundDocId {
			return false
		}
	}

	return true
}

// 在逻辑或反向表中对 docid 进行查找， 若有一个找到则返回 true,
// 都找不到则返回 false
// 如果 table 为空， 则返回 true
func (indexer *Indexer) findInShouldTable(table []*KeywordIndices, docId uint64) bool {
	for i := 0; i < len(table); i++ {
		_, foundDocId := indexer.searchIndex(table[i],
			0, indexer.getIndexLen(table[i])-1, docId)
		if foundDocId {
			return true
		}
	}

	if len(table) == 0 {
		return true
	}

	return false
}

// findInNotInTable 在逻辑非反向表中对 docid 进行查找,
// 若有一个找到则返回 true, 都找不到则返回 false
// 如果 table 为空, 则返回 false
func (indexer *Indexer) findInNotInTable(table []*KeywordIndices, docId uint64) bool {
	for i := 0; i < len(table); i++ {
		_, foundDocId := indexer.searchIndex(table[i],
			0, indexer.getIndexLen(table[i])-1, docId)
		if foundDocId {
			return true
		}
	}

	return false
}

// unionTable 如果不存在与逻辑检索， 则需要对逻辑或反向表求并集
// 先求差集再求并集， 可以减小内存占用
// docid 要保序
func (indexer *Indexer) unionTable(table []*KeywordIndices,
	notInTable []*KeywordIndices, countDocsOnly bool) (
	docs []types.IndexedDoc, numDocs int) {
	docIds := make([]uint64, 0)
	// 求并集
	for i := 0; i < len(table); i++ {
		for _, docid := range table[i].docIds {
			if !indexer.findInNotInTable(notInTable, docid) {
				found := false
				for _, v := range docIds {
					if v == docid {
						found = true
						break
					}
				}
				if !found {
					docIds = append(docIds, docid)
				}
			}
		}
	}
	// 排序
	// sortUint64.StableDesc(docIds)
	StableDesc(docIds)

	numDocs = 0
	for _, doc := range docIds {
		indexedDoc := types.IndexedDoc{}
		indexedDoc.DocId = doc
		if !countDocsOnly {
			docs = append(docs, indexedDoc)
		}
		numDocs++
	}

	return
}
