package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "http.go",
        "titans.go",
    ],
    importpath = "go-common/app/service/live/resource/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/resource/api/http/v1:go_default_library",
        "//app/service/live/resource/conf:go_default_library",
        "//app/service/live/resource/service:go_default_library",
        "//app/service/live/resource/service/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
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
