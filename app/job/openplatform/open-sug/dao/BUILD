package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bfs.go",
        "dao.go",
        "image.go",
        "item.go",
        "item_mysql.go",
        "redis.go",
        "season_es.go",
    ],
    importpath = "go-common/app/job/openplatform/open-sug/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/openplatform/sug/model:go_default_library",
        "//app/job/openplatform/open-sug/conf:go_default_library",
        "//app/job/openplatform/open-sug/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/code.google.com/p/graphics-go/graphics:go_default_library",
        "//vendor/github.com/golang/freetype:go_default_library",
        "//vendor/golang.org/x/image/font:go_default_library",
        "//vendor/golang.org/x/image/math/fixed:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
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
