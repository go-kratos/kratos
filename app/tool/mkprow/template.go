package main

import (
	"strings"
)

func JobName(template, label string) string {
	switch template {
	case "__bazel_build_job_name__":
		return strings.Replace(label, "/", "-", -1) + "-bazel-build"
	case "__bazel_test_job_name__":
		return strings.Replace(label, "/", "-", -1) + "-bazel-test"
	case "__go_linter_job_name__":
		return strings.Replace(label, "/", "-", -1) + "-lint"
	default:
		return strings.Replace(label, "/", "-", -1)
	}
}

func JobBazelPath(result, label string) string {
	if strings.Contains(label, "tool/") ||
		strings.Contains(label, "admin/") ||
		strings.Contains(label, "common/") ||
		strings.Contains(label, "infra/") ||
		strings.Contains(label, "interface/") ||
		strings.Contains(label, "job/") ||
		strings.Contains(label, "service/") ||
		strings.Contains(label, "tool/") {
		return strings.Replace(result, "<<bazel_dir_param>>", "app/"+label, -1)
	} else {
		if strings.Contains(label, "library/") {
			return strings.Replace(result, "<<bazel_dir_param>>", label, -1)
		} else {
			return "app"
		}
	}
}

func JobImage(template string) string {
	image, ok := GlobalStatue.Image[template]
	if ok {
		return image
	}
	return ""
}

func Trigger(triagger, label string) string {
	if strings.Contains(triagger, "__bazel_build_job_name__") {
		return strings.Replace(triagger, "__bazel_build_job_name__", strings.Replace(label, "/", "-", -1)+"-build", -1)
	}
	if strings.Contains(triagger, "__bazel_test_job_name__") {
		return strings.Replace(triagger, "__bazel_test_job_name__", strings.Replace(label, "/", "-", -1)+"-test", -1)
	}
	if strings.Contains(triagger, "__go_linter_job_name__") {
		return strings.Replace(triagger, "__go_linter_job_name__", strings.Replace(label, "/", "-", -1)+"-lint", -1)
	}
	return triagger
}
