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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    data = ["search-admin-test.toml"],
    importpath = "go-common/app/admin/main/search/cmd",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/search/conf:go_default_library",
        "//app/admin/main/search/http:go_default_library",
        "//app/admin/main/search/service:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
    ],
)

go_binary(
    name = "cmd",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
