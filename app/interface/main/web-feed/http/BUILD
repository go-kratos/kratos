package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "feed.go",
        "http.go",
    ],
    importpath = "go-common/app/interface/main/web-feed/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-feed/conf:go_default_library",
        "//app/interface/main/web-feed/model:go_default_library",
        "//app/interface/main/web-feed/service:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/limit/aqm:go_default_library",
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
