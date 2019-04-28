package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "player_test.go",
        "policy_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/player/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "player.go",
        "playurl.go",
        "policy.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/player/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/dm2/model:go_default_library",
        "//app/interface/main/dm2/rpc/client:go_default_library",
        "//app/interface/main/history/model:go_default_library",
        "//app/interface/main/history/rpc/client:go_default_library",
        "//app/interface/main/player/conf:go_default_library",
        "//app/interface/main/player/dao:go_default_library",
        "//app/interface/main/player/model:go_default_library",
        "//app/interface/main/tag/model:go_default_library",
        "//app/interface/main/tag/rpc/client:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/rpc/client:go_default_library",
        "//app/service/main/location/model:go_default_library",
        "//app/service/main/location/rpc/client:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//app/service/main/resource/rpc/client:go_default_library",
        "//app/service/main/ugcpay/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
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
