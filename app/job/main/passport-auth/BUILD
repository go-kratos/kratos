package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
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
        "//app/job/main/passport-auth/cmd:all-srcs",
        "//app/job/main/passport-auth/conf:all-srcs",
        "//app/job/main/passport-auth/dao:all-srcs",
        "//app/job/main/passport-auth/http:all-srcs",
        "//app/job/main/passport-auth/model:all-srcs",
        "//app/job/main/passport-auth/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
