package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importmap = "go-common/app/tool/rider/sendEmail/cmd",
    importpath = "go-common/app/tool/ci/cmd/mail",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/ci/lib/mail:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
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

go_binary(
    name = "mail",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
