load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["passport-job.toml"],
    importpath = "go-common/app/job/main/passport/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = [
        "//app/job/main/passport/conf:go_default_library",
        "//app/job/main/passport/http:go_default_library",
        "//app/job/main/passport/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
    ],
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
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
