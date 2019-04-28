/*
Copyright 2017 The Kubernetes Authors.

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

package clonerefs

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/test-infra/prow/kube"
)

// ParseRefs parses a human-provided string into the repo
// that should be cloned and the refs that need to be
// checked out once it is. The format is:
//   org,repo=base-ref[:base-sha][,pull-id[:pull-sha[:pull-ref]]]...
// For the base ref and pull IDs, a SHA may optionally be
// provided or may be omitted for the latest available SHA.
// Examples:
//   kubernetes,test-infra=master
//   kubernetes,test-infra=master:abcde12
//   kubernetes,test-infra=master:abcde12,34
//   kubernetes,test-infra=master:abcde12,34:fghij56
//   kubernetes,test-infra=master,34:fghij56
//   kubernetes,test-infra=master:abcde12,34:fghij56,78
//   gerrit,test-infra=master:abcde12,34:fghij56:refs/changes/00/123/1
func ParseRefs(value string) (*kube.Refs, error) {
	gitRef := &kube.Refs{}
	values := strings.SplitN(value, "=", 2)
	if len(values) != 2 {
		return gitRef, fmt.Errorf("refspec %s invalid: does not contain '='", value)
	}
	info := values[0]
	allRefs := values[1]

	infoValues := strings.SplitN(info, ",", 2)
	if len(infoValues) != 2 {
		return gitRef, fmt.Errorf("refspec %s invalid: does not contain 'org,repo' as prefix", value)
	}
	gitRef.Org = infoValues[0]
	gitRef.Repo = infoValues[1]

	refValues := strings.Split(allRefs, ",")
	if len(refValues) == 1 && refValues[0] == "" {
		return gitRef, fmt.Errorf("refspec %s invalid: does not contain any refs", value)
	}
	baseRefParts := strings.Split(refValues[0], ":")
	if len(baseRefParts) != 1 && len(baseRefParts) != 2 {
		return gitRef, fmt.Errorf("refspec %s invalid: malformed base ref", refValues[0])
	}
	gitRef.BaseRef = baseRefParts[0]
	if len(baseRefParts) == 2 {
		gitRef.BaseSHA = baseRefParts[1]
	}
	for _, refValue := range refValues[1:] {
		refParts := strings.Split(refValue, ":")
		if len(refParts) == 0 || len(refParts) > 3 {
			return gitRef, fmt.Errorf("refspec %s invalid: malformed pull ref", refValue)
		}
		pullNumber, err := strconv.Atoi(refParts[0])
		if err != nil {
			return gitRef, fmt.Errorf("refspec %s invalid: pull request identifier not a number: %v", refValue, err)
		}
		pullRef := kube.Pull{
			Number: pullNumber,
		}
		if len(refParts) > 1 {
			pullRef.SHA = refParts[1]
		}
		if len(refParts) > 2 {
			pullRef.Ref = refParts[2]
		}
		gitRef.Pulls = append(gitRef.Pulls, pullRef)
	}

	return gitRef, nil
}
