/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1alpha1

import (
	"strings"

	"github.com/knative/pkg/apis"
)

type pathTree struct {
	nodeMap map[string]map[string]string
}

// insertNode functions checks the path does not have overlap with existing
// paths in path.nodeMap. If not it creates a key for path and adds
func insertNode(path string, pathtree pathTree) *apis.FieldError {
	err := apis.ErrMultipleOneOf("b.spec.sources.targetPath")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	for nodePath, nodeMap := range pathtree.nodeMap {
		if len(nodeMap) > len(parts) {
			if strings.HasPrefix(nodePath, path) {
				return err
			}
		}

		if len(nodeMap) == len(parts) {
			if path == nodePath {
				return err
			}
		}
		if len(nodeMap) < len(parts) {
			if strings.HasPrefix(path, nodePath) {
				return err
			}
		}
	}
	// path is trimmed with "/"
	addNode(path, pathtree)
	return nil
}

func addNode(path string, tree pathTree) {
	parts := strings.Split(path, "/")
	nm := map[string]string{}

	for _, part := range parts {
		nm[part] = part
	}
	tree.nodeMap[path] = nm
}
