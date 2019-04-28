package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "breaker.go",
        "collector.go",
        "detect.go",
        "operation_name.go",
        "process.go",
    ],
    importpath = "go-common/app/service/main/dapper/collector",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper/conf:go_default_library",
        "//app/service/main/dapper/dao:go_default_library",
        "//app/service/main/dapper/model:go_default_library",
        "//app/service/main/dapper/pkg/batchwrite:go_default_library",
        "//app/service/main/dapper/pkg/collect:go_default_library",
        "//app/service/main/dapper/pkg/collect/kafkacollect:go_default_library",
        "//app/service/main/dapper/pkg/collect/tcpcollect:go_default_library",
        "//app/service/main/dapper/pkg/pointwrite:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "breaker_test.go",
        "collector_test.go",
        "detect_test.go",
        "operation_name_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/dapper/model:go_default_library",
        "//library/net/trace/proto:go_default_library",
    ],
)
