package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "feature.go",
        "instance.go",
        "model.go",
        "rank.go",
    ],
    importpath = "go-common/app/service/bbq/recsys/service/rank",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//app/service/bbq/recsys/model:go_default_library",
        "//app/service/bbq/recsys/service/rank/treelite:go_default_library",
        "//app/service/bbq/recsys/service/retrieve:go_default_library",
        "//app/service/bbq/recsys/service/util:go_default_library",
        "//library/log:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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
        "//app/service/bbq/recsys/service/rank/treelite:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
