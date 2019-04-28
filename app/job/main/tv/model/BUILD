package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
        "//app/job/main/tv/model/common:all-srcs",
        "//app/job/main/tv/model/pgc:all-srcs",
        "//app/job/main/tv/model/report:all-srcs",
        "//app/job/main/tv/model/ugc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
