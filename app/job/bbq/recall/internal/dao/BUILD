package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "berserker.go",
        "bloomfilter.go",
        "dao.go",
        "hdfs.go",
        "inverted_index.go",
        "video.go",
    ],
    importpath = "go-common/app/job/bbq/recall/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/internal/conf:go_default_library",
        "//app/job/bbq/recall/internal/model:go_default_library",
        "//app/job/bbq/recall/proto:go_default_library",
        "//app/job/bbq/recall/proto/quality:go_default_library",
        "//app/service/bbq/recsys-recall/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//vendor/github.com/Dai0522/go-hash/bloomfilter:go_default_library",
        "//vendor/github.com/golang/snappy:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
