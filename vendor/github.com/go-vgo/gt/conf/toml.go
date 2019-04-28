// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build toml

package conf

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Init toml config
func Init(fliePath string, config interface{}) {
	confLock.Lock()
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		log.Println("toml.DecodeFile error: ", err)
		return
	}
	confLock.Unlock()
}
