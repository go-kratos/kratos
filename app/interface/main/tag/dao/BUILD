package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "grpc_channel_test.go",
        "grpc_tag_test.go",
        "rpc_test.go",
        "spam_redis_test.go",
        "tagarc_prid_redis_test.go",
        "tagarc_redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tag/conf:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "grpc_channel.go",
        "grpc_tag.go",
        "rpc.go",
        "spam_redis.go",
        "tagarc_prid_redis.go",
        "tagarc_redis.go",
    ],
    importpath = "go-common/app/interface/main/tag/dao",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tag/conf:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/tag/api:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
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
