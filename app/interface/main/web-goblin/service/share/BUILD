package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
        "share.go",
    ],
    importpath = "go-common/app/interface/main/web-goblin/service/share",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/web-goblin/conf:go_default_library",
        "//app/interface/main/web-goblin/dao/share:go_default_library",
        "//app/interface/main/web-goblin/model/share:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//app/service/main/usersuit/rpc/client:go_default_library",
        "//library/cache:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
