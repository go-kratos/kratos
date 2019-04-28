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

package sidecar

import (
	"encoding/json"
	"flag"

	"k8s.io/test-infra/prow/gcsupload"
	"k8s.io/test-infra/prow/pod-utils/wrapper"
)

// NewOptions returns an empty Options with no nil fields
func NewOptions() *Options {
	return &Options{
		GcsOptions:     gcsupload.NewOptions(),
		WrapperOptions: &wrapper.Options{},
	}
}

// Options exposes the configuration necessary
// for defining the process being watched and
// where in GCS an upload will land.
type Options struct {
	GcsOptions     *gcsupload.Options `json:"gcs_options"`
	WrapperOptions *wrapper.Options   `json:"wrapper_options"`
}

// Validate ensures that the set of options are
// self-consistent and valid
func (o *Options) Validate() error {
	if err := o.WrapperOptions.Validate(); err != nil {
		return err
	}

	return o.GcsOptions.Validate()
}

const (
	// JSONConfigEnvVar is the environment variable that
	// utilities expect to find a full JSON configuration
	// in when run.
	JSONConfigEnvVar = "SIDECAR_OPTIONS"
)

// ConfigVar exposese the environment variable used
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
	o.GcsOptions.AddFlags(flags)
	o.WrapperOptions.AddFlags(flags)
}

// Complete internalizes command line arguments
func (o *Options) Complete(args []string) {
	o.GcsOptions.Complete(args)
}

// Encode will encode the set of options in the format that
// is expected for the configuration environment variable
func Encode(options Options) (string, error) {
	encoded, err := json.Marshal(options)
	return string(encoded), err
}
