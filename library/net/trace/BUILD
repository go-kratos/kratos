load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

package(default_visibility = ["//visibility:public"])

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
        "//library/net/trace/proto:all-srcs",
    ],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "go-common/library/net/trace",
    tags = ["manual"],
    deps = [
        "//library/conf/dsn:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/trace/proto:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "config_test.go",
        "context_test.go",
        "dapper_test.go",
        "marshal_test.go",
        "report_test.go",
        "sample_test.go",
        "span_test.go",
        "tag_test.go",
        "util_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)
