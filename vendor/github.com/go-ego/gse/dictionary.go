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

package gse

import (
	"github.com/go-ego/cedar"
)

// Dictionary 结构体实现了一个字串前缀树，
// 一个分词可能出现在叶子节点也有可能出现在非叶节点
type Dictionary struct {
	trie           *cedar.Cedar // Cedar 前缀树
	maxTokenLen    int          // 词典中最长的分词
	tokens         []Token      // 词典中所有的分词，方便遍历
	totalFrequency int64        // 词典中所有分词的频率之和
}

// NewDict new dictionary
func NewDict() *Dictionary {
	return &Dictionary{trie: cedar.New()}
}

// MaxTokenLen 词典中最长的分词
func (dict *Dictionary) MaxTokenLen() int {
	return dict.maxTokenLen
}

// NumTokens 词典中分词数目
func (dict *Dictionary) NumTokens() int {
	return len(dict.tokens)
}

// TotalFrequency 词典中所有分词的频率之和
func (dict *Dictionary) TotalFrequency() int64 {
	return dict.totalFrequency
}

// addToken 向词典中加入一个分词
func (dict *Dictionary) addToken(token Token) {
	bytes := textSliceToBytes(token.text)
	_, err := dict.trie.Get(bytes)
	if err == nil {
		return
	}

	dict.trie.Insert(bytes, dict.NumTokens())
	dict.tokens = append(dict.tokens, token)
	dict.totalFrequency += int64(token.frequency)
	if len(token.text) > dict.maxTokenLen {
		dict.maxTokenLen = len(token.text)
	}
}

// lookupTokens 在词典中查找和字元组 words 可以前缀匹配的所有分词
// 返回值为找到的分词数
func (dict *Dictionary) lookupTokens(words []Text,
	tokens []*Token) (numOfTokens int) {
	var (
		id, value int
		err       error
	)

	for _, word := range words {
		id, err = dict.trie.Jump(word, id)
		if err != nil {
			break
		}
		value, err = dict.trie.Value(id)
		if err == nil {
			tokens[numOfTokens] = &dict.tokens[value]
			numOfTokens++
		}
	}

	return
}
