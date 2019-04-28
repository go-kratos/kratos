package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["http.go"],
    importpath = "go-common/app/interface/live/app-interface/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/app-interface/api/http/v1:go_default_library",
        "//app/interface/live/app-interface/api/http/v2:go_default_library",
        "//app/interface/live/app-interface/conf:go_default_library",
        "//app/interface/live/app-interface/dao:go_default_library",
        "//app/interface/live/app-interface/service:go_default_library",
        "//app/interface/live/app-interface/service/v1:go_default_library",
        "//app/interface/live/app-interface/service/v1/app_conf:go_default_library",
        "//app/interface/live/app-interface/service/v2:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
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
