package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["agent.go"],
    importpath = "go-common/app/service/main/dapper/agent",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper/conf:go_default_library",
        "//app/service/main/dapper/pkg/deliver:go_default_library",
        "//app/service/main/dapper/pkg/diskqueue:go_default_library",
        "//app/service/main/dapper/server/udpcollect:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace/proto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
