package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/model/search:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "reply.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/creative/service/reply",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/dao/account:go_default_library",
        "//app/interface/main/creative/dao/archive:go_default_library",
        "//app/interface/main/creative/dao/article:go_default_library",
        "//app/interface/main/creative/dao/elec:go_default_library",
        "//app/interface/main/creative/dao/music:go_default_library",
        "//app/interface/main/creative/dao/reply:go_default_library",
        "//app/interface/main/creative/dao/search:go_default_library",
        "//app/interface/main/creative/model/music:go_default_library",
        "//app/interface/main/creative/model/reply:go_default_library",
        "//app/interface/main/creative/model/search:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
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
