load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bloomfilter.go",
        "filter.go",
        "parallel.go",
        "ranker.go",
        "recall.go",
        "result.go",
        "score.go",
        "service.go",
        "task.go",
    ],
    importpath = "go-common/app/service/bbq/recsys-recall/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/proto:go_default_library",
        "//app/service/bbq/recsys-recall/api/grpc/v1:go_default_library",
        "//app/service/bbq/recsys-recall/conf:go_default_library",
        "//app/service/bbq/recsys-recall/dao:go_default_library",
        "//app/service/bbq/recsys-recall/model:go_default_library",
        "//app/service/bbq/recsys-recall/service/index:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/Dai0522/go-hash/bloomfilter:go_default_library",
        "//vendor/github.com/Dai0522/workpool:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
        "//app/service/bbq/recsys-recall/service/index:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
