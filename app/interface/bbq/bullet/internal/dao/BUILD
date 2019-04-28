package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "filter.go",
    ],
    importpath = "go-common/app/interface/bbq/bullet/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/bbq/bullet/api:go_default_library",
        "//app/interface/bbq/bullet/internal/conf:go_default_library",
        "//app/interface/bbq/bullet/internal/model:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/main/filter/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/warden:go_default_library",
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
