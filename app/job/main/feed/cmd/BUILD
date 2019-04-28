package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = [
        "feed-job-example.toml",
        "feed-job-test.toml",
    ],
    importpath = "go-common/app/job/main/feed/cmd",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/feed/conf:go_default_library",
        "//app/job/main/feed/http:go_default_library",
        "//app/job/main/feed/service:go_default_library",
        "//library/log:go_default_library",
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
