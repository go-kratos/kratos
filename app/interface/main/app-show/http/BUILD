package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "banner.go",
        "daily.go",
        "http.go",
        "local.go",
        "rank.go",
        "region.go",
        "show.go",
    ],
    importpath = "go-common/app/interface/main/app-show/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-show/conf:go_default_library",
        "//app/interface/main/app-show/model:go_default_library",
        "//app/interface/main/app-show/model/region:go_default_library",
        "//app/interface/main/app-show/model/show:go_default_library",
        "//app/interface/main/app-show/service/banner:go_default_library",
        "//app/interface/main/app-show/service/daily:go_default_library",
        "//app/interface/main/app-show/service/ping:go_default_library",
        "//app/interface/main/app-show/service/rank:go_default_library",
        "//app/interface/main/app-show/service/region:go_default_library",
        "//app/interface/main/app-show/service/show:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/net/metadata:go_default_library",
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
