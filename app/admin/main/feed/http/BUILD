load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bangumi.go",
        "bfs.go",
        "channel.go",
        "common.go",
        "egg.go",
        "http.go",
        "popular.go",
        "search.go",
    ],
    importpath = "go-common/app/admin/main/feed/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/feed/conf:go_default_library",
        "//app/admin/main/feed/model/channel:go_default_library",
        "//app/admin/main/feed/model/common:go_default_library",
        "//app/admin/main/feed/model/egg:go_default_library",
        "//app/admin/main/feed/model/search:go_default_library",
        "//app/admin/main/feed/model/show:go_default_library",
        "//app/admin/main/feed/service/bfs:go_default_library",
        "//app/admin/main/feed/service/channel:go_default_library",
        "//app/admin/main/feed/service/common:go_default_library",
        "//app/admin/main/feed/service/egg:go_default_library",
        "//app/admin/main/feed/service/pgc:go_default_library",
        "//app/admin/main/feed/service/popular:go_default_library",
        "//app/admin/main/feed/service/search:go_default_library",
        "//app/admin/main/feed/util:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/episode/v1:go_default_library",
        "//app/service/openplatform/pgc-season/api/grpc/season/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
