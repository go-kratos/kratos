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

package kube

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	prowJobs = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "prowjobs",
		Help: "Number of prowjobs in the system",
	}, []string{
		// name of the job
		"job_name",
		// type of the prowjob: presubmit, postsubmit, periodic, batch
		"type",
		// state of the prowjob: triggered, pending, success, failure, aborted, error
		"state",
	})
)

func init() {
	prometheus.MustRegister(prowJobs)
}

// GatherProwJobMetrics gathers prometheus metrics for prowjobs.
func GatherProwJobMetrics(pjs []ProwJob) {
	// map of job to job type to state to count
	metricMap := make(map[string]map[string]map[string]float64)

	for _, pj := range pjs {
		if metricMap[pj.Spec.Job] == nil {
			metricMap[pj.Spec.Job] = make(map[string]map[string]float64)
		}
		if metricMap[pj.Spec.Job][string(pj.Spec.Type)] == nil {
			metricMap[pj.Spec.Job][string(pj.Spec.Type)] = make(map[string]float64)
		}
		metricMap[pj.Spec.Job][string(pj.Spec.Type)][string(pj.Status.State)]++
	}

	// This may be racing with the prometheus server but we need to remove
	// stale metrics like triggered or pending jobs that are now complete.
	prowJobs.Reset()

	for job, jobMap := range metricMap {
		for jobType, typeMap := range jobMap {
			for state, count := range typeMap {
				prowJobs.WithLabelValues(job, jobType, state).Set(count)
			}
		}
	}
}
