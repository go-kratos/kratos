package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "build_defaults.go",
        "build_template_interface.go",
        "build_template_types.go",
        "build_template_validation.go",
        "build_types.go",
        "build_validation.go",
        "cluster_build_template_types.go",
        "cluster_build_template_validation.go",
        "doc.go",
        "metadata_validation.go",
        "register.go",
        "target_path_validation.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "github.com/knative/build/pkg/apis/build/v1alpha1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/knative/build/pkg/apis/build:go_default_library",
        "//vendor/github.com/knative/pkg/apis:go_default_library",
        "//vendor/github.com/knative/pkg/apis/duck/v1alpha1:go_default_library",
        "//vendor/github.com/knative/pkg/kmeta:go_default_library",
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
