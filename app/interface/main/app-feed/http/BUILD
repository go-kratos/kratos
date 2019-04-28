package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "dynamic.go",
        "feed.go",
        "http.go",
        "local.go",
        "region.go",
    ],
    importpath = "go-common/app/interface/main/app-feed/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card:go_default_library",
        "//app/interface/main/app-card/model/card/ai:go_default_library",
        "//app/interface/main/app-card/model/card/operate:go_default_library",
        "//app/interface/main/app-feed/conf:go_default_library",
        "//app/interface/main/app-feed/model:go_default_library",
        "//app/interface/main/app-feed/model/feed:go_default_library",
        "//app/interface/main/app-feed/service/external:go_default_library",
        "//app/interface/main/app-feed/service/feed:go_default_library",
        "//app/interface/main/app-feed/service/ping:go_default_library",
        "//app/interface/main/app-feed/service/region:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
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
