load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
    srcs = ["conf.go"],
    importpath = "go-common/app/admin/main/search/conf",
    tags = ["automanaged"],
    deps = [
        "//library/conf:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode/tip:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/trace:go_default_library",
        "//vendor/github.com/BurntSushi/toml:go_default_library",
    ],
)
