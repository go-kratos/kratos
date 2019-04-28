package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "abtest.go",
        "filter.go",
        "preprocess.go",
        "rank.go",
        "recsys.go",
        "relatedrec.go",
        "service.go",
        "store.go",
        "upsrec.go",
    ],
    importpath = "go-common/app/service/bbq/recsys/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//app/service/bbq/recsys/conf:go_default_library",
        "//app/service/bbq/recsys/dao:go_default_library",
        "//app/service/bbq/recsys/model:go_default_library",
        "//app/service/bbq/recsys/service/postprocess:go_default_library",
        "//app/service/bbq/recsys/service/rank:go_default_library",
        "//app/service/bbq/recsys/service/retrieve:go_default_library",
        "//app/service/bbq/recsys/service/util:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/go-ego/murmur:go_default_library",
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
        "//app/service/bbq/recsys/service/postprocess:all-srcs",
        "//app/service/bbq/recsys/service/rank:all-srcs",
        "//app/service/bbq/recsys/service/retrieve:all-srcs",
        "//app/service/bbq/recsys/service/util:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
