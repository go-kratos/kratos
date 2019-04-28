load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

package(default_visibility = ["//visibility:public"])

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
        "//app/tool/gorpc/goparser:all-srcs",
        "//app/tool/gorpc/input:all-srcs",
        "//app/tool/gorpc/model:all-srcs",
    ],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "go-common/app/tool/gorpc",
    tags = ["automanaged"],
    deps = [
        "//app/tool/gorpc/goparser:go_default_library",
        "//app/tool/gorpc/input:go_default_library",
        "//app/tool/gorpc/model:go_default_library",
        "//vendor/golang.org/x/tools/imports:go_default_library",
    ],
)

go_binary(
    name = "gorpc",
    embed = [":go_default_library"],
)
