package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "local.go",
        "tag.go",
    ],
    importpath = "go-common/app/interface/main/app-tag/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-tag/conf:go_default_library",
        "//app/interface/main/app-tag/model:go_default_library",
        "//app/interface/main/app-tag/model/region:go_default_library",
        "//app/interface/main/app-tag/service/ping:go_default_library",
        "//app/interface/main/app-tag/service/tag:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
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
