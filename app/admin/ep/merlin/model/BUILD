package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bilihub.go",
        "constants.go",
        "dashboard.go",
        "device.go",
        "dto.go",
        "merlin.go",
        "paas.go",
        "tree.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/merlin/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
