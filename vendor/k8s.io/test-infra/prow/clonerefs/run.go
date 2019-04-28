/*
Copyright 2018 The Kubernetes Authors.

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/kube"
	"k8s.io/test-infra/prow/pod-utils/clone"
)

var cloneFunc = clone.Run

// Run clones the configured refs
func (o Options) Run() error {
	var env []string
	if len(o.KeyFiles) > 0 {
		var err error
		env, err = addSSHKeys(o.KeyFiles)
		if err != nil {
			logrus.WithError(err).Error("Failed to add SSH keys.")
			// Continue on error. Clones will fail with an appropriate error message
			// that initupload can consume whereas quitting without writing the clone
			// record log is silent and results in an errored prow job instead of a
			// failed one.
		}
	}
	if len(o.HostFingerprints) > 0 {
		if err := addHostFingerprints(o.HostFingerprints); err != nil {
			logrus.WithError(err).Error("failed to add host fingerprints")
		}
	}

	var numWorkers int
	if o.MaxParallelWorkers != 0 {
		numWorkers = o.MaxParallelWorkers
	} else {
		numWorkers = len(o.GitRefs)
	}

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	input := make(chan kube.Refs)
	output := make(chan clone.Record, len(o.GitRefs))
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for ref := range input {
				output <- cloneFunc(ref, o.SrcRoot, o.GitUserName, o.GitUserEmail, o.CookiePath, env)
			}
		}()
	}

	for _, ref := range o.GitRefs {
		input <- ref
	}

	close(input)
	wg.Wait()
	close(output)

	var results []clone.Record
	for record := range output {
		results = append(results, record)
	}

	logData, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal clone records: %v", err)
	}

	if err := ioutil.WriteFile(o.Log, logData, 0755); err != nil {
		return fmt.Errorf("failed to write clone records: %v", err)
	}

	return nil
}

func addHostFingerprints(fingerprints []string) error {
	path := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not create/append to %s: %v", path, err)
	}
	if _, err := f.Write([]byte(strings.Join(fingerprints, "\n"))); err != nil {
		return fmt.Errorf("failed to write fingerprints to %s: %v", path, err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close %s: %v", path, err)
	}
	return nil
}

// addSSHKeys will start the ssh-agent and add all the specified
// keys, returning the ssh-agent environment variables for reuse
func addSSHKeys(paths []string) ([]string, error) {
	vars, err := exec.Command("ssh-agent").CombinedOutput()
	if err != nil {
		return []string{}, fmt.Errorf("failed to start ssh-agent: %v", err)
	}
	logrus.Info("Started SSH agent")
	// ssh-agent will output three lines of text, in the form:
	// SSH_AUTH_SOCK=xxx; export SSH_AUTH_SOCK;
	// SSH_AGENT_PID=xxx; export SSH_AGENT_PID;
	// echo Agent pid xxx;
	// We need to parse out the environment variables from that.
	parts := strings.Split(string(vars), ";")
	env := []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[2])}
	for _, keyPath := range paths {
		// we can be given literal paths to keys or paths to dirs
		// that are mounted from a secret, so we need to check which
		// we have
		if err := filepath.Walk(keyPath, func(path string, info os.FileInfo, err error) error {
			if strings.HasPrefix(info.Name(), "..") {
				// kubernetes volumes also include files we
				// should not look be looking into for keys
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if info.IsDir() {
				return nil
			}

			cmd := exec.Command("ssh-add", path)
			cmd.Env = append(cmd.Env, env...)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to add ssh key at %s: %v: %s", path, err, output)
			}
			logrus.Infof("Added SSH key at %s", path)
			return nil
		}); err != nil {
			return env, fmt.Errorf("error walking path %q: %v", keyPath, err)
		}
	}
	return env, nil
}
