package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "archive.go",
        "article.go",
        "bnj.go",
        "broadcast.go",
        "coin.go",
        "dynamic.go",
        "elec.go",
        "feedback.go",
        "help.go",
        "http.go",
        "icon.go",
        "ip.go",
        "like.go",
        "nav.go",
        "newlist.go",
        "online.go",
        "ranking.go",
        "resource.go",
        "search.go",
        "tag.go",
        "wechat.go",
    ],
    importpath = "go-common/app/interface/main/web/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web/conf:go_default_library",
        "//app/interface/main/web/model:go_default_library",
        "//app/interface/main/web/service:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/cache:go_default_library",
        "//library/net/http/blademaster/middleware/cache/store:go_default_library",
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
