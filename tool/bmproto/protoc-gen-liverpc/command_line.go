// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"strings"
)

type commandLineParams struct {
	importPrefix string            // String to prefix to imported package file names.
	importMap    map[string]string // Mapping from .proto file name to import path.
}

// parseCommandLineParams breaks the comma-separated list of key=value pairs
// in the parameter (a member of the request protobuf) into a key/value map.
// It then sets command line parameter mappings defined by those entries.
func parseCommandLineParams(parameter string) (*commandLineParams, error) {
	ps := make(map[string]string)
	for _, p := range strings.Split(parameter, ",") {
		if p == "" {
			continue
		}
		i := strings.Index(p, "=")
		if i < 0 {
			return nil, fmt.Errorf("invalid parameter %q: expected format of parameter to be k=v", p)
		}
		k := p[0:i]
		v := p[i+1:]
		if v == "" {
			return nil, fmt.Errorf("invalid parameter %q: expected format of parameter to be k=v", k)
		}
		ps[k] = v
	}

	clp := &commandLineParams{
		importMap: make(map[string]string),
	}
	for k, v := range ps {
		switch {
		case k == "import_prefix":
			clp.importPrefix = v
			// Support import map 'M' prefix per https://github.com/golang/protobuf/blob/6fb5325/protoc-gen-go/generator/generator.go#L497.
		case len(k) > 0 && k[0] == 'M':
			clp.importMap[k[1:]] = v // 1 is the length of 'M'.
		case len(k) > 0 && strings.HasPrefix(k, "go_import_mapping@"):
			clp.importMap[k[18:]] = v // 18 is the length of 'go_import_mapping@'.
		default:
			return nil, fmt.Errorf("unknown parameter %q", k)
		}
	}
	return clp, nil
}
