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

package initupload

import (
	"encoding/json"
	"flag"

	"k8s.io/test-infra/prow/gcsupload"
)

const (
	// JSONConfigEnvVar is the environment variable that
	// utilities expect to find a full JSON configuration
	// in when run.
	JSONConfigEnvVar = "INITUPLOAD_OPTIONS"
)

// NewOptions returns an empty Options with no nil fields
func NewOptions() *Options {
	return &Options{
		Options: gcsupload.NewOptions(),
	}
}

type Options struct {
	*gcsupload.Options

	// Log is the log file to which clone records are written.
	// If unspecified, no clone records are uploaded.
	Log string `json:"log,omitempty"`
}

// ConfigVar exposes the environment variable used
// to store serialized configuration
func (o *Options) ConfigVar() string {
	return JSONConfigEnvVar
}

// LoadConfig loads options from serialized config
func (o *Options) LoadConfig(config string) error {
	return json.Unmarshal([]byte(config), o)
}

// AddFlags binds flags to options
func (o *Options) AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&o.Log, "clone-log", "", "Path to the clone records log")
	o.Options.AddFlags(flags)
}

// Complete internalizes command line arguments
func (o *Options) Complete(args []string) {
	o.Options.Complete(args)
}

// Encode will encode the set of options in the format
// that is expected for the configuration environment variable
func Encode(options Options) (string, error) {
	encoded, err := json.Marshal(options)
	return string(encoded), err
}
