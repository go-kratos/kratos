package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["res_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/web-show/conf:go_default_library",
        "//app/interface/main/web-show/model/resource:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "ad.go",
        "refresh.go",
        "res.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/web-show/service/resource",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/web-show/conf:go_default_library",
        "//app/interface/main/web-show/dao/ad:go_default_library",
        "//app/interface/main/web-show/dao/bangumi:go_default_library",
        "//app/interface/main/web-show/dao/data:go_default_library",
        "//app/interface/main/web-show/dao/resource:go_default_library",
        "//app/interface/main/web-show/model/resource:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/rpc/client:go_default_library",
        "//app/service/main/resource/rpc/client:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
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
