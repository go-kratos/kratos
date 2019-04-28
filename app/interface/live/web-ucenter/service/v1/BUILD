package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "anchorTask.go",
        "capsule.go",
        "history.go",
        "user.go",
    ],
    importpath = "go-common/app/interface/live/web-ucenter/service/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/web-ucenter/api/http:go_default_library",
        "//app/interface/live/web-ucenter/api/http/v1:go_default_library",
        "//app/interface/live/web-ucenter/conf:go_default_library",
        "//app/interface/live/web-ucenter/dao:go_default_library",
        "//app/interface/live/web-ucenter/dao/capsule:go_default_library",
        "//app/interface/live/web-ucenter/dao/history:go_default_library",
        "//app/interface/live/web-ucenter/dao/user:go_default_library",
        "//app/interface/live/web-ucenter/model:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//app/service/live/xrewardcenter/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
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
        "anchorTask_test.go",
        "history_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/web-ucenter/api/http/v1:go_default_library",
        "//app/interface/live/web-ucenter/conf:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
