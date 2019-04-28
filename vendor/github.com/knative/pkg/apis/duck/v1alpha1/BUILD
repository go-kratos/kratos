package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "addressable_types.go",
        "condition_set.go",
        "conditions_types.go",
        "doc.go",
        "generational_types.go",
        "legacy_targetable_types.go",
        "register.go",
        "retired_targetable_types.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "github.com/knative/pkg/apis/duck/v1alpha1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/knative/pkg/apis:go_default_library",
        "//vendor/github.com/knative/pkg/apis/duck:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
