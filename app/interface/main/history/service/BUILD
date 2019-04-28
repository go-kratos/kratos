package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "history_test.go",
        "service_test.go",
        "toview_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/history/conf:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "history.go",
        "history_service.go",
        "rpc.go",
        "service.go",
        "toview.go",
    ],
    importpath = "go-common/app/interface/main/history/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/history/conf:go_default_library",
        "//app/interface/main/history/dao/history:go_default_library",
        "//app/interface/main/history/dao/toview:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/favorite/api/gorpc:go_default_library",
        "//app/service/main/favorite/model:go_default_library",
        "//app/service/main/history/api/grpc:go_default_library",
        "//app/service/main/history/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
