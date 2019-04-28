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

package riot

import (
	"bufio"
	"log"
	"os"
)

// StopTokens stop tokens map
type StopTokens struct {
	stopTokens map[string]bool
}

// Init 从 stopTokenFile 中读入停用词，一个词一行
// 文档索引建立时会跳过这些停用词
func (st *StopTokens) Init(stopTokenFile string) {
	st.stopTokens = make(map[string]bool)
	if stopTokenFile == "" {
		return
	}

	file, err := os.Open(stopTokenFile)
	if err != nil {
		log.Fatal("Open stop token file error: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			st.stopTokens[text] = true
		}
	}

}

// IsStopToken to determine whether to stop token
func (st *StopTokens) IsStopToken(token string) bool {
	_, found := st.stopTokens[token]
	return found
}
