package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "anniversary.go",
        "apply.go",
        "archive.go",
        "article.go",
        "author.go",
        "cards.go",
        "categories.go",
        "complaint.go",
        "creative.go",
        "favorite.go",
        "hotspots.go",
        "http.go",
        "like.go",
        "list.go",
        "notice.go",
        "rank.go",
        "read.go",
        "recommends.go",
        "search.go",
        "sentinel.go",
        "share.go",
        "users.go",
        "view.go",
    ],
    importpath = "go-common/app/interface/openplatform/article/http",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/conf:go_default_library",
        "//app/interface/openplatform/article/dao:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/service:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/antispam:go_default_library",
        "//library/net/http/blademaster/middleware/auth:go_default_library",
        "//library/net/http/blademaster/middleware/cache:go_default_library",
        "//library/net/http/blademaster/middleware/cache/store:go_default_library",
        "//library/net/http/blademaster/middleware/verify:go_default_library",
        "//library/net/metadata:go_default_library",
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
