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

package decorate

import (
	"fmt"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation"

	"k8s.io/test-infra/prow/clonerefs"
	"k8s.io/test-infra/prow/entrypoint"
	"k8s.io/test-infra/prow/gcsupload"
	"k8s.io/test-infra/prow/initupload"
	"k8s.io/test-infra/prow/kube"
	"k8s.io/test-infra/prow/pod-utils/clone"
	"k8s.io/test-infra/prow/pod-utils/downwardapi"
	"k8s.io/test-infra/prow/pod-utils/wrapper"
	"k8s.io/test-infra/prow/sidecar"
)

const (
	logMountName            = "logs"
	logMountPath            = "/logs"
	artifactsEnv            = "ARTIFACTS"
	artifactsPath           = logMountPath + "/artifacts"
	codeMountName           = "code"
	codeMountPath           = "/home/prow/go"
	gopathEnv               = "GOPATH"
	toolsMountName          = "tools"
	toolsMountPath          = "/tools"
	gcsCredentialsMountName = "gcs-credentials"
	gcsCredentialsMountPath = "/secrets/gcs"
)

// Labels returns a string slice with label consts from kube.
func Labels() []string {
	return []string{kube.ProwJobTypeLabel, kube.CreatedByProw, kube.ProwJobIDLabel}
}

// VolumeMounts returns a string slice with *MountName consts in it.
func VolumeMounts() []string {
	return []string{logMountName, codeMountName, toolsMountName, gcsCredentialsMountName}
}

// VolumeMountPaths returns a string slice with *MountPath consts in it.
func VolumeMountPaths() []string {
	return []string{logMountPath, codeMountPath, toolsMountPath, gcsCredentialsMountPath}
}

// LabelsAndAnnotationsForSpec returns a minimal set of labels to add to prowjobs or its owned resources.
//
// User-provided extraLabels and extraAnnotations values will take precedence over auto-provided values.
func LabelsAndAnnotationsForSpec(spec kube.ProwJobSpec, extraLabels, extraAnnotations map[string]string) (map[string]string, map[string]string) {
	jobNameForLabel := spec.Job
	if len(jobNameForLabel) > validation.LabelValueMaxLength {
		// TODO(fejta): consider truncating middle rather than end.
		jobNameForLabel = strings.TrimRight(spec.Job[:validation.LabelValueMaxLength], "-")
		logrus.Warnf("Cannot use full job name '%s' for '%s' label, will be truncated to '%s'",
			spec.Job,
			kube.ProwJobAnnotation,
			jobNameForLabel,
		)
	}
	labels := map[string]string{
		kube.CreatedByProw:     "true",
		kube.ProwJobTypeLabel:  string(spec.Type),
		kube.ProwJobAnnotation: jobNameForLabel,
	}
	if spec.Type != kube.PeriodicJob && spec.Refs != nil {
		labels[kube.OrgLabel] = spec.Refs.Org
		labels[kube.RepoLabel] = spec.Refs.Repo
		if len(spec.Refs.Pulls) > 0 {
			labels[kube.PullLabel] = strconv.Itoa(spec.Refs.Pulls[0].Number)
		}
	}

	for k, v := range extraLabels {
		labels[k] = v
	}

	// let's validate labels
	for key, value := range labels {
		if errs := validation.IsValidLabelValue(value); len(errs) > 0 {
			// try to use basename of a path, if path contains invalid //
			base := filepath.Base(value)
			if errs := validation.IsValidLabelValue(base); len(errs) == 0 {
				labels[key] = base
				continue
			}
			logrus.Warnf("Removing invalid label: key - %s, value - %s, error: %s", key, value, errs)
			delete(labels, key)
		}
	}

	annotations := map[string]string{
		kube.ProwJobAnnotation: spec.Job,
	}
	for k, v := range extraAnnotations {
		annotations[k] = v
	}

	return labels, annotations
}

// LabelsAndAnnotationsForJob returns a standard set of labels to add to pod/build/etc resources.
func LabelsAndAnnotationsForJob(pj kube.ProwJob) (map[string]string, map[string]string) {
	var extraLabels map[string]string
	if extraLabels = pj.ObjectMeta.Labels; extraLabels == nil {
		extraLabels = map[string]string{}
	}
	extraLabels[kube.ProwJobIDLabel] = pj.ObjectMeta.Name
	return LabelsAndAnnotationsForSpec(pj.Spec, extraLabels, nil)
}

// ProwJobToPod converts a ProwJob to a Pod that will run the tests.
func ProwJobToPod(pj kube.ProwJob, buildID string) (*v1.Pod, error) {
	if pj.Spec.PodSpec == nil {
		return nil, fmt.Errorf("prowjob %q lacks a pod spec", pj.Name)
	}

	rawEnv, err := downwardapi.EnvForSpec(downwardapi.NewJobSpec(pj.Spec, buildID, pj.Name))
	if err != nil {
		return nil, err
	}

	spec := pj.Spec.PodSpec.DeepCopy()
	spec.RestartPolicy = "Never"
	spec.Containers[0].Name = kube.TestContainerName

	// we treat this as false if unset, while kubernetes treats it as true if
	// unset because it was added in v1.6
	if spec.AutomountServiceAccountToken == nil {
		myFalse := false
		spec.AutomountServiceAccountToken = &myFalse
	}

	if pj.Spec.DecorationConfig == nil {
		spec.Containers[0].Env = append(spec.Containers[0].Env, kubeEnv(rawEnv)...)
	} else {
		if err := decorate(spec, &pj, rawEnv); err != nil {
			return nil, fmt.Errorf("error decorating podspec: %v", err)
		}
	}

	podLabels, annotations := LabelsAndAnnotationsForJob(pj)
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        pj.ObjectMeta.Name,
			Labels:      podLabels,
			Annotations: annotations,
		},
		Spec: *spec,
	}, nil
}

const cloneLogPath = "clone.json"

// CloneLogPath returns the path to the clone log file in the volume mount.
func CloneLogPath(logMount kube.VolumeMount) string {
	return filepath.Join(logMount.MountPath, cloneLogPath)
}

// Exposed for testing
const (
	cloneRefsName    = "clonerefs"
	cloneRefsCommand = "/clonerefs"
)

// cloneEnv encodes clonerefs Options into json and puts it into an environment variable
func cloneEnv(opt clonerefs.Options) ([]v1.EnvVar, error) {
	// TODO(fejta): use flags
	cloneConfigEnv, err := clonerefs.Encode(opt)
	if err != nil {
		return nil, err
	}
	return kubeEnv(map[string]string{clonerefs.JSONConfigEnvVar: cloneConfigEnv}), nil
}

// sshVolume converts a secret holding ssh keys into the corresponding volume and mount.
//
// This is used by CloneRefs to attach the mount to the clonerefs container.
func sshVolume(secret string) (kube.Volume, kube.VolumeMount) {
	var sshKeyMode int32 = 0400 // this is octal, so symbolic ref is `u+r`
	name := strings.Join([]string{"ssh-keys", secret}, "-")
	mountPath := path.Join("/secrets/ssh", secret)
	v := kube.Volume{
		Name: name,
		VolumeSource: kube.VolumeSource{
			Secret: &kube.SecretSource{
				SecretName:  secret,
				DefaultMode: &sshKeyMode,
			},
		},
	}

	vm := kube.VolumeMount{
		Name:      name,
		MountPath: mountPath,
		ReadOnly:  true,
	}

	return v, vm
}

// cookiefileVolumes converts a secret holding cookies into the corresponding volume and mount.
//
// Secret can be of the form secret-name/base-name or just secret-name.
// Here secret-name refers to the kubernetes secret volume to mount, and base-name refers to the key in the secret
// where the cookies are stored. The secret-name pattern is equivalent to secret-name/secret-name.
//
// This is used by CloneRefs to attach the mount to the clonerefs container.
// The returned string value is the path to the cookiefile for use with --cookiefile.
func cookiefileVolume(secret string) (kube.Volume, kube.VolumeMount, string) {
	// Separate secret-name/key-in-secret
	parts := strings.SplitN(secret, "/", 2)
	cookieSecret := parts[0]
	var base string
	if len(parts) == 1 {
		base = parts[0] // Assume key-in-secret == secret-name
	} else {
		base = parts[1]
	}
	var cookiefileMode int32 = 0400 // u+r
	vol := kube.Volume{
		Name: "cookiefile",
		VolumeSource: kube.VolumeSource{
			Secret: &kube.SecretSource{
				SecretName:  cookieSecret,
				DefaultMode: &cookiefileMode,
			},
		},
	}
	mount := kube.VolumeMount{
		Name:      vol.Name,
		MountPath: "/secrets/cookiefile", // append base to flag
		ReadOnly:  true,
	}
	return vol, mount, path.Join(mount.MountPath, base)
}

// CloneRefs constructs the container and volumes necessary to clone the refs requested by the ProwJob.
//
// The container checks out repositories specified by the ProwJob Refs to `codeMount`.
// A log of what it checked out is written to `clone.json` in `logMount`.
//
// The container may need to mount SSH keys and/or cookiefiles in order to access private refs.
// CloneRefs returns a list of volumes containing these secrets required by the container.
func CloneRefs(pj kube.ProwJob, codeMount, logMount kube.VolumeMount) (*kube.Container, []kube.Refs, []kube.Volume, error) {
	if pj.Spec.DecorationConfig == nil {
		return nil, nil, nil, nil
	}
	if skip := pj.Spec.DecorationConfig.SkipCloning; skip != nil && *skip {
		return nil, nil, nil, nil
	}
	var cloneVolumes []kube.Volume
	var refs []kube.Refs // Do not return []*kube.Refs which we do not own
	if pj.Spec.Refs != nil {
		refs = append(refs, *pj.Spec.Refs)
	}
	for _, r := range pj.Spec.ExtraRefs {
		refs = append(refs, r)
	}
	if len(refs) == 0 { // nothing to clone
		return nil, nil, nil, nil
	}
	if codeMount.Name == "" || codeMount.MountPath == "" {
		return nil, nil, nil, fmt.Errorf("codeMount must set Name and MountPath")
	}
	if logMount.Name == "" || logMount.MountPath == "" {
		return nil, nil, nil, fmt.Errorf("logMount must set Name and MountPath")
	}

	var cloneMounts []kube.VolumeMount
	var sshKeyPaths []string
	for _, secret := range pj.Spec.DecorationConfig.SSHKeySecrets {
		volume, mount := sshVolume(secret)
		cloneMounts = append(cloneMounts, mount)
		sshKeyPaths = append(sshKeyPaths, mount.MountPath)
		cloneVolumes = append(cloneVolumes, volume)
	}

	var cloneArgs []string
	var cookiefilePath string

	if cp := pj.Spec.DecorationConfig.CookiefileSecret; cp != "" {
		v, vm, vp := cookiefileVolume(cp)
		cloneMounts = append(cloneMounts, vm)
		cloneVolumes = append(cloneVolumes, v)
		cookiefilePath = vp
		cloneArgs = append(cloneArgs, "--cookiefile="+cookiefilePath)
	}

	env, err := cloneEnv(clonerefs.Options{
		CookiePath:       cookiefilePath,
		GitRefs:          refs,
		GitUserEmail:     clonerefs.DefaultGitUserEmail,
		GitUserName:      clonerefs.DefaultGitUserName,
		HostFingerprints: pj.Spec.DecorationConfig.SSHHostFingerprints,
		KeyFiles:         sshKeyPaths,
		Log:              CloneLogPath(logMount),
		SrcRoot:          codeMount.MountPath,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("clone env: %v", err)
	}

	container := kube.Container{
		Name:         cloneRefsName,
		Image:        pj.Spec.DecorationConfig.UtilityImages.CloneRefs,
		Command:      []string{cloneRefsCommand},
		Args:         cloneArgs,
		Env:          env,
		VolumeMounts: append([]kube.VolumeMount{logMount, codeMount}, cloneMounts...),
	}
	return &container, refs, cloneVolumes, nil
}

func decorate(spec *kube.PodSpec, pj *kube.ProwJob, rawEnv map[string]string) error {
	rawEnv[artifactsEnv] = artifactsPath
	rawEnv[gopathEnv] = codeMountPath
	logMount := kube.VolumeMount{
		Name:      logMountName,
		MountPath: logMountPath,
	}
	logVolume := kube.Volume{
		Name: logMountName,
		VolumeSource: kube.VolumeSource{
			EmptyDir: &kube.EmptyDirVolumeSource{},
		},
	}

	codeMount := kube.VolumeMount{
		Name:      codeMountName,
		MountPath: codeMountPath,
	}
	codeVolume := kube.Volume{
		Name: codeMountName,
		VolumeSource: kube.VolumeSource{
			EmptyDir: &kube.EmptyDirVolumeSource{},
		},
	}

	toolsMount := kube.VolumeMount{
		Name:      toolsMountName,
		MountPath: toolsMountPath,
	}
	toolsVolume := kube.Volume{
		Name: toolsMountName,
		VolumeSource: kube.VolumeSource{
			EmptyDir: &kube.EmptyDirVolumeSource{},
		},
	}

	gcsCredentialsMount := kube.VolumeMount{
		Name:      gcsCredentialsMountName,
		MountPath: gcsCredentialsMountPath,
	}
	gcsCredentialsVolume := kube.Volume{
		Name: gcsCredentialsMountName,
		VolumeSource: kube.VolumeSource{
			Secret: &kube.SecretSource{
				SecretName: pj.Spec.DecorationConfig.GCSCredentialsSecret,
			},
		},
	}

	cloner, refs, cloneVolumes, err := CloneRefs(*pj, codeMount, logMount)
	if err != nil {
		return fmt.Errorf("could not create clonerefs container: %v", err)
	}
	if cloner != nil {
		spec.InitContainers = append([]kube.Container{*cloner}, spec.InitContainers...)
	}

	gcsOptions := gcsupload.Options{
		// TODO: pass the artifact dir here too once we figure that out
		GCSConfiguration:   pj.Spec.DecorationConfig.GCSConfiguration,
		GcsCredentialsFile: fmt.Sprintf("%s/service-account.json", gcsCredentialsMountPath),
		DryRun:             false,
	}

	initUploadOptions := initupload.Options{
		Options: &gcsOptions,
	}
	if cloner != nil {
		initUploadOptions.Log = CloneLogPath(logMount)
	}

	// TODO(fejta): use flags
	initUploadConfigEnv, err := initupload.Encode(initUploadOptions)
	if err != nil {
		return fmt.Errorf("could not encode initupload configuration as JSON: %v", err)
	}

	entrypointLocation := fmt.Sprintf("%s/entrypoint", toolsMountPath)

	spec.InitContainers = append(spec.InitContainers,
		kube.Container{
			Name:    "initupload",
			Image:   pj.Spec.DecorationConfig.UtilityImages.InitUpload,
			Command: []string{"/initupload"},
			Env: kubeEnv(map[string]string{
				initupload.JSONConfigEnvVar: initUploadConfigEnv,
				downwardapi.JobSpecEnv:      rawEnv[downwardapi.JobSpecEnv], // TODO: shouldn't need this?
			}),
			VolumeMounts: []kube.VolumeMount{logMount, gcsCredentialsMount},
		},
		kube.Container{
			Name:         "place-tools",
			Image:        pj.Spec.DecorationConfig.UtilityImages.Entrypoint,
			Command:      []string{"/bin/cp"},
			Args:         []string{"/entrypoint", entrypointLocation},
			VolumeMounts: []kube.VolumeMount{toolsMount},
		},
	)

	wrapperOptions := wrapper.Options{
		ProcessLog: fmt.Sprintf("%s/process-log.txt", logMountPath),
		MarkerFile: fmt.Sprintf("%s/marker-file.txt", logMountPath),
	}
	// TODO(fejta): use flags
	entrypointConfigEnv, err := entrypoint.Encode(entrypoint.Options{
		Args:        append(spec.Containers[0].Command, spec.Containers[0].Args...),
		Options:     &wrapperOptions,
		Timeout:     pj.Spec.DecorationConfig.Timeout,
		GracePeriod: pj.Spec.DecorationConfig.GracePeriod,
		ArtifactDir: artifactsPath,
	})
	if err != nil {
		return fmt.Errorf("could not encode entrypoint configuration as JSON: %v", err)
	}
	allEnv := rawEnv
	allEnv[entrypoint.JSONConfigEnvVar] = entrypointConfigEnv

	spec.Containers[0].Command = []string{entrypointLocation}
	spec.Containers[0].Args = []string{}
	spec.Containers[0].Env = append(spec.Containers[0].Env, kubeEnv(allEnv)...)
	spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, logMount, toolsMount)

	gcsOptions.Items = append(gcsOptions.Items, artifactsPath)
	// TODO(fejta): use flags
	sidecarConfigEnv, err := sidecar.Encode(sidecar.Options{
		GcsOptions:     &gcsOptions,
		WrapperOptions: &wrapperOptions,
	})
	if err != nil {
		return fmt.Errorf("could not encode sidecar configuration as JSON: %v", err)
	}

	spec.Containers = append(spec.Containers, kube.Container{
		Name:    "sidecar",
		Image:   pj.Spec.DecorationConfig.UtilityImages.Sidecar,
		Command: []string{"/sidecar"},
		Env: kubeEnv(map[string]string{
			sidecar.JSONConfigEnvVar: sidecarConfigEnv,
			downwardapi.JobSpecEnv:   rawEnv[downwardapi.JobSpecEnv], // TODO: shouldn't need this?
		}),
		VolumeMounts: []kube.VolumeMount{logMount, gcsCredentialsMount},
	})
	spec.Volumes = append(spec.Volumes, logVolume, toolsVolume, gcsCredentialsVolume)

	if len(refs) > 0 {
		spec.Containers[0].WorkingDir = clone.PathForRefs(codeMount.MountPath, refs[0])
		spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, codeMount)
		spec.Volumes = append(spec.Volumes, append(cloneVolumes, codeVolume)...)
	}

	return nil
}

// kubeEnv transforms a mapping of environment variables
// into their serialized form for a PodSpec, sorting by
// the name of the env vars
func kubeEnv(environment map[string]string) []v1.EnvVar {
	var keys []string
	for key := range environment {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var kubeEnvironment []v1.EnvVar
	for _, key := range keys {
		kubeEnvironment = append(kubeEnvironment, v1.EnvVar{
			Name:  key,
			Value: environment[key],
		})
	}

	return kubeEnvironment
}
