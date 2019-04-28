package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "confkv.go",
        "dao.go",
        "resource.go",
        "titans.go",
        "tree.go",
        "user_resource.go",
    ],
    importpath = "go-common/app/service/live/resource/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/resource/api/grpc/v1:go_default_library",
        "//app/service/live/resource/api/http/v1:go_default_library",
        "//app/service/live/resource/conf:go_default_library",
        "//app/service/live/resource/lrucache:go_default_library",
        "//app/service/live/resource/model:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
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
