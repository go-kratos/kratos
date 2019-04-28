// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build !toml

package conf

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

// Init toml config
func Init(filePath string, config interface{}) {
	confLock.Lock()
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("ioutil.ReadFile error: ", err)
	}
	toml.Unmarshal(fileBytes, config)
	confLock.Unlock()
}
