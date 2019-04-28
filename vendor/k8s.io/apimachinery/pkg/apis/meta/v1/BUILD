load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "controller_ref.go",
        "conversion.go",
        "doc.go",
        "duration.go",
        "generated.pb.go",
        "group_version.go",
        "helpers.go",
        "labels.go",
        "meta.go",
        "micro_time.go",
        "micro_time_proto.go",
        "register.go",
        "time.go",
        "time_proto.go",
        "types.go",
        "types_swagger_doc_generated.go",
        "watch.go",
        "zz_generated.deepcopy.go",
        "zz_generated.defaults.go",
    ],
    importmap = "k8s.io/test-infra/vendor/k8s.io/apimachinery/pkg/apis/meta/v1",
    importpath = "k8s.io/apimachinery/pkg/apis/meta/v1",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/google/gofuzz:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/conversion:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/fields:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/labels:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/selection:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/types:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/intstr:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/watch:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//sortkeys:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
