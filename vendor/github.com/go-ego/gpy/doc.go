// Copyright 2017 The go-ego Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-ego/ego/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

/*

Package gpy : Chinese Pinyin conversion tool; 汉语拼音转换工具.

Installation:
	go get -u github.com/go-ego/gpy

Usage :

	package main

	import (
		"fmt"

		"github.com/go-ego/gpy"
	)

	func main() {
		hans := "中国人"
		// 默认
		a := gpy.NewArgs()
		fmt.Println(gpy.Pinyin(hans, a))
		// [[zhong] [guo] [ren]]

		// 包含声调
		a.Style = gpy.Tone
		fmt.Println(gpy.Pinyin(hans, a))
		// [[zhōng] [guó] [rén]]

		// 声调用数字表示
		a.Style = gpy.Tone2
		fmt.Println(gpy.Pinyin(hans, a))
		// [[zho1ng] [guo2] [re2n]]

		// 开启多音字模式
		a = gpy.NewArgs()
		a.Heteronym = true
		fmt.Println(gpy.Pinyin(hans, a))
		// [[zhong zhong] [guo] [ren]]
		a.Style = gpy.Tone2
		fmt.Println(gpy.Pinyin(hans, a))
		// [[zho1ng zho4ng] [guo2] [re2n]]
	}
*/
package gpy
