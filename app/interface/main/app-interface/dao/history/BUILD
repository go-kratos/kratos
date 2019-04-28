package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["dao.go"],
    importpath = "go-common/app/interface/main/app-interface/dao/history",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-interface/conf:go_default_library",
        "//app/interface/main/app-interface/model/history:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//app/interface/main/history/rpc/client:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
