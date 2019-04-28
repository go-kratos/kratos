package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["dao.go"],
    importpath = "go-common/app/interface/live/web-ucenter/dao/user",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/web-ucenter/conf:go_default_library",
        "//app/interface/live/web-ucenter/dao:go_default_library",
        "//app/service/live/rankdb/api/liverpc/v1:go_default_library",
        "//app/service/live/rc/api/liverpc/v1:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/web-ucenter/conf:go_default_library",
        "//app/interface/live/web-ucenter/dao:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
