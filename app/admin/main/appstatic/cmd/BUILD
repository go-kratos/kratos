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
    data = ["appstatic-admin-test.toml"],
    importpath = "go-common/app/admin/main/appstatic/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/appstatic/conf:go_default_library",
        "//app/admin/main/appstatic/http:go_default_library",
        "//app/admin/main/appstatic/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/os/signal:go_default_library",
        "//library/syscall:go_default_library",
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
