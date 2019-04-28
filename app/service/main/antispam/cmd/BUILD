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
        "antispam-integration.toml",
        "antispam-test.toml",
    ],
    importpath = "go-common/app/service/main/antispam/cmd",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/antispam/conf:go_default_library",
        "//app/service/main/antispam/http:go_default_library",
        "//app/service/main/antispam/rpc/server:go_default_library",
        "//app/service/main/antispam/service:go_default_library",
        "//library/ecode/tip:go_default_library",
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
