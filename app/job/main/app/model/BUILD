package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "const.go",
        "model.go",
        "retry.go",
    ],
    importpath = "go-common/app/job/main/app/model",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/app/model/space:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/time:go_default_library",
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
        "//app/job/main/app/model/feed:all-srcs",
        "//app/job/main/app/model/space:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
