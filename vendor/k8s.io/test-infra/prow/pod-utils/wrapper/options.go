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

package wrapper

import (
	"errors"
	"flag"
)

// Options exposes the configuration options
// used when wrapping test execution
type Options struct {
	// ProcessLog will contain std{out,err} from the
	// wrapped test process
	ProcessLog string `json:"process_log"`

	// MarkerFile will be written with the exit code
	// of the test process or an internal error code
	// if the entrypoint fails.
	MarkerFile string `json:"marker_file"`
}

// AddFlags adds flags to the FlagSet that populate
// the wrapper options struct provided.
func (o *Options) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.ProcessLog, "process-log", "", "path to the log where stdout and stderr are streamed for the process we execute")
	fs.StringVar(&o.MarkerFile, "marker-file", "", "file we write the return code of the process we execute once it has finished running")
}

// Validate ensures that the set of options are
// self-consistent and valid
func (o *Options) Validate() error {
	if o.ProcessLog == "" {
		return errors.New("no log file specified with --process-log")
	}

	if o.MarkerFile == "" {
		return errors.New("no marker file specified with --marker-file")
	}

	return nil
}
