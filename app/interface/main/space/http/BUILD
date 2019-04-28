package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "album.go",
        "app.go",
        "archive.go",
        "article.go",
        "assist.go",
        "bangumi.go",
        "blacklist.go",
        "channel.go",
        "coin.go",
        "dynamic.go",
        "fav.go",
        "http.go",
        "setting.go",
        "shop.go",
        "tag.go",
        "web.go",
    ],
    importpath = "go-common/app/interface/main/space/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/space/conf:go_default_library",
        "//app/interface/main/space/model:go_default_library",
        "//app/interface/main/space/service:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/supervisor:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/xstr:go_default_library",
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

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
