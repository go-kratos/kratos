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
    data = ["databus-test.toml"],
    importpath = "go-common/app/infra/databus/cmd",
    tags = ["automanaged"],
    deps = [
        "//app/infra/databus/conf:go_default_library",
        "//app/infra/databus/http:go_default_library",
        "//app/infra/databus/service:go_default_library",
        "//app/infra/databus/tcp:go_default_library",
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
