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
    data = ["tv-job-test.toml"],
    importpath = "go-common/app/job/main/tv/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/http:go_default_library",
        "//app/job/main/tv/service/pgc:go_default_library",
        "//app/job/main/tv/service/ugc:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
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
