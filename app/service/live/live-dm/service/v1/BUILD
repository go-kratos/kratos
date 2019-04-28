package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dM.go",
        "dmcheck.go",
        "getCheckInfo.go",
        "getDMConf.go",
        "limit.go",
        "send.go",
    ],
    importpath = "go-common/app/service/live/live-dm/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/av/api/liverpc/v1:go_default_library",
        "//app/service/live/banned_service/api/liverpc/v1:go_default_library",
        "//app/service/live/live-dm/api/grpc/v1:go_default_library",
        "//app/service/live/live-dm/conf:go_default_library",
        "//app/service/live/live-dm/dao:go_default_library",
        "//app/service/live/live-dm/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/ipipdotnet/ipdb-go:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dm_test.go",
        "dmcheck_test.go",
        "getCheckInfo_test.go",
        "getDMConf_test.go",
        "limit_test.go",
        "send_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/live-dm/api/grpc/v1:go_default_library",
        "//app/service/live/live-dm/conf:go_default_library",
        "//app/service/live/live-dm/dao:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/metadata:go_default_library",
    ],
)
