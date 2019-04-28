// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package pb

import (
	"github.com/golang/protobuf/proto"
)

// MustMarshal is like proto.Marshal except it panic()'s if the protobuf
// couldn't be serialized.
func MustMarshal(pb proto.Message) []byte {
	b, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return b
}
