package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "api_test.go",
        "dao_test.go",
        "redis_test.go",
        "type_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/videoup/conf:go_default_library",
        "//app/interface/main/videoup/model/archive:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "dao.go",
        "redis.go",
        "type.go",
    ],
    importpath = "go-common/app/interface/main/videoup/dao/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/dao/tool:go_default_library",
        "//app/interface/main/videoup/conf:go_default_library",
        "//app/interface/main/videoup/model/archive:go_default_library",
        "//app/interface/main/videoup/model/porder:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
