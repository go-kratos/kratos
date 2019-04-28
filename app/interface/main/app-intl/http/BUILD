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
        "local.go",
        "playurl.go",
        "search.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/main/app-intl/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-card/model:go_default_library",
        "//app/interface/main/app-card/model/card:go_default_library",
        "//app/interface/main/app-card/model/card/ai:go_default_library",
        "//app/interface/main/app-intl/conf:go_default_library",
        "//app/interface/main/app-intl/model:go_default_library",
        "//app/interface/main/app-intl/model/feed:go_default_library",
        "//app/interface/main/app-intl/model/player:go_default_library",
        "//app/interface/main/app-intl/model/view:go_default_library",
        "//app/interface/main/app-intl/service/feed:go_default_library",
        "//app/interface/main/app-intl/service/player:go_default_library",
        "//app/interface/main/app-intl/service/search:go_default_library",
        "//app/interface/main/app-intl/service/view:go_default_library",
        "//library/conf/env:go_default_library",
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
