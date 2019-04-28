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

package config

import (
	"fmt"
	"regexp"
	"time"

	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	"k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/test-infra/prow/kube"
)

// Preset is intended to match the k8s' PodPreset feature, and may be removed
// if that feature goes beta.
type Preset struct {
	Labels       map[string]string `json:"labels"`
	Env          []v1.EnvVar       `json:"env"`
	Volumes      []v1.Volume       `json:"volumes"`
	VolumeMounts []v1.VolumeMount  `json:"volumeMounts"`
}

func mergePreset(preset Preset, labels map[string]string, pod *v1.PodSpec) error {
	if pod == nil {
		return nil
	}
	for l, v := range preset.Labels {
		if v2, ok := labels[l]; !ok || v2 != v {
			return nil
		}
	}
	for _, e1 := range preset.Env {
		for i := range pod.Containers {
			for _, e2 := range pod.Containers[i].Env {
				if e1.Name == e2.Name {
					return fmt.Errorf("env var duplicated in pod spec: %s", e1.Name)
				}
			}
			pod.Containers[i].Env = append(pod.Containers[i].Env, e1)
		}
	}
	for _, v1 := range preset.Volumes {
		for _, v2 := range pod.Volumes {
			if v1.Name == v2.Name {
				return fmt.Errorf("volume duplicated in pod spec: %s", v1.Name)
			}
		}
		pod.Volumes = append(pod.Volumes, v1)
	}
	for _, vm1 := range preset.VolumeMounts {
		for i := range pod.Containers {
			for _, vm2 := range pod.Containers[i].VolumeMounts {
				if vm1.Name == vm2.Name {
					return fmt.Errorf("volume mount duplicated in pod spec: %s", vm1.Name)
				}
			}
			pod.Containers[i].VolumeMounts = append(pod.Containers[i].VolumeMounts, vm1)
		}
	}
	return nil
}

// JobBase contains attributes common to all job types
type JobBase struct {
	// The name of the job.
	// e.g. pull-test-infra-bazel-build
	Name string `json:"name"`
	// Labels are added to prowjobs and pods created for this job.
	Labels map[string]string `json:"labels,omitempty"`
	// MaximumConcurrency of this job, 0 implies no limit.
	MaxConcurrency int `json:"max_concurrency,omitempty"`
	// Agent that will take care of running this job.
	Agent string `json:"agent"`
	// Cluster is the alias of the cluster to run this job in.
	// (Default: kube.DefaultClusterAlias)
	Cluster string `json:"cluster,omitempty"`
	// Namespace is the namespace in which pods schedule.
	//   nil: results in config.PodNamespace (aka pod default)
	//   empty: results in config.ProwJobNamespace (aka same as prowjob)
	Namespace *string `json:"namespace,omitempty"`
	// ErrorOnEviction indicates that the ProwJob should be completed and given
	// the ErrorState status if the pod that is executing the job is evicted.
	// If this field is unspecified or false, a new pod will be created to replace
	// the evicted one.
	ErrorOnEviction bool `json:"error_on_eviction,omitempty"`
	// SourcePath contains the path where this job is defined
	SourcePath string `json:"-"`
	// Spec is the Kubernetes pod spec used if Agent is kubernetes.
	Spec *v1.PodSpec `json:"spec,omitempty"`
	// BuildSpec is the Knative build spec used if Agent is knative-build.
	BuildSpec *buildv1alpha1.BuildSpec `json:"build_spec,omitempty"`

	UtilityConfig
}

// Presubmit runs on PRs.
type Presubmit struct {
	JobBase

	// AlwaysRun automatically for every PR, or only when a comment triggers it.
	AlwaysRun bool `json:"always_run"`
	// RunIfChanged automatically run if the PR modifies a file that matches this regex.
	RunIfChanged string `json:"run_if_changed,omitempty"`
	// TrustedLabels automatically run if the PR has label in TrustedLabels
	TrustedLabels []string `json:"trusted_labels,omitempty"`
	// UntrustedLabels automatically not run if the PR has label in UntrustedLabels
	UntrustedLabels []string `json:"untrusted_labels,omitempty"`
	// RunPRPushed automatically run if the source branch pushed
	RunPRPushed bool `json:"run_pr_pushed"`

	// Context is the name of the GitHub status context for the job.
	Context string `json:"context"`
	// Optional indicates that the job's status context should not be required for merge.
	Optional bool `json:"optional,omitempty"`
	// SkipReport skips commenting and setting status on GitHub.
	SkipReport bool `json:"skip_report,omitempty"`

	// Trigger is the regular expression to trigger the job.
	// e.g. `@k8s-bot e2e test this`
	// RerunCommand must also be specified if this field is specified.
	// (Default: `(?m)^/test (?:.*? )?<job name>(?: .*?)?$`)
	Trigger string `json:"trigger"`
	// The RerunCommand to give users. Must match Trigger.
	// Trigger must also be specified if this field is specified.
	// (Default: `/test <job name>`)
	RerunCommand string `json:"rerun_command"`

	// RunAfterSuccess is a list of jobs to run after successfully running this one.
	RunAfterSuccess []Presubmit `json:"run_after_success,omitempty"`

	Brancher

	// We'll set these when we load it.
	re        *regexp.Regexp // from Trigger.
	reChanges *regexp.Regexp // from RunIfChanged
}

// Postsubmit runs on push events.
type Postsubmit struct {
	JobBase

	RegexpChangeMatcher

	Brancher

	// Run these jobs after successfully running this one.
	RunAfterSuccess []Postsubmit `json:"run_after_success,omitempty"`
}

// Periodic runs on a timer.
type Periodic struct {
	JobBase

	// (deprecated)Interval to wait between two runs of the job.
	Interval string `json:"interval"`
	// Cron representation of job trigger time
	Cron string `json:"cron"`
	// Tags for config entries
	Tags []string `json:"tags,omitempty"`
	// Run these jobs after successfully running this one.
	RunAfterSuccess []Periodic `json:"run_after_success,omitempty"`

	interval time.Duration
}

// SetInterval updates interval, the frequency duration it runs.
func (p *Periodic) SetInterval(d time.Duration) {
	p.interval = d
}

// GetInterval returns interval, the frequency duration it runs.
func (p *Periodic) GetInterval() time.Duration {
	return p.interval
}

// RegexpChangeMatcher is for code shared between jobs that run only when certain files are changed.
type RegexpChangeMatcher struct {
	// RunIfChanged defines a regex used to select which subset of file changes should trigger this job.
	// If any file in the changeset matches this regex, the job will be triggered
	RunIfChanged string         `json:"run_if_changed,omitempty"`
	reChanges    *regexp.Regexp // from RunIfChanged
}

// RunsAgainstChanges returns true if any of the changed input paths match the run_if_changed regex.
func (cm RegexpChangeMatcher) RunsAgainstChanges(changes []string) bool {
	if cm.RunIfChanged == "" {
		return true
	}
	for _, change := range changes {
		if cm.reChanges.MatchString(change) {
			return true
		}
	}
	return false
}

// Brancher is for shared code between jobs that only run against certain
// branches. An empty brancher runs against all branches.
type Brancher struct {
	// Do not run against these branches. Default is no branches.
	SkipBranches []string `json:"skip_branches,omitempty"`
	// Only run against these branches. Default is all branches.
	Branches []string `json:"branches,omitempty"`

	// We'll set these when we load it.
	re     *regexp.Regexp
	reSkip *regexp.Regexp
}

// RunsAgainstAllBranch returns true if there are both branches and skip_branches are unset
func (br Brancher) RunsAgainstAllBranch() bool {
	return len(br.SkipBranches) == 0 && len(br.Branches) == 0
}

// RunsAgainstBranch returns true if the input branch matches, given the whitelist/blacklist.
func (br Brancher) RunsAgainstBranch(branch string) bool {
	if br.RunsAgainstAllBranch() {
		return true
	}

	// Favor SkipBranches over Branches
	if len(br.SkipBranches) != 0 && br.reSkip.MatchString(branch) {
		return false
	}
	if len(br.Branches) == 0 || br.re.MatchString(branch) {
		return true
	}
	return false
}

// Intersects checks if other Brancher would trigger for the same branch.
func (br Brancher) Intersects(other Brancher) bool {
	if br.RunsAgainstAllBranch() || other.RunsAgainstAllBranch() {
		return true
	}
	if len(br.Branches) > 0 {
		baseBranches := sets.NewString(br.Branches...)
		if len(other.Branches) > 0 {
			otherBranches := sets.NewString(other.Branches...)
			if baseBranches.Intersection(otherBranches).Len() > 0 {
				return true
			}
			return false
		}
		if !baseBranches.Intersection(sets.NewString(other.SkipBranches...)).Equal(baseBranches) {
			return true
		}
		return false
	}
	if len(other.Branches) == 0 {
		// There can only be one Brancher with skip_branches.
		return true
	}
	return other.Intersects(br)
}

// RunsAgainstChanges returns true if any of the changed input paths match the run_if_changed regex.
func (ps Presubmit) RunsAgainstChanges(changes []string) bool {
	for _, change := range changes {
		if ps.reChanges.MatchString(change) {
			return true
		}
	}
	return false
}

// TriggerMatches returns true if the comment body should trigger this presubmit.
//
// This is usually a /test foo string.
func (ps Presubmit) TriggerMatches(body string) bool {
	return ps.re.MatchString(body)
}

// ContextRequired checks whether a context is required from github points of view (required check).
func (ps Presubmit) ContextRequired() bool {
	if ps.Optional || ps.SkipReport {
		return false
	}
	return true
}

// ChangedFilesProvider returns a slice of modified files.
type ChangedFilesProvider func() ([]string, error)

func matching(j Presubmit, body string, testAll bool) []Presubmit {
	// When matching ignore whether the job runs for the branch or whether the job runs for the
	// PR's changes. Even if the job doesn't run, it still matches the PR and may need to be marked
	// as skipped on github.
	var result []Presubmit
	if (testAll && (j.AlwaysRun || j.RunIfChanged != "")) || j.TriggerMatches(body) {
		result = append(result, j)
	}
	for _, child := range j.RunAfterSuccess {
		result = append(result, matching(child, body, testAll)...)
	}
	return result
}

// MatchingPresubmits returns a slice of presubmits to trigger based on the repo and a comment text.
func (c *JobConfig) MatchingPresubmits(fullRepoName, body string, testAll bool) []Presubmit {
	var result []Presubmit
	if jobs, ok := c.Presubmits[fullRepoName]; ok {
		for _, job := range jobs {
			result = append(result, matching(job, body, testAll)...)
		}
	}
	return result
}

// UtilityConfig holds decoration metadata, such as how to clone and additional containers/etc
type UtilityConfig struct {
	// Decorate determines if we decorate the PodSpec or not
	Decorate bool `json:"decorate,omitempty"`

	// PathAlias is the location under <root-dir>/src
	// where the repository under test is cloned. If this
	// is not set, <root-dir>/src/github.com/org/repo will
	// be used as the default.
	PathAlias string `json:"path_alias,omitempty"`
	// CloneURI is the URI that is used to clone the
	// repository. If unset, will default to
	// `https://github.com/org/repo.git`.
	CloneURI string `json:"clone_uri,omitempty"`
	// SkipSubmodules determines if submodules should be
	// cloned when the job is run. Defaults to true.
	SkipSubmodules bool `json:"skip_submodules,omitempty"`

	// ExtraRefs are auxiliary repositories that
	// need to be cloned, determined from config
	ExtraRefs []kube.Refs `json:"extra_refs,omitempty"`

	// DecorationConfig holds configuration options for
	// decorating PodSpecs that users provide
	DecorationConfig *kube.DecorationConfig `json:"decoration_config,omitempty"`
}

// RetestPresubmits returns all presubmits that should be run given a /retest command.
// This is the set of all presubmits intersected with ((alwaysRun + runContexts) - skipContexts)
func (c *JobConfig) RetestPresubmits(fullRepoName string, skipContexts, runContexts map[string]bool) []Presubmit {
	var result []Presubmit
	if jobs, ok := c.Presubmits[fullRepoName]; ok {
		for _, job := range jobs {
			if skipContexts[job.Context] {
				continue
			}
			if job.AlwaysRun || job.RunIfChanged != "" || runContexts[job.Context] {
				result = append(result, job)
			}
		}
	}
	return result
}

// GetPresubmit returns the presubmit job for the provided repo and job name.
func (c *JobConfig) GetPresubmit(repo, jobName string) *Presubmit {
	presubmits := c.AllPresubmits([]string{repo})
	for i := range presubmits {
		ps := presubmits[i]
		if ps.Name == jobName {
			return &ps
		}
	}
	return nil
}

// SetPresubmits updates c.Presubmits to jobs, after compiling and validating their regexes.
func (c *JobConfig) SetPresubmits(jobs map[string][]Presubmit) error {
	nj := map[string][]Presubmit{}
	for k, v := range jobs {
		nj[k] = make([]Presubmit, len(v))
		copy(nj[k], v)
		if err := SetPresubmitRegexes(nj[k]); err != nil {
			return err
		}
	}
	c.Presubmits = nj
	return nil
}

// SetPostsubmits updates c.Postsubmits to jobs, after compiling and validating their regexes.
func (c *JobConfig) SetPostsubmits(jobs map[string][]Postsubmit) error {
	nj := map[string][]Postsubmit{}
	for k, v := range jobs {
		nj[k] = make([]Postsubmit, len(v))
		copy(nj[k], v)
		if err := SetPostsubmitRegexes(nj[k]); err != nil {
			return err
		}
	}
	c.Postsubmits = nj
	return nil
}

// listPresubmits list all the presubmit for a given repo including the run after success jobs.
func listPresubmits(ps []Presubmit) []Presubmit {
	var res []Presubmit
	for _, p := range ps {
		res = append(res, p)
		res = append(res, listPresubmits(p.RunAfterSuccess)...)
	}
	return res
}

// AllPresubmits returns all prow presubmit jobs in repos.
// if repos is empty, return all presubmits.
func (c *JobConfig) AllPresubmits(repos []string) []Presubmit {
	var res []Presubmit

	for repo, v := range c.Presubmits {
		if len(repos) == 0 {
			res = append(res, listPresubmits(v)...)
		} else {
			for _, r := range repos {
				if r == repo {
					res = append(res, listPresubmits(v)...)
					break
				}
			}
		}
	}

	return res
}

// listPostsubmits list all the postsubmits for a given repo including the run after success jobs.
func listPostsubmits(ps []Postsubmit) []Postsubmit {
	var res []Postsubmit
	for _, p := range ps {
		res = append(res, p)
		res = append(res, listPostsubmits(p.RunAfterSuccess)...)
	}
	return res
}

// AllPostsubmits returns all prow postsubmit jobs in repos.
// if repos is empty, return all postsubmits.
func (c *JobConfig) AllPostsubmits(repos []string) []Postsubmit {
	var res []Postsubmit

	for repo, v := range c.Postsubmits {
		if len(repos) == 0 {
			res = append(res, listPostsubmits(v)...)
		} else {
			for _, r := range repos {
				if r == repo {
					res = append(res, listPostsubmits(v)...)
					break
				}
			}
		}
	}

	return res
}

// AllPeriodics returns all prow periodic jobs.
func (c *JobConfig) AllPeriodics() []Periodic {
	var listPeriodic func(ps []Periodic) []Periodic
	listPeriodic = func(ps []Periodic) []Periodic {
		var res []Periodic
		for _, p := range ps {
			res = append(res, p)
			res = append(res, listPeriodic(p.RunAfterSuccess)...)
		}
		return res
	}

	return listPeriodic(c.Periodics)
}
