package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bloomfilter.go",
        "dao.go",
        "parallel.go",
        "recall.go",
        "relation.go",
        "user.go",
    ],
    importpath = "go-common/app/service/bbq/recsys/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/recsys-recall/api/grpc/v1:go_default_library",
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//app/service/bbq/recsys/conf:go_default_library",
        "//app/service/bbq/recsys/dao/parallel:go_default_library",
        "//app/service/bbq/recsys/model:go_default_library",
        "//app/service/bbq/recsys/service/retrieve:go_default_library",
        "//app/service/bbq/search/api/grpc/v1:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/main/relation/api:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//vendor/github.com/Dai0522/go-hash/bloomfilter:go_default_library",
        "//vendor/github.com/Dai0522/workpool:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
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
        "//app/service/bbq/recsys/dao/parallel:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
