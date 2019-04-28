package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "invite_test.go",
        "medal_test.go",
        "notify_test.go",
        "pendant_test.go",
        "service_test.go",
        "third_pendant_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/usersuit/conf:go_default_library",
        "//app/service/main/usersuit/dao/pendant:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "invite.go",
        "medal.go",
        "notify.go",
        "pendant.go",
        "point_flag.go",
        "service.go",
        "third_pendant.go",
    ],
    importpath = "go-common/app/service/main/usersuit/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/member/api/gorpc:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//app/service/main/member/model/block:go_default_library",
        "//app/service/main/point/model:go_default_library",
        "//app/service/main/point/rpc/client:go_default_library",
        "//app/service/main/usersuit/conf:go_default_library",
        "//app/service/main/usersuit/dao/invite:go_default_library",
        "//app/service/main/usersuit/dao/medal:go_default_library",
        "//app/service/main/usersuit/dao/pendant:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
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
