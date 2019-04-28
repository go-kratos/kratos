package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "depend.go",
        "service.go",
        "util.go",
    ],
    importpath = "go-common/app/service/main/dapper-query/service",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper-query/api/v1:go_default_library",
        "//app/service/main/dapper-query/conf:go_default_library",
        "//app/service/main/dapper-query/dao:go_default_library",
        "//app/service/main/dapper-query/model:go_default_library",
        "//app/service/main/dapper-query/pkg/cltclient:go_default_library",
        "//app/service/main/dapper-query/service/opslog:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/dapper-query/service/opslog:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
