# opencc - Golang version OpenCC

## Introduction 介紹
opencc is a golang port of OpenCC([Open Chinese Convert 開放中文轉換](https://github.com/BYVoid/OpenCC/)) which is a project for conversion between Traditional and Simplified Chinese developed by [BYVoid](https://www.byvoid.com/).

opencc stands for "**Go**lang version Open**CC**", it is a total rewrite version of OpenCC in Go. It just borrows the dict files and config files of OpenCC, so it may not produce the same output with the original OpenCC.

## Usage 使用
```go
package main

import (
	"fmt"
	"log"
	"context"

	"go-common/library/text/translate/chinese"
)

func main() {
	chinese.Init()
	in := `请不要怀疑,这是一个由人工智能推荐的频道。`
	out, err := chinese.Convert(context.Background(),in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:%s\n", in, out)
}
// 请不要怀疑,这是一个由人工智能推荐的频道。
// 請不要懷疑,這是一個由人工智慧推薦的頻道。
```


## Conversions
* `s2t` Simplified Chinese to Traditional Chinese
* `t2s` Traditional Chinese to Simplified Chinese
* `s2tw` Simplified Chinese to Traditional Chinese (Taiwan Standard)
* `tw2s` Traditional Chinese (Taiwan Standard) to Simplified Chinese
* `s2hk` Simplified Chinese to Traditional Chinese (Hong Kong Standard)
* `hk2s` Traditional Chinese (Hong Kong Standard) to Simplified Chinese
* `s2twp` Simplified Chinese to Traditional Chinese (Taiwan Standard) with Taiwanese idiom
* `tw2sp` Traditional Chinese (Taiwan Standard) to Simplified Chinese with Mainland Chinese idiom
* `t2tw` Traditional Chinese (OpenCC Standard) to Taiwan Standard
* `t2hk` Traditional Chinese (OpenCC Standard) to Hong Kong Standard
