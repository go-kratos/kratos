package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "message.go",
        "result.go",
        "retry.go",
        "stat.go",
    ],
    importpath = "go-common/app/job/main/archive-shjd/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/service/main/archive/api:go_default_library"],
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
